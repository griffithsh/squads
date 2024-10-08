package embark

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/griffithsh/squads/game/item"
	"github.com/griffithsh/squads/graph"
	"github.com/griffithsh/squads/mathx"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/ui"
)

// Archive is what is required by embark of any archive data provider.
type Archive interface {
	Profession(profession string) *game.ProfessionDetails
	Names() map[string][]string
	Appearance(profession string, sex game.CharacterSex, hair string, skin string) *game.Appearance
	HairVariations() []string
	SkinVariations() []string
	PedestalAppearances(sinister bool) []int
}

const (
	groundSpritesZ  = 0
	grassSpriteZ    = 10
	roadwaySpritesZ = 50
	pathwaySpritesZ = 60
	houseSpritesZ   = 80
	uiSpritesZ      = 1000
)

var embarkPoints = []geom.Key{{M: 0, N: 1}, {M: -1, N: 0}, {M: -2, N: 0}, {M: 1, N: 0}, {M: 2, N: 0}}

type housePosition struct {
	x, y           float64
	villagerEntity ecs.Entity
}

// Manager holds state and provides methods to control that state for an embark
// screen. This screen allows the player to configure the Characters they will
// start their run with.
type Manager struct {
	mgr     *ecs.World
	bus     *event.Bus
	archive Archive

	screenW, screenH int

	houses []*housePosition

	taken    map[geom.Key]hexType
	field    *geom.Field
	searcher *graph.Searcher
}

// NewManager creates a new Manager in a default state. You should call Begin to start the Manager.
func NewManager(mgr *ecs.World, bus *event.Bus, archive Archive) *Manager {
	em := Manager{
		mgr:     mgr,
		bus:     bus,
		archive: archive,
		taken:   map[geom.Key]hexType{},
		field:   geom.NewField(10, 5, 12),
	}

	costs := func(gv1 graph.Vertex, gv2 graph.Vertex) float64 {
		k1 := gv1.(geom.Key)
		k2 := gv2.(geom.Key)
		cost := func(k geom.Key) float64 {
			ht, ok := em.taken[k]
			if !ok {
				return 500.0
			}
			switch ht {
			case blocked:
				return math.Inf(0)
			case roadway:
				return 0
			case pathway:
				return 25.0
			}
			return 500.0
		}

		return mathx.MinF64(cost(k1), cost(k2))
	}
	edges := func(gv graph.Vertex) []graph.Vertex {
		k := gv.(geom.Key)

		result := []graph.Vertex{}

		for k := range k.Neighbors() {
			if t, ok := em.taken[k]; ok && t == blocked {
				continue
			}
			result = append(result, k)
		}
		return result
	}
	guesses := func(v1, v2 graph.Vertex) float64 {
		squareDiff := func(a, b int) float64 {
			diff := a - b
			if diff < 0 {
				diff = -diff
			}
			return float64(diff * diff)
		}
		a, b := v1.(geom.Key), v2.(geom.Key)

		return squareDiff(a.M, b.M) + squareDiff(a.N, b.N)
	}
	em.searcher = graph.NewSearcher(costs, edges, guesses)

	bus.Subscribe(game.WindowSizeChanged{}.Type(), em.handleWindowSizeChanged)
	return &em
}

func randEdge(w, h int) geom.Key {
	switch rand.Intn(4) {
	case 1:
		return geom.Key{M: -w / 2, N: rand.Intn(h) - h/2}
	case 2:
		return geom.Key{M: w / 2, N: rand.Intn(h) - h/2}
	case 3:
		return geom.Key{M: rand.Intn(w) - w/2, N: h / 2}
	default:
		return geom.Key{M: rand.Intn(w) - w/2, N: -h / 2}
	}
}

func twoRandEdges(w, h int) (geom.Key, geom.Key) {
	a := randEdge(w, h)

	for {
		b := randEdge(w, h)

		if a.M == b.M || a.N == b.N {
			continue
		}
		return a, b
	}
}

func (em *Manager) addRoadwayEntity(key geom.Key, dir geom.DirectionType) {
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "embark")
	x, y := em.field.Ktow(key)
	layer := roadwaySpritesZ + int(dir)
	sx, sy := 0, 0
	switch dir {
	case geom.S:
		sx, sy = 0, 155
	case geom.N:
		sx, sy = 40, 167
	case geom.NW:
		sx, sy = 40, 155
	case geom.NE:
		sx, sy = 20, 167
	case geom.SE:
		sx, sy = 0, 167
	case geom.SW:
		sx, sy = 20, 155
	}
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x, Y: y,
		},
		Layer: layer,
	})
	em.mgr.AddComponent(e, &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: sx, Y: sy,
		W: 20, H: 12,
	})
}

func (em *Manager) addFeatureEntity(feat Feature, m, n, layer int) {
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "embark")
	x, y := em.field.Ktow(geom.Key{M: m, N: n})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x, Y: y,
		},
		Layer: layer,
	})
	em.mgr.AddComponent(e, &feat.Sprite)
}

func (em *Manager) rollHouses(villageW, villageH int) {
	// adds blocked and initial pathway values to taken
	// stores a Feature and a key to place it at - map[geom.Key]Feature
	for i := 0; i < 7; i++ {
		feat := houseFeatures[rand.Intn(len(houseFeatures))]
		func(feat Feature) {
			speculations := []geom.Key{
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
			}
			best := math.Inf(0)
			m, n := 0, 0
			for _, key := range speculations {
				// Is this key better than our previous best?
				x, y := em.field.Ktow(key)
				guess := (x * x) + (y * y)
				if best < guess {
					continue
				}

				// Will this feature fit at this key?
				for _, part := range feat.Occupies(key) {
					if ht, ok := em.taken[part]; ok {
						switch ht {
						case blocked, roadway, pathway:
							goto next
						}
					}
				}

				// Is the pathway connector of this feature open?
				if len(feat.PathConnect) > 0 {
					pathStart := feat.StartOfPathFor(key)
					if ht, ok := em.taken[pathStart]; ok {
						switch ht {
						case blocked:
							continue
						}
					}
				}

				// This speculation is the new best key to use!
				best = guess
				m, n = key.M, key.N
			next:
			}
			// We were very unlucky to not find any available places for this feature ...?
			if best == math.Inf(0) {
				return
			}

			// Add blocked and pathway values to em.taken.
			for _, part := range feat.Occupies(geom.Key{M: m, N: n}) {
				em.taken[part] = blocked
			}
			// If not a roadway, set it to pathway.
			if ht, ok := em.taken[feat.StartOfPathFor(geom.Key{M: m, N: n})]; !ok || ht != roadway {
				em.taken[feat.StartOfPathFor(geom.Key{M: m, N: n})] = pathway
			}

			// Add entity and components for the new feature.
			em.addFeatureEntity(feat, m, n, 100)

			// Save the location of this house for later ...
			x, y := em.field.Ktow(geom.Key{M: m, N: n})
			em.houses = append(em.houses, &housePosition{
				x: x,
				y: y,
			})
		}(feat)
	}
}

func (em *Manager) rollFlavor(villageW, villageH int) {
	for i := 0; i < int(float64(len(flavorFeatures))*1.5); i++ {
		feat := flavorFeatures[rand.Intn(len(flavorFeatures))]
		func(feat Feature) {
			speculations := []geom.Key{
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
				{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
			}
			best := math.Inf(0)
			m, n := 0, 0
			for _, key := range speculations {
				// Is this key better than our previous best?
				x, y := em.field.Ktow(key)
				guess := (x * x) + (y * y)
				if best < guess {
					continue
				}

				// Will this feature fit at this key?
				for _, part := range feat.Occupies(key) {
					if ht, ok := em.taken[part]; ok {
						switch ht {
						case blocked, roadway, pathway:
							goto next
						}
					}
				}

				// Is the pathway connector of this feature open?
				if len(feat.PathConnect) > 0 {
					pathStart := feat.StartOfPathFor(key)
					if ht, ok := em.taken[pathStart]; ok {
						switch ht {
						case blocked:
							continue
						}
					}
				}

				// This speculation is the new best key to use!
				best = guess
				m, n = key.M, key.N
			next:
			}
			// We were very unlucky to not find any available places for this feature ...?
			if best == math.Inf(0) {
				return
			}

			// Add blocked and pathway values to em.taken.
			for _, part := range feat.Occupies(geom.Key{M: m, N: n}) {
				em.taken[part] = blocked
			}

			// Add entity and components for the new feature.
			em.addFeatureEntity(feat, m, n, houseSpritesZ)
		}(feat)
	}
}

func (em *Manager) addPathways() {
	// By iterating through all in taken now, and picking the ones set to
	// pathway, we can find every unconnected feature. We have to make the
	// assumption that there are no other pathways right now though.
	origins := []geom.Key{}
	for k, ht := range em.taken {
		if ht == pathway {
			origins = append(origins, k)
		}
	}

	// TODO: should probably sort origins here to prevent random map iteration
	// causing differences.

	// Set every step to pathway unless it's already a roadway (roadways are
	// bigger and faster so they should take precedence).
	for _, origin := range origins {
		path := em.searcher.Search(origin, geom.Key{M: 0, N: 1})
		for _, step := range path {
			if ht, ok := em.taken[step.V.(geom.Key)]; ok && ht == roadway {
				continue
			}
			em.taken[step.V.(geom.Key)] = pathway
		}
	}

	// Add entities and components for every pathway hex that is now in
	// em.taken.
	for k, originHexType := range em.taken {
		if originHexType != pathway && originHexType != roadway {
			continue
		}

		// TODO: determine the correct combination of sprites to use given this
		// key's neighbors!
		// ...
		adjacent := k.Adjacent()

		for dir, neighbor := range adjacent {
			if ht, ok := em.taken[neighbor]; ok && ht == pathway || ht == roadway {
				if ht == roadway && originHexType == roadway {
					continue
				}
				e := em.mgr.NewEntity()
				em.mgr.Tag(e, "embark")
				x, y := em.field.Ktow(k)
				em.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: x, Y: y,
					},
					Layer: pathwaySpritesZ + int(dir),
				})
				em.mgr.AddComponent(e, spriteForPathway(dir, hashKeys(neighbor, k)))
			}
		}
	}
}

func spriteForPathway(dir geom.DirectionType, version int) *game.Sprite {
	alternate := version%2 != 0
	ne := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 0, Y: 104,
		W: 26, H: 18,
	}
	e := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 26, Y: 104,
		W: 26, H: 18,
	}
	se := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 52, Y: 104,
		W: 26, H: 18,
	}
	sw := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 52, Y: 122,
		W: 26, H: 18,
	}
	w := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 26, Y: 122,
		W: 26, H: 18,
	}
	nw := &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 0, Y: 122,
		W: 26, H: 18,
	}
	switch dir {
	case geom.NE:
		if alternate {
			return ne
		}
		return e
	case geom.SW:
		if alternate {
			return w
		}
		return sw
	case geom.NW:
		if alternate {
			return w
		}
		return nw
	case geom.SE:
		if alternate {
			return se
		}
		return e
	case geom.N:
		if alternate {
			return ne
		}
		return nw
	default: // geom.S
		if alternate {
			return se
		}
		return sw
	}
}

func hashKeys(k1, k2 geom.Key) int {
	joined := []geom.Key{k1, k2}
	sort.Slice(joined, func(i, j int) bool {
		if joined[i].M != joined[j].M {
			return joined[i].M < joined[j].M
		}
		return joined[i].N < joined[j].N
	})
	h := fnv.New128()
	bytes := func(i int) []byte {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(i))

		return b
	}
	h.Write(bytes(joined[0].M))
	h.Write(bytes(joined[0].N))
	h.Write(bytes(joined[1].M))
	h.Write(bytes(joined[1].N))
	sum := h.Sum([]byte{})
	return int(binary.LittleEndian.Uint64(sum))
}

func (em *Manager) rollRoadway(villageW, villageH int) {
	edge1, edge2 := twoRandEdges(villageW, villageH)
	path := em.searcher.Search(geom.Key{M: 0, N: 1}, edge1)
	if path == nil {
		panic("there was no path for the roadway")
	}
	for _, step := range path {
		key := step.V.(geom.Key)
		em.taken[key] = roadway
	}
	path = em.searcher.Search(geom.Key{M: 0, N: 1}, edge2)
	if path == nil {
		panic("there was no path for the roadway")
	}
	for _, step := range path {
		key := step.V.(geom.Key)
		em.taken[key] = roadway
	}
	for origin, ht := range em.taken {
		if ht != roadway {
			continue
		}

		connected := false
		for dir, key := range origin.Adjacent() {
			neighborType, ok := em.taken[key]
			if ok && neighborType == roadway {
				em.addRoadwayEntity(key, dir)
			}
		}

		if !connected {
			// create entity, add sprite and pos components
			e := em.mgr.NewEntity()
			em.mgr.Tag(e, "embark")
			x, y := em.field.Ktow(origin)
			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x, Y: y,
				},
				Layer: roadwaySpritesZ,
			})
			em.mgr.AddComponent(e, &game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 0, Y: 143,
				W: 20, H: 12,
			})
		}
	}
}

func (em *Manager) addFaeGate() {
	m, n := 0, 0
	key := geom.Key{M: m, N: n}
	em.taken[key] = blocked
	for _, dir := range faeGateFeature.Coverage {
		key = key.Adjacent()[dir]
		em.taken[key] = blocked
	}
	for _, k := range embarkPoints {
		em.taken[k] = roadway
	}
	em.addFeatureEntity(faeGateFeature, m, n, houseSpritesZ)
}

func (em *Manager) addMulliganHouse(villageW, villageH int) {
	speculations := []geom.Key{
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
		{M: rand.Intn(villageW) - villageW/2, N: rand.Intn(villageH) - villageH/2},
	}
	feat := windmillFeature
	best := math.Inf(0)
	m, n := 0, 0
	for _, key := range speculations {
		// Is this key better than our previous best?
		x, y := em.field.Ktow(key)
		guess := (x * x) + (y * y)
		if best < guess {
			continue
		}

		// Will this feature fit at this key?
		for _, part := range feat.Occupies(key) {
			if ht, ok := em.taken[part]; ok {
				switch ht {
				case blocked, roadway, pathway:
					goto next
				}
			}
		}

		// Is the pathway connector of this feature open?
		if len(feat.PathConnect) > 0 {
			pathStart := feat.StartOfPathFor(key)
			if ht, ok := em.taken[pathStart]; ok {
				switch ht {
				case blocked:
					continue
				}
			}
		}

		// This speculation is the new best key to use!
		best = guess
		m, n = key.M, key.N
	next:
	}
	// We were very unlucky to not find any available places for this feature ...?
	if best == math.Inf(0) {
		return
	}

	// Add blocked values to em.taken for every hex occupied.
	for _, part := range feat.Occupies(geom.Key{M: m, N: n}) {
		em.taken[part] = blocked
	}

	em.addFeatureEntity(feat, m, n, houseSpritesZ)
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "embark")
	em.mgr.Tag(e, "embark-villager-buttons")
	x, y := em.field.Ktow(geom.Key{M: m, N: n})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x, Y: y,
		},
		Layer: uiSpritesZ,
	})
	em.mgr.AddComponent(e, &ui.Interactive{
		W: 32, H: 24,
		Trigger: func(float64, float64) {
			fmt.Println("mulligan!")
		},
	})

	em.mgr.AddComponent(e, &game.Sprite{
		Texture: "embark/village-tiles.png",

		X: 16, Y: 244,
		W: 16, H: 12,
	})
	em.mgr.AddComponent(e, &game.HoverAnimation{})
}

// Begin an embark Manager, setting up Entities required to display and interact
// with the embark screen.
func (em *Manager) Begin() {
	const villageW = 48
	const villageH = 36
	em.taken = map[geom.Key]hexType{}

	em.addFaeGate()

	em.rollHouses(villageW, villageH)
	em.rollRoadway(villageW, villageH)
	em.addPathways()
	em.rollFlavor(villageW, villageH)
	em.addMulliganHouse(villageW, villageH)
	em.rollVillagers(5)

	em.repaint()

	// Add background tiles.
	// expand taken keys.
	for i := 0; i < 6; i++ {
		toAdd := []geom.Key{}
		for k := range em.taken {
			for neighbor := range k.Neighbors() {
				if _, ok := em.taken[neighbor]; ok {
					continue
				}
				if neighbor.M > villageW/2 || neighbor.M < -villageW/2 {
					continue
				}
				if neighbor.N > villageH/2 || neighbor.N < -villageH/2 {
					continue
				}
				// Cache additions to add after a full sweep.
				toAdd = append(toAdd, neighbor)
			}
		}
		var ht hexType = clear
		if i >= 2 {
			ht = grassy
		}
		if i >= 3 {
			ht = bushes
		}
		if i >= 5 {
			ht = trees
		}
		// Add the cached keys.
		for _, key := range toAdd {
			em.taken[key] = ht
		}
	}
	// Add some short-grassy textured hexagonal tiles for the background.
	for key, ht := range em.taken {
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")

		x, y := em.field.Ktow(key)
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x, Y: y,
			},
			Layer: groundSpritesZ,
		})
		spr := game.Sprite{
			Texture: "embark/village-tiles.png",

			X: 60, Y: 0,
			W: 20, H: 24,

			OffsetY: -6,
		}
		switch rand.Intn(4) {
		case 1:
			spr.X = 80
		case 2:
			spr.X = 100
		case 3:
			spr.X = 120
		}
		em.mgr.AddComponent(e, &spr)

		if ht == clear {
			luck := rand.Intn(10)
			if luck < 2 {
				e := em.mgr.NewEntity()
				em.mgr.Tag(e, "embark")

				em.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: x, Y: y,
					},
					Layer: grassSpriteZ,
				})
				spr := game.Sprite{
					Texture: "embark/village-tiles.png",

					X: 0, Y: 179,
					W: 20, H: 24,

					OffsetY: -6,
				}
				switch luck {
				case 1:
					spr.X = 20
				case 2:
					spr.X = 40
				}
				em.mgr.AddComponent(e, &spr)
			}
		} else if ht == grassy {
			e := em.mgr.NewEntity()
			em.mgr.Tag(e, "embark")

			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x, Y: y,
				},
				Layer: grassSpriteZ,
			})
			spr := game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 140, Y: 0,
				W: 20, H: 24,

				OffsetY: -6,
			}
			switch rand.Intn(3) {
			case 1:
				spr.X = 160
			case 2:
				spr.X = 180
			}
			em.mgr.AddComponent(e, &spr)
		} else if ht == bushes {
			e := em.mgr.NewEntity()
			em.mgr.Tag(e, "embark")

			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x, Y: y,
				},
				Layer: grassSpriteZ,
			})
			spr := game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 0, Y: 0,
				W: 20, H: 24,

				OffsetY: -6,
			}
			switch rand.Intn(8) {
			case 1:
				spr.X = 20
			case 2:
				spr.X = 40
			case 3:
				spr.X = 140
			case 4:
				spr.X = 160
			case 5:
				spr.X = 180
			case 6:
				spr.X = 200
			case 7:
				spr.X = 220
			}
			em.mgr.AddComponent(e, &spr)
		} else if ht == trees {
			e := em.mgr.NewEntity()
			em.mgr.Tag(e, "embark")

			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x, Y: y,
				},
				Layer: grassSpriteZ,
			})
			spr := game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 0, Y: 0,
				W: 20, H: 24,

				OffsetY: -6,
			}
			switch rand.Intn(5) {
			case 1:
				spr.X = 20
			case 2:
				spr.X = 40
			case 3:
				spr.X = 200
			case 4:
				spr.X = 220
			}
			em.mgr.AddComponent(e, &spr)
		}
	}
}

// End an embark Manager, resetting it to a default state.
func (em *Manager) End() {
	for _, e := range em.mgr.Tagged("embark") {
		em.mgr.DestroyEntity(e)
	}
}

func (em *Manager) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	em.screenW, em.screenH = wsc.NewW, wsc.NewH
}

// repaint synchronises the embarkation status of the villagers. It should be
// called after a change is made to who will embark.
func (em *Manager) repaint() {
	for _, e := range em.mgr.Tagged("embark-villager-buttons") {
		em.mgr.DestroyEntity(e)
	}

	takenEmbarkPoints := 0
	for _, house := range em.houses {
		if house.villagerEntity == 0 {
			continue
		}

		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.Tag(e, "embark-villager-buttons")

		embarking := em.mgr.Component(house.villagerEntity, "Embarking").(*Embarking)
		if embarking.Value {
			// Draw this villager near the gate...
			x, y := em.field.Ktow(embarkPoints[takenEmbarkPoints])
			takenEmbarkPoints++
			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x, Y: y,
				},
				Layer: uiSpritesZ,
			})
			em.mgr.AddComponent(e, &game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 0, Y: 244,
				W: 16, H: 12,
			})
			em.mgr.AddComponent(e, &ui.Interactive{
				W: 32, H: 24,
				Trigger: func(villager ecs.Entity) func(float64, float64) {
					return func(float64, float64) {
						em.mgr.Tag(villager, "embark-focus-villager")
						em.repaint()
					}
				}(house.villagerEntity),
			})
		} else {
			// Draw this villager at their house.
			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: house.x, Y: house.y,
				},
				Layer: uiSpritesZ,
			})
			em.mgr.AddComponent(e, &ui.Interactive{
				W: 32, H: 24,
				Trigger: func(villager ecs.Entity) func(float64, float64) {
					return func(float64, float64) {
						em.mgr.Tag(villager, "embark-focus-villager")
						em.repaint()
					}
				}(house.villagerEntity),
			})

			em.mgr.AddComponent(e, &game.Sprite{
				Texture: "embark/village-tiles.png",

				X: 0, Y: 244,
				W: 16, H: 12,
			})
			em.mgr.AddComponent(e, &game.HoverAnimation{})
		}
	}

	/*
		Houses may have a villager, villagers may be flagged for embarkation.
		Houses have an x/y.
		Ready spots have an x/y too!

		When we repaint, we iterate through every house, if a villager is
		present it is either drawn at the house or the embark zone near the
		gate. Each would have an Interactive that pops a modal that focuses
		on the villager.

		We must destroy every entity that visually represents a villager, but
		not the entities that *are* the villagers.
	*/

	// Then, if any villager is popped/focused/etc, we need to paint a modal window ...
	if villager := em.mgr.AnyTagged("embark-focus-villager"); villager != 0 {
		f, err := os.Open("output/demo.ui.xml")
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		nui := ui.NewUI(f)
		char := em.mgr.Component(villager, "Character").(*game.Character)
		equip := em.mgr.Component(villager, "Equipment").(*item.Equipment)
		prof := em.archive.Profession(char.Profession)
		app := em.archive.Appearance(char.Profession, char.Sex, char.Hair, char.Skin)

		data := AsCharacterSheetData(char, equip, prof, app)
		data.HandleCancel = func(string) {
			// Cancel - destroy the UI.
			em.mgr.RemoveTag(villager, "embark-focus-villager")
			em.mgr.DestroyEntity(e)
			em.repaint()
		}
		embarking := em.mgr.Component(villager, (&Embarking{}).Type()).(*Embarking)
		data.ActionButton = "Prepare"
		if embarking.Value {
			data.ActionButton = "Return"
		}
		data.HandleAction = func(string) {
			// Prepare or unprepare - toggle embarking status and destroy the UI.
			em.mgr.RemoveTag(villager, "embark-focus-villager")
			em.mgr.DestroyEntity(e)
			embarking.Value = !embarking.Value
			em.repaint()
		}
		nui.Data = data
		em.mgr.AddComponent(e, nui)

		// FIXME: Is this the best way to disable(or remove?) all the buttons in the embark?
		for _, e := range em.mgr.Tagged("embark-villager-buttons") {
			em.mgr.RemoveComponent(e, &ui.Interactive{})
		}

	} else if takenEmbarkPoints > 0 {
		// Else if no-one is popped, then check how many are embarked. If > 0, show an embark/go! button.
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.Tag(e, "embark-villager-buttons")
		f, err := os.Open("output/embark-start.xml")
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		uic := ui.NewUI(f)
		uic.Data = struct{ HandleStart func(string) }{func(string) {
			uiEntity := e
			em.bus.Publish(&SquadSelected{})

			e := em.mgr.NewEntity()
			em.mgr.Tag(e, "player")
			em.mgr.AddComponent(e, &game.Squad{})
			squad := em.mgr.Component(e, "Squad").(*game.Squad)
			players := game.NewTeam()
			apps := em.archive.PedestalAppearances(false)
			players.PedestalAppearance = apps[rand.Intn(len(apps))]
			em.mgr.AddComponent(e, players)

			// Add prepared villagers to the team and squad
			for _, house := range em.houses {
				if house.villagerEntity == 0 {
					continue
				}
				e := house.villagerEntity
				embarking := em.mgr.Component(e, "Embarking").(*Embarking)
				if !embarking.Value {
					continue
				}
				em.mgr.AddComponent(e, players)
				squad.Members = append(squad.Members, e)
				em.mgr.RemoveTag(e, "embark")
			}

			// Destroy the entity for the ui here.
			em.mgr.DestroyEntity(uiEntity)

			e = em.mgr.NewEntity()
			em.mgr.AddComponent(e, &game.DiagonalMatrixWipe{
				W: em.screenW, H: em.screenH,
				Obscuring: true,
				OnComplete: func() {
					em.bus.Publish(&Embarked{})
				},
			})
		}}
		em.mgr.AddComponent(e, uic)
	}
}

// rollVillagers removes any rolled Characters in this village and rolls new ones.
func (em *Manager) rollVillagers(num int) {
	if num > len(em.houses) {
		panic(fmt.Sprintf("insufficient houses(%d) for %d villagers", len(em.houses), num))
	}

	for _, house := range em.houses {
		if house.villagerEntity == 0 {
			continue
		}
		em.mgr.DestroyEntity(house.villagerEntity)
		house.villagerEntity = 0
	}

	g := newGenerator(em.archive)
	for i := 0; i < num; i++ {
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.AddComponent(e, g.generateChar())
		em.mgr.AddComponent(e, &item.Equipment{
			Weapon: g.generateWeapon(),
		})

		em.mgr.AddComponent(e, &Embarking{false})

		em.houses[i].villagerEntity = e
	}
}

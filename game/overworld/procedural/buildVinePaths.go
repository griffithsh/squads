package procedural

import (
	"fmt"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// buildVinePaths generates long maps, as a single correct pathway with
// alternative dead-ending branches is generated trending in a single direction.
// It is intended to implement beaches.
func buildVinePaths(seed int64, level int) (Paths, error) {
	prng := rand.New(rand.NewSource(seed))
	fmt.Println("buildVinePaths")
	directions := []geom.DirectionType{
		geom.NE,
		geom.SE,
		geom.S,
		geom.SW,
		geom.NW,
		geom.N,
	}

	// sunshine is the direction that vines grow in.
	sunshine := directions[DeterministicIndexOf(prng, directions)]
	fl := geom.Actualize(sunshine, geom.ForwardLeft)
	fr := geom.Actualize(sunshine, geom.ForwardRight)

	type tip struct {
		remaining int
		key       geom.Key
	}

	type branchChance struct {
		chance int
		relDir geom.RelativeDirection
		key    geom.Key
	}

	// rollOnward randomly picks a direction from rolled sunshine direction
	rollOnward := func(prng *rand.Rand) geom.DirectionType {
		roll := prng.Intn(10)
		if roll < 2 {
			return fl
		} else if roll < 4 {
			return fr
		} else {
			return sunshine
		}
	}

	filterTaken := func(chances []branchChance, result Paths, current geom.Key) []branchChance {
		i := 0
		l := len(chances)
		for i < l {
			rel := chances[i].relDir
			dir := geom.Actualize(sunshine, rel)
			chances[i].key = current.ToDirection(dir)
			if _, ok := result.Nodes[chances[i].key]; ok {
				chances = append(chances[:i], chances[i+1:]...)
				l--
			} else {
				i++
			}
		}
		return chances[:i]
	}

	randomFrom := func(prng *rand.Rand, chances []branchChance) (branchChance, []branchChance) {
		sum := 0
		for _, chance := range chances {
			sum += chance.chance
		}
		roll := prng.Intn(sum)
		i := 0
		l := len(chances)
		running := 0
		for i < l {
			if roll < running+chances[i].chance {
				// got it!
				break
			} else {
				running = running + chances[i].chance
				i++
			}
		}
		return chances[i], append(chances[0:i], chances[i+1:]...)
	}

	result := Paths{
		Nodes: map[geom.Key]Placement{},
	}

	current := geom.Key{}
	// start := current

	tips := []tip{}
	nextBranch := 2 + prng.Intn(4)
	for i := 0; i < 21+prng.Intn(7); i++ {
		nextBranch--
		if nextBranch > 0 {
			// grow the vine in the general direction of sunshine
			chances := []branchChance{
				{chance: 29, relDir: geom.ForwardLeft},
				{chance: 17, relDir: geom.ForwardRight},
				{chance: 60, relDir: geom.Forward},
			}
			chances = filterTaken(chances, result, current)
			if len(chances) == 0 {
				fmt.Printf("\tmain trunk is dead; all adjacent hexes are unavailable\n")
				break
			}
			trunk, _ := randomFrom(prng, chances)
			result.Connect(current, trunk.key)
			current = trunk.key
			continue

		}

		// this is a branching point! Will we succeed??? who knows
		nextBranch = 2 + prng.Intn(4)

		chances := []branchChance{
			{chance: 10, relDir: geom.BackLeft},
			{chance: 10, relDir: geom.BackRight},
			{chance: 20, relDir: geom.ForwardLeft},
			{chance: 20, relDir: geom.ForwardRight},
			{chance: 50, relDir: geom.Forward},
		}
		// filter taken spots
		chances = filterTaken(chances, result, current)
		if len(chances) == 0 {
			fmt.Printf("\tmain trunk is dead; all adjacent hexes are unavailable\n")
			break
		}
		trunk, chances := randomFrom(prng, chances)
		result.Connect(current, trunk.key)

		if len(chances) > 0 {
			branch, _ := randomFrom(prng, chances)
			tips = append(tips, tip{
				remaining: 2 + prng.Intn(4),
				key:       branch.key,
			})

			result.Connect(current, branch.key)
		} else {
			fmt.Printf("\tbranch from %v is dead; all adjacent hexes are unavailable", current)
		}
		current = trunk.key
	}

	for i, tip := range tips {
		for tip.remaining > 0 {
			tip.remaining--

			onward := rollOnward(prng)
			next := tip.key.ToDirection(onward)
			if _, ok := result.Nodes[next]; ok {
				// no good!
				fmt.Printf("\ttip %d/%d cannot grow, blocked at %v\n", i+1, len(tips), next)
				continue
			}
			result.Connect(tip.key, next)
			tip.key = next
		}
	}
	return result, nil
}

package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

/*
gradients:

6 hexes of terrain A, remainder terrain B (grassland leading to forest)

At least one hex of impassable water,2-3 hexes of sand, remainder grassland (coastline)

gradient edges with a straight line (not very useful!?)
gradient edges with a straight line plus random number from the range (X1 to X2)
gradient edges with increasing density of A, and decreasing density of B over X hexes (like dithering)

fill strategy:

noise pattern - 30% chance of A, 70% chance of B
cellular automata (forming swamps in grassland)

strategy: {
	// defines whether the gradient info runs across the long or the short side. LtR vs RtL etc are rolled.
	orientation: long, short

	// defines how terrain for arbitrary hex M,N are selected, referring to named fill strategies or terrains
	gradient:

	// map of names to fill strategy definitions
	fills:

	// map of array of terrain options keyed by an identifier
	terrains:
}
*/

/*
Could you have terrain selections differ based on whether there is a road there or not?

What if you wanted to have a dark forest level, with a river running through the
middle, and a bridge joining the sides? Could doodads implement that?
*/

/*
Radial gradients???!??

find center of all paths by averaging the x,y of all keys (and the strategy
could have an offset, for "ascending the caldera", OR force the center on the
start or end hex for "entering the forest")

center is the zero point, OR you might need to be able to set the zero point to
key closest to center for when you want to do a caldera curve etc.

the max point is (always, I guess) the furthest distance any path is from the center
*/

/*
What if the dark forest bridge problem had a "double-blobber" buildPath solution?

0,0 is the point at which you cross the bridge, two sets of paths at opposing
sides of that origin could be made to fan out in some way, and pick a start on
one side, and a goal on the other.

when terrain generation time came, there would have to be some way of specifying
a "river" that intersects 0,0 in some way, but does not cross paths in any other
location. A special bridge doodad could be used to overwrite 0,0?
*/

// Code uniquely identifies either a Terrain or a Fill for an overworld tile.
type Code string

// TerrainBuilder is anything that can generate terrain for a set of paths.
type TerrainBuilder interface {
	Build(prng *rand.Rand, paths Paths) map[geom.Key]Code
}

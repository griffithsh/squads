# Goals fof the squads project

## TODO

- [x] Add a HUD to prove the viability of faiface/pixel a little further
- [x] Add pathfinding so that Actors can navigate the game-board
- [ ] Command Queue for Actors
- [ ] Animations

## MAYDO

- [ ] Cache sprites in the renderer

## Scratch pad & ideas

I think there is a concept like an _intent_ or a command. This concept would capture the idea of a an Actor intending to "retreat from danger" or "approach the nearest opponent" or more concrete things like "move to M,N". These intent would also be things like "Use skill X", or "Use skill Y on nearest ally".

Intents would be translated to Actions like "move to the hex to the SW, then SW, then SW, then S" and "use skill Y on M,N".

I guess the line is that intents can be non-specific, Actions must be concrete.

Intents map to configurable pseudo-script elements in a UI that the player can set up to automate the behaviours their heroes perform.

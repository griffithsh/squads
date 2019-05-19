# Goals of the squads project

## TODO

- [x] Add a HUD to prove the viability of faiface/pixel a little further
- [x] Add pathfinding so that Actors can navigate the game-board
- [x] Replace faiface/pixel with ebiten
- [ ] Facing
- [x] Screen-relative rendering of ecs Components for a HUD
- [x] More construction in combat, rather than main
- [ ] Animations
- [ ] Command Queue for Actors
- [ ] Add debouncing of click events to main - Interaction should be called once per click

## MAYDO

## Scratch pad & ideas

I think there is a concept like an _intent_ or a command. This concept would capture the idea of a an Actor intending to "retreat from danger" or "approach the nearest opponent" or more concrete things like "move to M,N". These intent would also be things like "Use skill X", or "Use skill Y on nearest ally".

Intents would be translated to Actions like "move to the hex to the SW, then SW, then SW, then S" and "use skill Y on M,N".

I guess the line is that intents can be non-specific, Actions must be concrete.

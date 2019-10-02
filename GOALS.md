# Goals of the squads project

## TODO

- [x] Add a HUD to prove the viability of faiface/pixel a little further
- [x] Add pathfinding so that Actors can navigate the game-board
- [x] Replace faiface/pixel with ebiten
- [x] Screen-relative rendering of ecs Components for a HUD
- [x] More construction in combat, rather than main
- [x] Add debouncing of click events to main - Interaction should be called once per click
- [x] Add an "End Turn" button to the HUD
- [x] Preparation, ActionPoints, and anything else needed for turns-based combat
- [x] HUD shows updated information about Action Points remaining for the current actor
- [x] HUD shows a skills section
- [x] HUD shows updated information about turn order of the actors
- [x] A way to generate starting positions for every actor
- [x] Actors have a team
- [x] Animations scaffolding
  - [x] Frame/Sprite Animations
  - [x] TranslationAnimations
    - [x] Merge SpriteOffset into Sprite
    - [x] Add RenderOffset to replace it
    - [x] Floating Animation - makes the entity hover up and down (for a floating cursor?)
- [x] Add Move and cancel buttons to go in and out of move mode
- [x] Actors always have a boundary line to demarcate the hexes they occupy
- [x] Facing (for backstab damage)
- [x] Profession and Sex for every Actor and a mapping between those and a set of animations
- [x] Animations at a game-concept level
  - [x] A _Performance_ System
- [x] Moving animations sped up and slowed down by the Entity's Mover's speed
- [ ] Computer-controlled Teams
- [ ] Transitions between Combat, the Overworld, a village/embark mode, and a main menu/splash screen
  - [ ] Combat can signal its completion
  - [ ] Main menu can load a saved game

## MAYDO
- [x] TurnToken is a field of the combat manager
- [ ] Medium and Large Actor art
- [ ] Intent Queue(?) for AI Actors
- [x] Move Hex, Hex4, Hex7 split up to Field1, Field4, Field7, so that they can all return LogicalHex for At(), Get() etc.
- [ ] Negative coordinates should no longer wrap absolutely positioned Sprites - bottom or right-aligned renderable should be positioned via the HUDs copy of the screen center
- [ ] Structuralise the way systems are registered with an ecs.World, so that all things that update can be found in a consistent place
- [ ] Other Animation types
  - [ ] Fix the way hover animation goes wild when the game loses focus
  - [ ] Jump Animation - makes the entity go up then down, then auto-ends
  - [ ] Shake animation - makes the entity shake left and right (like for taking hits?)
  - [ ] Float-away animation - (like for damage amount floaters)

## Scratch pad & ideas

### Intents?

I think there is a concept like an _intent_ or a command. This concept would capture the idea of a an Actor intending to "retreat from danger" or "approach the nearest opponent" or more concrete things like "move to M,N". These intent would also be things like "Use skill X", or "Use skill Y on nearest ally".

Intents would be translated to Actions like "move to the hex to the SW, then SW, then SW, then S" and "use skill Y on M,N".

I guess the line is that intents can be non-specific, Actions must be concrete.

### Ownership versus effect

In an Entity Component System, is there some way of talking about the ownership of a thing versus the application of a thing? One Actor (a paladin?) might have an aura that affects allies to give them additional action points. The paladin owns the aura; the aura's lifecycle is tied to the paladin. But, the allies are affected by the aura. This seems like a separation between affect? and ownership.

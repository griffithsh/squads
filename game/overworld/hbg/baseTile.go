package hbg

import "github.com/griffithsh/squads/game/overworld/procedural"

type BaseTile struct {
	Code    procedural.Code `json:"code"`
	Texture string          `json:"texture"`

	Variations Options `json:"variations"`
}

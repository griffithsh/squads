package hbg

import (
	"encoding/json"
	"fmt"

	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
)

// EdgeCollection maps the direction of neighbors to the appropriate encroachment info.
// type EdgeCollection map[geom.DirectionType]Options
type EdgeCollection struct {
	Texture string                         `json:"texture"`
	Options map[geom.DirectionType]Options `json:"options"`
}

func (e EdgeCollection) MarshalJSON() ([]byte, error) {
	options := map[string][]Option{}
	for k, v := range e.Options {
		options[k.String()] = v
	}
	dummy := struct {
		Texture string              `json:"texture"`
		Options map[string][]Option `json:"options"`
	}{
		Texture: e.Texture,
		Options: options,
	}

	return json.Marshal(dummy)
}

// UnmarshalJSON
func (e *EdgeCollection) UnmarshalJSON(data []byte) error {
	var dummy struct {
		Texture string              `json:"texture"`
		Options map[string][]Option `json:"options"`
	}
	err := json.Unmarshal(data, &dummy)
	if err != nil {
		return err
	}
	options := map[geom.DirectionType]Options{}
	for k, v := range dummy.Options {
		dir, err := geom.DirectionTypeString(k)
		if err != nil {
			return fmt.Errorf("convert %q to direction: %v", k, err)
		}

		options[dir] = v
	}

	e.Texture = dummy.Texture
	e.Options = options

	return nil
}

type Offset struct {
	X int
	Y int
}

type Corner struct {
	Options Options
	Offset  Offset
	W       int
	H       int
}

type CornerCollection struct {
	Texture string
	Corners map[geom.DirectionType]Corner
}

func (e CornerCollection) MarshalJSON() ([]byte, error) {
	corners := map[string]Corner{}
	for k, v := range e.Corners {
		corners[k.String()] = v
	}
	dummy := struct {
		Texture string            `json:"texture"`
		Corners map[string]Corner `json:"corners"`
	}{
		Texture: e.Texture,
		Corners: corners,
	}

	return json.Marshal(dummy)
}

func (e *CornerCollection) UnmarshalJSON(data []byte) error {
	var dummy struct {
		Texture string            `json:"texture"`
		Corners map[string]Corner `json:"corners"`
	}
	err := json.Unmarshal(data, &dummy)
	if err != nil {
		return err
	}
	corners := map[geom.DirectionType]Corner{}
	for k, v := range dummy.Corners {
		dir, err := geom.DirectionTypeString(k)
		if err != nil {
			return fmt.Errorf("convert %q to direction: %v", k, err)
		}

		corners[dir] = v
	}

	e.Texture = dummy.Texture
	e.Corners = corners

	return nil
}

type Encroachment struct {
	Description string          `json:"description"`
	Over        procedural.Code `json:"over"`
	Adjacent    procedural.Code `json:"adjacent"`

	Edges   EdgeCollection   `json:"edges"`
	Corners CornerCollection `json:"corners"`
}

// EncroachmentsCollection is a mapping between the over and the adjacent codes
// of an encroachment, and the Encroachment itself. Keys are formatted as
// fmt.Sprintf("%s<-%s",over,adjacent). The idea is that there is one collection
// that is loaded in layers. First from core stuff built-in to the base game,
// then from a datafile if present, and also from an uncompressed directory
// structure potentially.
type EncroachmentsCollection map[string]Encroachment

func (ec EncroachmentsCollection) Get(over, adjacent procedural.Code) *Encroachment {
	key := fmt.Sprintf("%s<-%s", over, adjacent)

	if e, ok := ec[key]; ok {
		return &e
	}
	return nil
}

func (ec EncroachmentsCollection) Put(e Encroachment) {
	key := fmt.Sprintf("%s<-%s", e.Over, e.Adjacent)
	ec[key] = e
}

package data

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseSkills(t *testing.T) {
	r := strings.NewReader(`{
    "id": "basic-slash",
    "name": "Slash",
    "explanation": "Slash the target",
    "tags": [
        "Attack"
    ],
    "icon": {
        "texture": "hud2.png",
        "x": 0,
        "y": 0,
        "w": 24,
        "h": 24
    },
    "targeting": {
        "selectable": {
            "type": "SelectWithin",
            "minRange": 1,
            "maxRange": 1
        },
        "brush": {
            "type": "SingleHex"
        }
    },
    "costs": {
        "CostsActionPoints": 20
    },
    "attackChanceToHitModifier": -0.1,
    "effects": [
        {
            "when": 100,
            "whenPoint": "AttackApexTimingPoint",
            "what": [{
                "_type": "DamageEffect",
                "min": [
                    {
                        "operator": "AddOp",
                        "variable": "$DMG-MIN"
                    }
                ],
                "max": [
                    {
                        "operator": "MultOp",
                        "variable": "$DMG-MAX"
                    }
                ],
                "classification": "Attack",
                "damageType": "FireDamage"
            }]
        }
    ]
}`)
	got, err := parseSkills(r)
	if err != nil {
		t.Fatalf("parseSkills: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("want 1 skills, got %d skills", len(got))
	}

	gotSkill := got[0]

	buf := []byte{}
	b := bytes.NewBuffer(buf)
	enc := json.NewEncoder(b)

	if err := enc.Encode(&gotSkill); err != nil {
		t.Fatalf("encode gotten skill: %v", err)
	}

	encoded := strings.TrimSpace(b.String())

	want := `{"ID":"basic-slash","Name":"Slash","Explanation":"Slash the target","Tags":["Attack"],"Icon":{"Frames":[{"texture":"hud2.png","x":0,"y":0,"w":24,"h":24,"offsetX":0,"offsetY":0}],"Timings":[5000000000],"Pointer":0,"EndBehavior":0},"Targeting":{"Selectable":{"Type":1,"MinRange":1,"MaxRange":1},"Brush":{"Type":0,"MinRange":0,"MaxRange":0}},"Effects":[{"When":{},"What":[{"Min":[{"Operator":"AddOp","Variable":"$DMG-MIN"}],"Max":[{"Operator":"MultOp","Variable":"$DMG-MAX"}],"Classification":"Attack","DamageType":"FireDamage"}]}],"Costs":{"0":20},"AttackChanceToHitModifier":-0.1}`
	if encoded != want {
		t.Errorf("want:\n\t%s\ngot:\n\t%s", want, encoded)
	}
}

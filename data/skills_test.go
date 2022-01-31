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
    "targeting": "TargetAdjacent",
    "targetingBrush": "SingleHex",
    "costs": {
        "CostsActionPoints": 20
    },
    "effects": [
        {
            "when": 100,
            "whenPoint": "AttackApexTimingPoint",
            "what": {
                "_type": "DamageEffect",
                "min": [
                    {
                        "operator": "AddOp",
                        "variable": "$DMG-MIN"
                    }
                ],
                "max": [
                    {
                        "operator": "AddOp",
                        "variable": "$DMG-MAX"
                    }
                ],
                "classification": "Attack"
            }
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

	encoded := strings.TrimSuffix(b.String(), "\n")

	want := `{"ID":"basic-slash","Name":"Slash","Explanation":"Slash the target","Tags":[1],"Icon":{"Frames":[{"texture":"hud2.png","x":0,"y":0,"w":24,"h":24,"offsetX":0,"offsetY":0}],"Timings":[5000000000],"Pointer":0,"EndBehavior":0},"Targeting":1,"TargetingBrush":0,"Effects":[{"When":{},"What":{"Min":[{"Operator":0,"Variable":"$DMG-MIN"}],"Max":[{"Operator":0,"Variable":"$DMG-MAX"}],"Classification":0,"DamageType":0}}],"Costs":{"0":20}}`
	if encoded != want {
		t.Errorf("unexpected: \n%v", encoded)
	}
}

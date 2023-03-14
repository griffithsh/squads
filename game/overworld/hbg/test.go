package hbg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// testMarshaling sends an object on a round trip through json.Marshal,
// Unmarshal, and Marshal again and proves that the operation is stable by
// checking the two JSON outputs are identical.
func testMarshaling[V any](original V) error {
	b1, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	var unmarshaled V
	err = json.Unmarshal(b1, &unmarshaled)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	b2, err := json.MarshalIndent(unmarshaled, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	if !reflect.DeepEqual(b1, b2) {
		divider := "-------------------------------"
		return fmt.Errorf("not equal:\nwant: %s\n%s\ngot:  %s\n%s\n%s\n", divider, b1, divider, b2, divider)
	}
	fmt.Println(string(b2))
	return nil
}

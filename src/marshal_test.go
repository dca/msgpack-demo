package msgpack

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	mp := NewMsgpack()

	jsonObj := map[string]interface{}{
		"name":    "Dca Hsu",
		"age":     30,
		"friends": []string{"someone", "測試"},
	}

	m, err := mp.Marshal(jsonObj)
	if err != nil {
		fmt.Println("Error marshaling to MessagePack:", err)
		return
	}

	fmt.Println("MessagePack:", string(m))
	fmt.Printf("MessagePack: %x\n", m)
}

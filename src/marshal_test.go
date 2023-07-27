package msgpack

import (
	"testing"
	"github.com/stretchr/testify/require"
	gomsgpack "github.com/shamaton/msgpack/v2"

)

func TestMarshal(t *testing.T) {
	mp := NewMsgpack()

	testCases := []struct {
		desc string
		jsonObj map[string]interface{}
	}{
		{
			desc: "Test Case 1",
			jsonObj: map[string]interface{}{
				"name":    "Dca Hsu",
				"age":     30,
				"friends": []string{"someone", "測試"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m1, err := mp.Marshal(tc.jsonObj)
			require.NoError(t, err)

			m2, err := gomsgpack.Marshal(tc.jsonObj)
			require.NoError(t, err)
	
			require.Equal(t, m1, m2, "The two MessagePack byte should be equal")
		})
	}
}

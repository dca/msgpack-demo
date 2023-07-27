package msgpack

import (
	"encoding/json"
	"testing"

	gomsgpack "github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	mp := NewMsgpack()

	testCases := []struct {
		desc    string
		jsonStr string
	}{
		{
			desc: "Test Case - number",
			jsonStr: `{
				"age": 18
			}`,
		},
		{
			desc: "Test Case - array of numbers",
			jsonStr: `{
				"amounts": [100, 99.99]
			}`,
		},
		{
			desc: "Test Case - string",
			jsonStr: `{
				"name": "Dca Hsu"
			}`,
		},
		{
			desc: "Test Case - array of strings",
			jsonStr: `{
				"friends": [
					"John",
					"Mary"
				]
			}`,
		},
		{
			desc: "Test Case - true",
			jsonStr: `{
				"isGood": true
			}`,
		},
		{
			desc: "Test Case - false",
			jsonStr: `{
				"isGood": false
			}`,
		},
		{
			desc: "Test Case - nil",
			jsonStr: `{
				"isGood": null
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			var jsonObj map[string]interface{}

			err := json.Unmarshal([]byte(tc.jsonStr), &jsonObj)
			require.NoError(t, err)

			m1, err := mp.Marshal(jsonObj)
			require.NoError(t, err)

			m2, err := gomsgpack.Marshal(jsonObj)
			require.NoError(t, err)

			require.Equal(t, m1, m2, "The two MessagePack byte should be equal")
		})
	}
}

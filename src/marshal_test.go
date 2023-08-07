package msgpack

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	mp := NewMsgpack()

	testCases := []struct {
		desc     string
		jsonStr  string
		expected []byte
	}{
		{
			desc: "Test Case - number",
			jsonStr: `{
				"age": 18
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0x12},
		},
		{
			desc: "Test Case - array of numbers",
			jsonStr: `{
				"amounts": [100, 99.99]
			}`,
			expected: []byte{0x81, 0xa7, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x92, 0x64, 0xcb, 0x40, 0x58, 0xff, 0x5c, 0x28, 0xf5, 0xc2, 0x8f},
		},
		{
			desc: "Test Case - string",
			jsonStr: `{
				"name": "Dca Hsu"
			}`,
			expected: []byte{0x81, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa7, 0x44, 0x63, 0x61, 0x20, 0x48, 0x73, 0x75},
		},
		{
			desc: "Test Case - array of strings",
			jsonStr: `{
				"friends": [
					"John",
					"Mary"
				]
			}`,
			expected: []byte{0x81, 0xa7, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x92, 0xa4, 0x4a, 0x6f, 0x68, 0x6e, 0xa4, 0x4d, 0x61, 0x72, 0x79},
		},
		{
			desc: "Test Case - true",
			jsonStr: `{
				"isGood": true
			}`,
			expected: []byte{0x81, 0xa6, 0x69, 0x73, 0x47, 0x6f, 0x6f, 0x64, 0xc3},
		},
		{
			desc: "Test Case - false",
			jsonStr: `{
				"isGood": false
			}`,
			expected: []byte{0x81, 0xa6, 0x69, 0x73, 0x47, 0x6f, 0x6f, 0x64, 0xc2},
		},
		{
			desc: "Test Case - nil",
			jsonStr: `{
				"isGood": null
			}`,
			expected: []byte{0x81, 0xa6, 0x69, 0x73, 0x47, 0x6f, 0x6f, 0x64, 0xc0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			var jsonObj map[string]interface{}
			jsonBytes := []byte(tc.jsonStr)

			decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
			decoder.UseNumber()
			decoder.Decode(&jsonObj)

			m1, err := mp.Marshal(jsonObj)
			require.NoError(t, err)

			require.Equal(t, tc.expected, m1, "The two MessagePack byte should be equal")
		})
	}
}

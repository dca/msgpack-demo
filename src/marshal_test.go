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
			desc: "Test Case - positive fixint",
			jsonStr: `{
				"age": 18
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0x12},
		},
		{
			desc: "Test Case - negative fixint",
			jsonStr: `{
				"age": -18
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xee},
		},
		{
			desc: "Test Case - negative uint8",
			jsonStr: `{
				"age": 255
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xcc, 0xff},
		},
		{
			desc: "Test Case - negative uint16",
			jsonStr: `{
				"age": 65535
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xcd, 0xff, 0xff},
		},
		{
			desc: "Test Case - negative uint32",
			jsonStr: `{
				"age": 4294967295
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xce, 0xff, 0xff, 0xff, 0xff},
		},
		{
			desc: "Test Case - negative uint64",
			jsonStr: `{
				"age": 18446744073709551615
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			desc: "Test Case - negative int8",
			jsonStr: `{
				"age": -128
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xd0, 0x80},
		},
		{
			desc: "Test Case - negative int16",
			jsonStr: `{
				"age": -32768
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xd1, 0x80, 0x00},
		},
		{
			desc: "Test Case - negative int32",
			jsonStr: `{
				"age": -2147483648
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xd2, 0x80, 0x00, 0x00, 0x00},
		},
		{
			desc: "Test Case - negative int64",
			jsonStr: `{
				"age": -9223372036854775808
			}`,
			expected: []byte{0x81, 0xa3, 0x61, 0x67, 0x65, 0xd3, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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

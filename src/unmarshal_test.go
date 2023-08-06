package msgpack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	mp := NewMsgpack()

	testCases := []struct {
		desc            string
		inputBytes      []byte
		expectedJsonObj map[string]interface{}
	}{
		{
			desc:       "Test Case - string",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xa3, 0x61, 0x73, 0x64},
			expectedJsonObj: map[string]interface{}{
				"compact": "asd",
			},
		},
		{
			desc:       "Test Case - nil",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xc0},
			expectedJsonObj: map[string]interface{}{
				"compact": nil,
			},
		},
		{
			desc:       "Test Case - True",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xc3},
			expectedJsonObj: map[string]interface{}{
				"compact": true,
			},
		},
		{
			desc:       "Test Case - False",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xc2},
			expectedJsonObj: map[string]interface{}{
				"compact": false,
			},
		},
		{
			desc:       "Test Case - Array",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x92, 0xa1, 0x61, 0xa1, 0x62},
			expectedJsonObj: map[string]interface{}{
				"compact": []interface{}{"a", "b"},
			},
		},
		{
			desc:       "Test Case - PositiveInt",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x05},
			expectedJsonObj: map[string]interface{}{
				"compact": 5,
			},
		},
		{
			desc:       "Test Case - NegativeInt",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xfb},
			expectedJsonObj: map[string]interface{}{
				"compact": -5,
			},
		},
		{
			desc:       "Test Case - uint8",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xcc, 0x05},
			expectedJsonObj: map[string]interface{}{
				"compact": 5,
			},
		},
		{
			desc:       "Test Case - uint16",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xcd, 0x00, 0x05},
			expectedJsonObj: map[string]interface{}{
				"compact": 5,
			},
		}, {
			desc:       "Test Case - uint32",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xce, 0x00, 0x00, 0x00, 0x05},
			expectedJsonObj: map[string]interface{}{
				"compact": 5,
			},
		},
		{
			desc:       "Test Case - uint64",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xcf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
			expectedJsonObj: map[string]interface{}{
				"compact": 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			output, err := mp.Unmarshal(tc.inputBytes)
			require.NoError(t, err)

			jsonExpected, err := json.Marshal(tc.expectedJsonObj)
			require.NoError(t, err)

			jsonOutput, err := json.Marshal(output)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, string(jsonOutput), string(jsonExpected), "The two JSON string should be equal")
		})
	}
}

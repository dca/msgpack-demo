package msgpack

import (
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
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			jsonOutput, err := mp.Unmarshal(tc.inputBytes)
			require.NoError(t, err)
			require.Equal(t, tc.expectedJsonObj, jsonOutput, "The two JSON object should be equal")

		})
	}
}

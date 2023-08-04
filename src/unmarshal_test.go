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
			desc:       "Test Case - 1",
			inputBytes: []byte{0x81, 0xa7, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0xa3, 0x61, 0x73, 0x64},
			expectedJsonObj: map[string]interface{}{
				"compact": "asd",
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

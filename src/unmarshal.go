package msgpack

import (
	"fmt"
	MsgPackTypes "msgpack/src/types"
)

// Unmarshal converts data from MessagePack format to JSON format.
func (m *Msgpack) Unmarshal(data []byte) (map[string]interface{}, error) {
	var i int
	var jsonObj map[string]interface{}

	jsonObjOutput, err := m.decodeMsgpack(data, &jsonObj, &i)

	if err != nil {
		return nil, err
	}
	return jsonObjOutput.(map[string]interface{}), nil
}

func (m *Msgpack) decodeMsgpack(data []byte, jsonObj *map[string]interface{}, i *int) (interface{}, error) {
	currentByte := data[*i]

	switch {
	case isMsgPackTypeNil(currentByte):
		*i++
		return nil, nil

	case isMsgPackTypeTrue(currentByte):
		*i++
		return true, nil

	case isMsgPackTypeFalse(currentByte):
		*i++
		return false, nil

	case isMsgPackTypeString(currentByte):
		str, err := m.handleMsgPackTypeString(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return str, nil

	case isMsgPackTypeMap(currentByte):
		obj, err := m.handleMsgPackTypeMap(data, jsonObj, i)
		if err != nil {
			return nil, err
		}

		return obj, nil
	}
	return *jsonObj, nil
}

// isMsgPackTypeNil checks if the byte represents a MsgPack nil type.
func isMsgPackTypeNil(b byte) bool {
	return b == MsgPackTypes.Nil
}

// isMsgPackTypeTrue checks if the byte represents a MsgPack true type.
func isMsgPackTypeTrue(b byte) bool {
	return b == MsgPackTypes.True
}

// isMsgPackTypeFalse checks if the byte represents a MsgPack false type.
func isMsgPackTypeFalse(b byte) bool {
	return b == MsgPackTypes.False
}

// isMsgPackTypeFloat64 checks if the byte represents a MsgPack float64 type.
func isMsgPackTypeFloat64(b byte) bool {
	return b == MsgPackTypes.Float64
}

// isMsgPackTypeString checks if the byte represents a MsgPack string type.
func isMsgPackTypeString(b byte) bool {
	// Using binary mask to filter out the 3 MSBs and compare with the pattern 101xxxxx
	return (b & 0xE0) == MsgPackTypes.FixStr
}

// isMsgPackTypeArray checks if the byte represents a MsgPack array type.
func isMsgPackTypeArray(b byte) bool {
	// Using binary mask to filter out the 4 MSBs and compare with the pattern 1001xxxx
	return (b & 0xF0) == MsgPackTypes.FixArray
}

// isMsgPackTypeMap checks if the byte represents a MsgPack map type.
func isMsgPackTypeMap(b byte) bool {
	// Using binary mask to filter out the 4 MSBs and compare with the pattern 1000xxxx
	return (b & 0xF0) == MsgPackTypes.FixMap
}

func (m *Msgpack) handleMsgPackTypeNil() {

}

func (m *Msgpack) handleMsgPackTypeTrue() {

}

func (m *Msgpack) handleMsgPackTypeFalse() {

}

func (m *Msgpack) handleMsgPackTypeFloat64() {

}

func (m *Msgpack) handleMsgPackTypeString(data []byte, jsonObj *map[string]interface{}, i *int) (string, error) {
	currentByte := data[*i]
	*i++

	// parse string length
	strLen := int(currentByte & 0x1F) // 0x1F = 00011111

	// Ensure strLen bytes are available in data
	if *i+strLen > len(data) {
		return "", fmt.Errorf("insufficient data for string of length %d", strLen)
	}

	// Read string bytes
	strBytes := data[*i : *i+strLen]

	// Increment pointer by strLen
	*i += strLen

	return string(strBytes), nil
}

func (m *Msgpack) handleMsgPackTypeArray() {

}

func (m *Msgpack) handleMsgPackTypeMap(data []byte, jsonObj *map[string]interface{}, i *int) (map[string]interface{}, error) {

	currentByte := data[*i]
	*i++

	// creat new map
	deepJsonObj := make(map[string]interface{})

	// parse map length
	mapLen := int(currentByte & 0x0F) // 0x0F = 00001111
	// parse map key
	for j := 0; j < mapLen; j++ {
		// parse map key
		key, err := m.handleMsgPackTypeString(data, jsonObj, i)
		if err != nil {
			return nil, err
		}

		// parse map value
		value, err := m.decodeMsgpack(data, jsonObj, i)
		if err != nil {
			return nil, err
		}

		// append key and value to new map
		deepJsonObj[key] = value
	}

	return deepJsonObj, nil
}

func (m *Msgpack) handleMsgPackTypeInt() {

}

func (m *Msgpack) handleMsgPackTypeUint() {

}

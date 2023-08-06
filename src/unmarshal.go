package msgpack

import (
	"encoding/binary"
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

	case isMsgPackTypePositiveInt(currentByte):
		integer, err := m.handleMsgPackTypePositiveInt(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeNegativeInt(currentByte):
		integer, err := m.handleMsgPackTypeNegativeInt(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint8(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeUint8(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint16(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeUint16(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint32(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeUint32(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUInt64(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeUint64(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeArray(currentByte):
		arr, err := m.handleMsgPackTypeArray(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return arr, nil

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

func isMsgPackTypePositiveInt(b byte) bool {
	return (b & 0x80) == MsgPackTypes.FixIntPos
}

func isMsgPackTypeNegativeInt(b byte) bool {
	return (b & 0xE0) == MsgPackTypes.FixIntNeg
}

func isMsgPackTypeUint8(b byte) bool {
	return b == MsgPackTypes.Uint8
}

func isMsgPackTypeUint16(b byte) bool {
	return b == MsgPackTypes.Uint16
}

func isMsgPackTypeUint32(b byte) bool {
	return b == MsgPackTypes.Uint32
}

func isMsgPackTypeUInt64(b byte) bool {
	return b == MsgPackTypes.Uint64
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

func (m *Msgpack) handleMsgPackTypeArray(data []byte, jsonObj *map[string]interface{}, i *int) ([]interface{}, error) {
	currentByte := data[*i]
	*i++

	// parse array length
	arrLen := int(currentByte & 0x0F) // 0x0F = 00001111

	// create new array
	arr := make([]interface{}, arrLen)

	// create error variable
	var err error

	// parse array elements
	for j := 0; j < arrLen; j++ {
		// parse array element
		arr[j], err = m.decodeMsgpack(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
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

func (m *Msgpack) handleMsgPackTypePositiveInt(data []byte, jsonObj *map[string]interface{}, i *int) (int, error) {
	currentByte := data[*i]
	*i++

	// parse int
	return int(currentByte & 0x7F), nil // 0x7F = 01111111
}

func (m *Msgpack) handleMsgPackTypeNegativeInt(data []byte, jsonObj *map[string]interface{}, i *int) (int, error) {
	currentByte := data[*i]
	*i++

	// parse int
	value := int(int8(currentByte))
	return value, nil
}

func (m *Msgpack) handleMsgPackTypeUint8(data []byte, jsonObj *map[string]interface{}, i *int) (uint8, error) {
	currentByte := data[*i]
	*i++

	// parse uint8
	return uint8(currentByte), nil
}

func (m *Msgpack) handleMsgPackTypeUint16(data []byte, jsonObj *map[string]interface{}, i *int) (uint16, error) {
	// Ensure 2 bytes are available in data
	if *i+2 > len(data) {
		return 0, fmt.Errorf("insufficient data for uint16")
	}

	// Read uint16 bytes
	value := binary.BigEndian.Uint16(data[*i : *i+2])

	// Increment pointer by 2
	*i += 2

	return value, nil
}

func (m *Msgpack) handleMsgPackTypeUint32(data []byte, jsonObj *map[string]interface{}, i *int) (uint32, error) {
	// Ensure 4 bytes are available in data
	if *i+4 > len(data) {
		return 0, fmt.Errorf("insufficient data for uint32")
	}

	// Read uint32 bytes
	value := binary.BigEndian.Uint32(data[*i : *i+4])

	// Increment pointer by 4
	*i += 4

	return value, nil
}

func (m *Msgpack) handleMsgPackTypeUint64(data []byte, jsonObj *map[string]interface{}, i *int) (uint64, error) {
	// Ensure 8 bytes are available in data
	if *i+8 > len(data) {
		return 0, fmt.Errorf("insufficient data for uint64")
	}

	// Read uint64 bytes
	value := binary.BigEndian.Uint64(data[*i : *i+8])

	// Increment pointer by 8
	*i += 8

	return value, nil
}

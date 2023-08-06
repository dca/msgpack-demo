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
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.FixIntPos, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeNegativeInt(currentByte):
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.FixIntNeg, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint8(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.Uint8, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint16(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.Uint16, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUint32(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.Uint32, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case isMsgPackTypeUInt64(currentByte):
		*i++
		integer, err := m.handleMsgPackTypeIntFamily(MsgPackTypes.Uint64, data, i)
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

// handleMsgPackTypeIntFamily handles int family types
func (m *Msgpack) handleMsgPackTypeIntFamily(t int, data []byte, i *int) (interface{}, error) {

	var value interface{}
	var err error
	var bytes []byte

	currentByte := data[*i]
	length := 0

	switch t {

	case MsgPackTypes.FixIntPos:
		value = int(currentByte & 0x7F) // 0x7F = 01111111

	case MsgPackTypes.FixIntNeg:
		value = int(int8(currentByte))

	case MsgPackTypes.Uint8:
		length = 1
		bytes, err = m.getNextBytes(data, i, length)
		value = uint8(bytes[0])

	case MsgPackTypes.Uint16:
		length = 2
		bytes, err = m.getNextBytes(data, i, length)
		value = binary.BigEndian.Uint16(bytes)

	case MsgPackTypes.Uint32:
		length = 4
		bytes, err = m.getNextBytes(data, i, length)
		value = binary.BigEndian.Uint32(bytes)

	case MsgPackTypes.Uint64:
		length = 8
		bytes, err = m.getNextBytes(data, i, length)
		value = binary.BigEndian.Uint64(bytes)

	default:
		return nil, fmt.Errorf("unknown int family type %d", t)
	}

	if err != nil {
		return nil, err
	}
	return value, nil
}

// getNextBytes returns the next n bytes from data
func (m *Msgpack) getNextBytes(data []byte, i *int, length int) ([]byte, error) {
	if *i+length > len(data) {
		return nil, fmt.Errorf("data out of range")
	}

	value := data[*i : *i+length]
	*i += length

	return value, nil
}

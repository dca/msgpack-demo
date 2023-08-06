package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
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
	case MsgPackTypes.IsMsgPackTypeNil(currentByte):
		*i++
		return nil, nil

	case MsgPackTypes.IsMsgPackTypeTrue(currentByte):
		*i++
		return true, nil

	case MsgPackTypes.IsMsgPackTypeFalse(currentByte):
		*i++
		return false, nil

	case MsgPackTypes.IsMsgPackTypeString(currentByte):
		str, err := m.handleMsgPackTypeString(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return str, nil

	case MsgPackTypes.IsMsgPackTypePositiveInt(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.FixIntPos, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeNegativeInt(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.FixIntNeg, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeUint8(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Uint8, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeUint16(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Uint16, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeUint32(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Uint32, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeUInt64(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Uint64, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeInt8(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Int8, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeInt16(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Int16, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeInt32(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Int32, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeInt64(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Int64, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeFloat32(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Float32, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeFloat64(currentByte):
		integer, err := m.handleMsgPackTypeNumberFamily(MsgPackTypes.Float64, data, i)
		if err != nil {
			return nil, err
		}
		return integer, nil

	case MsgPackTypes.IsMsgPackTypeArray(currentByte):
		arr, err := m.handleMsgPackTypeArray(data, jsonObj, i)
		if err != nil {
			return nil, err
		}
		return arr, nil

	case MsgPackTypes.IsMsgPackTypeMap(currentByte):
		obj, err := m.handleMsgPackTypeMap(data, jsonObj, i)
		if err != nil {
			return nil, err
		}

		return obj, nil
	}
	return *jsonObj, nil
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

// handleMsgPackTypeNumberFamily handles number family types
func (m *Msgpack) handleMsgPackTypeNumberFamily(t int, data []byte, i *int) (interface{}, error) {

	currentByte := data[*i]
	*i++

	var value interface{}
	var err error
	var bytes []byte
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

	case MsgPackTypes.Int8:
		length = 1
		bytes, err = m.getNextBytes(data, i, length)
		value = int8(bytes[0])

	case MsgPackTypes.Int16:
		length = 2
		bytes, err = m.getNextBytes(data, i, length)
		value = int16(binary.BigEndian.Uint16(bytes))

	case MsgPackTypes.Int32:
		length = 4
		bytes, err = m.getNextBytes(data, i, length)
		value = int32(binary.BigEndian.Uint32(bytes))

	case MsgPackTypes.Int64:
		length = 8
		bytes, err = m.getNextBytes(data, i, length)
		value = int64(binary.BigEndian.Uint64(bytes))

	case MsgPackTypes.Float32:
		length = 4
		bytes, err = m.getNextBytes(data, i, length)
		value = math.Float32frombits(binary.BigEndian.Uint32(bytes))

	case MsgPackTypes.Float64:
		length = 8
		bytes, err = m.getNextBytes(data, i, length)
		value = math.Float64frombits(binary.BigEndian.Uint64(bytes))

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

package msgpack

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	MsgPackTypes "msgpack/src/types"
)

// Marshal
func (m *Msgpack) Marshal(data map[string]interface{}) ([]byte, error) {
	msg, err := m.encodeJSON(data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (m *Msgpack) encodeJSON(data map[string]interface{}) ([]byte, error) {
	var result []byte
	m.handleValue(&result, data)
	return result, nil
}

func (m *Msgpack) handleValue(result *[]byte, data interface{}) error {
	switch v := reflect.ValueOf(data); v.Kind() {

	case reflect.Bool:
		// add value
		if v.Bool() {
			*result = append(*result, byte(MsgPackTypes.True))
		} else {
			*result = append(*result, byte(MsgPackTypes.False))
		}

	case reflect.Ptr, reflect.Invalid:
		// add value
		*result = append(*result, byte(MsgPackTypes.Nil))

	case reflect.Float64:
		// add type prefix
		*result = append(*result, byte(MsgPackTypes.Float64))

		// add value
		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, v.Float())
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		*result = append(*result, buf.Bytes()...)

	case reflect.String:
		//
		if num, ok := data.(json.Number); ok {
			// as number case
			if _, err := m.encodeMsgPackTypeNumberFamily(result, num); err != nil {
				return err
			}

		} else {
			// as string case

			valueAsBytes := []byte(v.String())

			// add type prefix
			*result = append(*result, byte(len(valueAsBytes))+MsgPackTypes.FixStr)

			// add value
			*result = append(*result, valueAsBytes...)
		}

	case reflect.Slice, reflect.Array:
		// add type prefix
		*result = append(*result, byte(v.Len())+MsgPackTypes.FixArray)

		// add value
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			m.handleValue(result, elem)
		}

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String {
			keys := v.MapKeys()

			// add type prefix
			*result = append(*result, byte(len(keys))+MsgPackTypes.FixMap)

			// add value
			for _, key := range keys {
				keyAsBytes := []byte(key.String())

				// add type prefix for key
				*result = append(*result, byte(len(keyAsBytes))+MsgPackTypes.FixStr)

				// add content of key
				*result = append(*result, keyAsBytes...)

				// add value
				elem := v.MapIndex(key).Interface()
				m.handleValue(result, elem)
			}
		} else {
			fmt.Println("Non-string keys in map are not supported.")
		}
	default:
		fmt.Println("unhandle case.", v.Kind(), v)
	}

	return nil
}

func (m *Msgpack) encodeMsgPackTypeNumberFamily(result *[]byte, num json.Number) (*[]byte, error) {

	var buf bytes.Buffer

	if ui, err := strconv.ParseUint(num.String(), 10, 64); err == nil {
		switch {
		case ui <= 127:
			*result = append(*result, byte(ui))

		case ui <= math.MaxUint8:
			*result = append(*result, byte(MsgPackTypes.Uint8))
			*result = append(*result, byte(uint8(ui)))

		case ui <= math.MaxUint16:
			*result = append(*result, byte(MsgPackTypes.Uint16))
			binary.Write(&buf, binary.BigEndian, uint16(ui))
			*result = append(*result, buf.Bytes()...)

		case ui <= math.MaxUint32:
			*result = append(*result, byte(MsgPackTypes.Uint32))
			binary.Write(&buf, binary.BigEndian, uint32(ui))
			*result = append(*result, buf.Bytes()...)

		default: // ===> case ui <= math.MaxUint64:
			*result = append(*result, byte(MsgPackTypes.Uint32))
			binary.Write(&buf, binary.BigEndian, uint64(ui))
			*result = append(*result, buf.Bytes()...)
		}
	} else if i, err := num.Int64(); err == nil {

		switch {
		case i >= math.MinInt8 && i <= math.MaxInt8:
			*result = append(*result, byte(MsgPackTypes.Int8))
			*result = append(*result, byte(int8(i)))

		case i >= math.MinInt16 && i <= math.MaxInt16:
			*result = append(*result, byte(MsgPackTypes.Int16))
			binary.Write(&buf, binary.BigEndian, int16(i))
			*result = append(*result, buf.Bytes()...)

		case i >= math.MinInt32 && i <= math.MaxInt32:
			*result = append(*result, byte(MsgPackTypes.Int32))
			binary.Write(&buf, binary.BigEndian, int32(i))
			*result = append(*result, buf.Bytes()...)

		default: // ===> case i >= math.MinInt64 && i <= math.MaxInt64:
			*result = append(*result, byte(MsgPackTypes.Int64))
			binary.Write(&buf, binary.BigEndian, int64(i))
			*result = append(*result, buf.Bytes()...)
		}
	} else if f, err := num.Float64(); err == nil {
		switch {
		// case f >= math.SmallestNonzeroFloat32 && f <= math.MaxFloat32:
		// 	*result = append(*result, byte(MsgPackTypes.Float32))
		// 	binary.Write(&buf, binary.BigEndian, float32(f))
		// 	*result = append(*result, buf.Bytes()...)

		default: // ===> case f >= math.SmallestNonzeroFloat64 && f <= math.MaxFloat64:
			*result = append(*result, byte(MsgPackTypes.Float64))
			binary.Write(&buf, binary.BigEndian, float64(f))
			*result = append(*result, buf.Bytes()...)
		}
	}

	return result, nil
}

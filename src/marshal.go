package msgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

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

func (m *Msgpack) handleValue(result *[]byte, data interface{}) {
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
		valueAsBytes := []byte(v.String())

		// add type prefix
		*result = append(*result, byte(len(valueAsBytes))+MsgPackTypes.FixStr)

		// add value
		*result = append(*result, valueAsBytes...)

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
}

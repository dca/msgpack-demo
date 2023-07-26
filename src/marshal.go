package msgpack

import (
	"fmt"
	"reflect"
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
	case reflect.Int:
		*result = append(*result, byte(v.Int()))
	case reflect.String:
		valueAsBytes := []byte(v.String())

		length := len(valueAsBytes)
		strPrefix := 0xa0 + length
		*result = append(*result, byte(strPrefix))

		*result = append(*result, valueAsBytes...)
		fmt.Println("string:", v)

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			m.handleValue(result, elem)
		}

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String {
			keys := v.MapKeys()

			mapPrefix := 0x80 + len(keys)
			*result = append(*result, byte(mapPrefix))

			for _, key := range keys {
				keyAsBytes := []byte(key.String())

				//
				length := len(keyAsBytes)
				strPrefix := 0xa0 + length
				*result = append(*result, byte(strPrefix))

				*result = append(*result, keyAsBytes...)

				elem := v.MapIndex(key).Interface()
				m.handleValue(result, elem)
			}
		} else {
			fmt.Println("Non-string keys in map are not supported.")
		}
	}
}

package MsgPackTypes

// IsMsgPackTypeNil checks if the byte represents a MsgPack nil type.
func IsMsgPackTypeNil(b byte) bool {
	return b == Nil
}

// IsMsgPackTypeTrue checks if the byte represents a MsgPack true type.
func IsMsgPackTypeTrue(b byte) bool {
	return b == True
}

// IsMsgPackTypeFalse checks if the byte represents a MsgPack false type.
func IsMsgPackTypeFalse(b byte) bool {
	return b == False
}

func IsMsgPackTypePositiveInt(b byte) bool {
	return (b & 0x80) == FixIntPos
}

func IsMsgPackTypeNegativeInt(b byte) bool {
	return (b & 0xE0) == FixIntNeg
}

func IsMsgPackTypeUint8(b byte) bool {
	return b == Uint8
}

func IsMsgPackTypeUint16(b byte) bool {
	return b == Uint16
}

func IsMsgPackTypeUint32(b byte) bool {
	return b == Uint32
}

func IsMsgPackTypeUInt64(b byte) bool {
	return b == Uint64
}

func IsMsgPackTypeInt8(b byte) bool {
	return b == Int8
}

func IsMsgPackTypeInt16(b byte) bool {
	return b == Int16
}

func IsMsgPackTypeInt32(b byte) bool {
	return b == Int32
}

func IsMsgPackTypeInt64(b byte) bool {
	return b == Int64
}

// IsMsgPackTypeFloat32 checks if the byte represents a MsgPack float32 type.
func IsMsgPackTypeFloat32(b byte) bool {
	return b == Float32
}

// IsMsgPackTypeFloat64 checks if the byte represents a MsgPack float64 type.
func IsMsgPackTypeFloat64(b byte) bool {
	return b == Float64
}

// IsMsgPackTypeString checks if the byte represents a MsgPack string type.
func IsMsgPackTypeString(b byte) bool {
	// Using binary mask to filter out the 3 MSBs and compare with the pattern 101xxxxx
	return (b & 0xE0) == FixStr
}

// IsMsgPackTypeArray checks if the byte represents a MsgPack array type.
func IsMsgPackTypeArray(b byte) bool {
	// Using binary mask to filter out the 4 MSBs and compare with the pattern 1001xxxx
	return (b & 0xF0) == FixArray
}

// IsMsgPackTypeMap checks if the byte represents a MsgPack map type.
func IsMsgPackTypeMap(b byte) bool {
	// Using binary mask to filter out the 4 MSBs and compare with the pattern 1000xxxx
	return (b & 0xF0) == FixMap
}

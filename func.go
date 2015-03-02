package flv

func bytes_int24(x []byte) int {
	return int(x[0])<<16 + int(x[1])<<8 + int(x[2])
}

func int24_bytes(x int) []byte {
	return []byte{byte(x >> 16), byte(x >> 8), byte(x)}
}

func bytes_int32(x []byte) int {
	return int(x[0])<<24 + int(x[1])<<16 + int(x[2])<<8 + int(x[3])
}

func int32_bytes(x int) []byte {
	return []byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}
}

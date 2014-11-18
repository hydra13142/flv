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

type Tag struct {
	Head []byte
	Data []byte
}

func (this *Tag) Audio() bool {
	return this.Head[0]&31 == 8
}

func (this *Tag) Video() bool {
	return this.Head[0]&31 == 9
}

func (this *Tag) Script() bool {
	return this.Head[0]&31 == 18
}

func (this *Tag) Size() int {
	return bytes_int24(this.Head[1:])
}

func (this *Tag) SetSize(sz int) {
	copy(this.Head[1:], int24_bytes(sz))
}

func (this *Tag) Time() int {
	s := bytes_int24(this.Head[4:])
	if h := this.Head[7]; h != 0 {
		s += int(h) << 24
	}
	return s
}

func (this *Tag) SetTime(ms int) {
	copy(this.Head[4:], int24_bytes(ms))
	if h := ms >> 24; h != 0 {
		this.Head[7] = byte(h)
	}
}

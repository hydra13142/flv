package flv

import "github.com/hydra13142/amf"

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
	return amf.Int24(this.Head[1:])
}

func (this *Tag) SetSize(sz int) {
	copy(this.Head[1:], amf.ToInt24(sz))
}

func (this *Tag) Time() int {
	s := amf.Int24(this.Head[4:])
	if h := this.Head[7]; h != 0 {
		s += int(h) << 24
	}
	return s
}

func (this *Tag) SetTime(ms int) {
	copy(this.Head[4:], amf.ToInt24(ms))
	if h := ms >> 24; h != 0 {
		this.Head[7] = byte(h)
	}
}

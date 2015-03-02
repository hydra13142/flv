package flv

import "strconv"

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

func (this *Tag) Keyframe() bool {
	a, b := this.Head[0], this.Data[0]
	return (a>>4 == 1) || (a == 9 && b == 23)
}

func (this *Tag) VideoType() string {
	var (
		a, b string
		r    = this.Head[0]
	)
	switch r >> 4 {
	case 1:
		a = "keyframe"
	case 2:
		a = "inter frame"
	case 3:
		a = "disposable inter frame"
	case 4:
		a = "generated keyframe"
	case 5:
		a = "video info/command frame"
	default:
		a = "unknown frame " + strconv.Itoa(int(r>>4))
	}
	switch r & 15 {
	case 1:
		b = "JPEG"
	case 2:
		b = "Sorenson H.263"
	case 3:
		b = "Screen video"
	case 4:
		b = "On2 VP6"
	case 5:
		b = "On2 VP6 with alpha channel"
	case 6:
		b = "Screen video version 2"
	case 7:
		b = "AVC"
	default:
		b = "unknown code " + strconv.Itoa(int(r&15))
	}
	return a + ", " + b

}

func (this *Tag) AudioType() string {
	var (
		a, b, c, d string
		r          = this.Head[0]
	)
	switch r >> 4 {
	case 0:
		a = "Linear PCM, platform endian"
	case 1:
		a = "ADPCM"
	case 2:
		a = "MP3"
	case 3:
		a = "Linear PCM little endian"
	case 4:
		a = "Nellymoser 16-kHz mono"
	case 5:
		a = "Nellymoser 8-kHz mono"
	case 6:
		a = "Nellymoser"
	case 7:
		a = "G.711 A-law logarithmic PCM"
	case 8:
		a = "G.711 mu-law logarithmic PCM"
	case 9:
		a = "reserved"
	case 10:
		a = "AAC"
	case 11:
		a = "Speex"
	case 14:
		a = "MP3 8-Khz"
	case 15:
		a = "Device-specific sound"
	}
	switch (r >> 2) & 3 {
	case 0:
		b = "5.5 kHz"
	case 1:
		b = "11 kHz"
	case 2:
		b = "22 kHz"
	case 3:
		b = "44 kHz"
	}
	if r&2 == 0 {
		c = "8 bits"
	} else {
		c = "16 bits"
	}
	if r&2 == 0 {
		d = "mono"
	} else {
		d = "stereo"
	}
	return a + ", " + b + ", " + c + ", " + d
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

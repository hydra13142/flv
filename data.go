package flv

func Keyframe(r byte) bool {
	return (r >> 4) == 1
}

func VideoType(r byte) string {
	var a, b string
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
		a = "unknown frame"
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
		b = "unknown codetype"
	}
	return a + ", " + b

}

func AudioType(r byte) string {
	var a, b, c, d string
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
	switch r & 2 {
	case 0:
		c = "8 bits"
	case 2:
		c = "16 bits"
	}
	switch r & 1 {
	case 0:
		d = "mono"
	case 1:
		d = "stereo"
	}
	return a + ", " + b + ", " + c + ", " + d
}

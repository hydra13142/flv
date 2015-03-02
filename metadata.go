package flv

type MetaData struct {
	HasAudio     bool `amf:"hasAudio"`
	HasVideo     bool `amf:"hasVideo"`
	HasMetaData  bool `amf:"hasMetadata"`
	CanSeekToEnd bool `amf:"canSeekToEnd"`
	Stereo       bool `amf:"stereo"`

	Duration  float64 `amf:"duration"`
	FileSize  float64 `amf:"filesize"`
	DataSize  float64 `amf:"datasize"`
	AudioSize float64 `amf:"audiosize"`
	VideoSize float64 `amf:"videosize"`

	VideoCodecID float64 `amf:"videocodecid"`
	AudioCodecID float64 `amf:"audiocodecid"`

	AudioDataRate   float64 `amf:"audiodatarate,omitempty"`
	AudioSampleSize float64 `amf:"audiosamplesize,omitempty"`
	AudioSampleRate float64 `amf:"audiosamplerate,omitempty"`

	VideoDataRate         float64 `amf:"videodatarate,omitempty"`
	FrameRate             float64 `amf:"framerate,omitempty"`
	Width                 float64 `amf:"width,omitempty"`
	Height                float64 `amf:"height,omitempty"`
	HasKeyframes          bool    `amf:"hasKeyframes"`
	Keyframes             Frames  `amf:"keyframes,omitempty"`
	LastTimeStamp         float64 `amf:"lasttimestamp,omitempty"`
	LastKeyframeTimeStamp float64 `amf:"lastkeyframetimestamp,omitempty"`
	LastKeyframeLocation  float64 `amf:"lastkeyframelocation,omitempty"`

	MetaDataCreator string `amf:"metadatacreator,omitempty"`
}

type Frames struct {
	Times         []float64 `amf:"times"`
	FilePositions []float64 `amf:"filepositions"`
}

func SameCodec(a, b *MetaData) bool {
	if a.HasAudio != b.HasAudio {
		return false
	}
	if a.HasAudio {
		if a.AudioCodecID != b.AudioCodecID {
			return false
		}
		if a.AudioDataRate != b.AudioDataRate {
			return false
		}
		if a.AudioSampleSize != b.AudioSampleSize {
			return false
		}
		if a.AudioSampleRate != b.AudioSampleRate {
			return false
		}
	}
	if a.HasVideo != b.HasVideo {
		return false
	}
	if a.HasVideo {
		if a.VideoCodecID != b.VideoCodecID {
			return false
		}
		if a.VideoDataRate != b.VideoDataRate {
			return false
		}
		if a.FrameRate != b.FrameRate {
			return false
		}
		if a.Width != b.Width || a.Height != b.Height {
			return false
		}
	}
	return true
}

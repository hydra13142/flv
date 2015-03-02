package flv

import (
	"bytes"
	"github.com/hydra13142/encoding/AMF"
)

type feed int

func (this *feed) Write(data []byte) (int, error) {
	l := len(data)
	*this += feed(l)
	return l, nil
}

func (this *FLV) Refresh() {
	A, V, D, F := 0, 0, 0, 0
	L := len(this.Tags)
	T := this.Tags[L-1].Time() - this.Tags[0].Time()
	KT := []float64{}
	KP := []float64{}
	for _, tag := range this.Tags[1:] {
		S := tag.Size()
		if tag.Audio() {
			A += S
		}
		if tag.Video() {
			V += S
			if tag.Keyframe() {
				KT = append(KT, float64(tag.Time())/1000)
				KP = append(KP, float64(F))
			}
		}
		D += S
		F += S + 15
	}
	this.Meta.Duration = float64(T) / 1000
	this.Meta.DataSize = float64(D)
	this.Meta.AudioSize = float64(A)
	this.Meta.VideoSize = float64(V)
	this.Meta.LastTimeStamp = this.Meta.Duration
	if len(KT) > 0 {
		this.Meta.HasKeyframes = true
		this.Meta.Keyframes = Frames{KT, KP}
		this.Meta.LastKeyframeTimeStamp = 99 // avoid being omited
		this.Meta.LastKeyframeLocation = 999 // avoid being omited
	} else {
		this.Meta.HasKeyframes = false
		this.Meta.Keyframes = Frames{}
	}
	p := feed(0)
	e := AMF.NewEncoder(&p)
	e.Encode(this.Meta, "AMF0")
	header := len(this.Head) + 4
	offset := header + (11 + (13 + int(p) + 3) + 4)
	this.Meta.FileSize = float64(F + offset)
	if len(KP) > 0 {
		for i, _ := range KP {
			KP[i] += float64(offset)
		}
		this.Meta.LastKeyframeTimeStamp = KT[len(KT)-1]
		this.Meta.LastKeyframeLocation = KP[len(KP)-1]
	} else {
		this.Meta.LastKeyframeTimeStamp = 0
		this.Meta.LastKeyframeLocation = 0
	}
	q := bytes.NewBuffer(nil)
	f := AMF.NewEncoder(q)
	f.Encode("onMetaData", "AMF0")
	f.Encode(this.Meta, "AMF0")
	q.Write([]byte{0, 0, 9})
	this.Tags[0].SetSize(q.Len())
	this.Tags[0].Data = q.Bytes()
}

func (this *FLV) Clip(bgn, end float64) error {
	var (
		i, j, k = 0, 0, 0
		L, R    = 0, 0
	)
	if !this.Meta.HasKeyframes {
		return nokeyframe
	}
	if bgn < 0 {
		bgn = 0
	}
	kf := this.Meta.Keyframes
	lk := len(kf.Times)
	for ; i < lk; i++ {
		if bgn < kf.Times[i] {
			break
		}
	}
	bgn = kf.Times[i-1]
	L = int(kf.FilePositions[i-1])
	for ; i < lk; i++ {
		if end <= kf.Times[i] {
			break
		}
	}
	if i < lk {
		R = int(kf.FilePositions[i])
	} else {
		R = int(this.Meta.FileSize)
	}

	lt := len(this.Tags)
	for i = 0; i < lt; i++ {
		if !this.Tags[i].Script() && this.Tags[i].Data[1] != 0 {
			break
		}
	}

	s := len(this.Head) + 4
	for k = 0; k < lt; k++ {
		if s == L {
			j = k
		}
		if s >= R {
			break
		}
		s += this.Tags[k].Size() + 15
	}
	copy(this.Tags[i:], this.Tags[j:k])
	this.Tags = this.Tags[:k-j+i]
	for _, tag := range this.Tags[i:] {
		tag.SetTime(tag.Time() - int(bgn*1000))
	}
	this.Refresh()
	return nil
}

func (this *FLV) Append(ano *FLV) error {
	if !SameCodec(&this.Meta, &ano.Meta) {
		return unsuitcode
	}
	offset := int(1000*this.Meta.Duration + 1000/this.Meta.FrameRate)
	a, b := len(this.Tags), len(ano.Tags)
	tag := make([]Tag, a+b-1)
	copy(tag, this.Tags)
	for i := 1; i < b; i++ {
		head := make([]byte, len(ano.Tags[i].Head))
		copy(head, ano.Tags[i].Head)
		tag[a+i-1] = Tag{head, ano.Tags[i].Data}
	}
	for i := 0; i < b-1; i++ {
		tag[a+i].SetTime(tag[a+i].Time() + offset)
	}
	this.Tags = tag
	this.Refresh()
	return nil
}

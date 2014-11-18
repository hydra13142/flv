package main

import (
	"bytes"
	"github.com/hydra13142/encoding/AMF"
	"github.com/hydra13142/flv"
	"log"
	"os"
)

func LogError(x ...interface{}) {
	if x[0] != nil {
		log.Fatalln(x...)
	}
}

func MergeFLV(list []string, file string) {
	var (
		meta  flv.MetaData
		video []flv.Tag
		audio []flv.Tag
	)

	// 读取所有flv文件
	flvs := make([]*flv.FLV, len(list))
	for i := 0; i < len(flvs); i++ {
		flvs[i] = flv.New()
	}
	for i, name := range list {
		r, err := os.Open(name)
		LogError(err)
		err = flvs[i].ReadFrom(r)
		LogError(err)
		r.Close()
	}

	// 读取第一个文件的元数据
	d := AMF.NewDecoder(bytes.NewReader(flvs[0].Tags[0].Data))
	LogError(d.Decode(nil))
	LogError(d.Decode(&meta))

	// 将所有文件的tag分类并顺序连接
	step := int(1000 / meta.FrameRate)
	move := 0
	for _, one := range flvs {
		if i := len(video) - 1; i > 0 {
			for j := 1; j < len(one.Tags); j++ {
				if one.Tags[j].Video() {
					move = video[i].Time() + step - one.Tags[j].Time()
					break
				}
			}
		}
		for _, tag := range one.Tags[1:] {
			tag.SetTime(tag.Time() + move)
			switch {
			case tag.Video():
				video = append(video, tag)
			case tag.Audio():
				audio = append(audio, tag)
			}
		}
	}

	// 对tag进行修复，去掉前黑和后黑
	var begin, over int
	if len(video) > 2 {
		a, b := video[0].Time(), video[1].Time()
		if b-a > step*3 {
			video = video[1:]
			begin = b
		} else {
			begin = a
		}
	}
	if i := len(video) - 1; i > 1 {
		a, b := video[i-1].Time(), video[i].Time()
		if b-a > step*3 {
			video = video[:i]
			over = a
		} else {
			over = b
		}
	}
	for {
		l := len(audio)
		if l <= 0 || audio[0].Time() >= begin {
			break
		}
		audio = audio[1:]
	}
	for {
		l := len(audio)
		if l <= 0 || audio[l-1].Time() <= over {
			break
		}
		audio = audio[:l-1]
	}

	// 将tag按时间顺序集中到一起
	whole := make([]flv.Tag, len(audio)+len(video)+1)
	times := []float64{}
	place := []float64{}
	A, V, D, F := 0, 0, 0, 0
	i, j, t := 0, 0, 1
	for {
		var a, v int
		if i >= len(audio) {
			a = 1 << 30
		} else {
			a = audio[i].Time()
		}
		if j >= len(video) {
			v = 1 << 30
		} else {
			v = video[j].Time()
		}
		if a == v && v == 1<<30 {
			break
		}
		if a <= v {
			whole[t] = audio[i]
			A += audio[i].Size()
			i++
		} else {
			if flv.Keyframe(video[j].Data[0]) {
				times = append(times, float64(video[j].Time())/1000)
				place = append(place, float64(F))
			}
			whole[t] = video[j]
			V += video[j].Size()
			j++
		}
		whole[t].SetTime(whole[t].Time() - begin)
		D += whole[t].Size()
		F += whole[t].Size() + 15
		t++
	}

	// 更新元数据
	meta.Duration = float64(over-begin) / 1000
	meta.DataSize = float64(D)
	meta.AudioSize = float64(A)
	meta.VideoSize = float64(V)
	meta.HasKeyframes = true
	meta.Keyframes = flv.Frames{times, place}
	meta.LastTimeStamp = meta.Duration
	meta.LastKeyframeTimeStamp = times[len(times)-1]
	p := bytes.NewBuffer(nil)
	e := AMF.NewEncoder(p)
	e.Encode(meta, "AMF0")
	header := len(flvs[0].Head) + 4
	offset := 11 + 13 + p.Len() + 3 + 4
	for i, v := range place {
		place[i] = v + float64(offset)
	}
	meta.FileSize = float64(F + header + offset)
	meta.LastKeyframeLocation = place[len(place)-1]

	// 创建元数据的标签
	p.Reset()
	e.Encode("onMetaData", "AMF0")
	e.Encode(meta, "AMF0")
	whole[0].Head = flvs[0].Tags[0].Head
	whole[0].Data = append(p.Bytes(), 0, 0, 9)
	whole[0].SetSize(p.Len() + 3)

	// 写入flv文件
	all := flv.FLV{flvs[0].Head, whole}
	w, err := os.Create(file)
	LogError(err)
	all.WriteTo(w)
	w.Close()
}

func main() {
	l := len(os.Args)
	if l <= 1 {
		log.Fatal("usage : executable [src...] dest")
	}
	MergeFLV(os.Args[1:], "./merge.flv")
}

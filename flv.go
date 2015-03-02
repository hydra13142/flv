package flv

import (
	"bytes"
	"errors"
	"github.com/hydra13142/encoding/AMF"
	"io"
)

var (
	uncomplete = errors.New("Uncomplete File")
	wrongcodec = errors.New("Wrong Codec")
	findnometa = errors.New("No Media Information")
	nokeyframe = errors.New("No Keyframe")
	unsuitcode = errors.New("Different Codec")
)

type FLV struct {
	Head []byte
	Meta MetaData
	Tags []Tag
}

func New() *FLV {
	return &FLV{[]byte{}, MetaData{}, []Tag{}}
}

func (this *FLV) Audio() bool {
	return len(this.Head) >= 9 && this.Head[4]&4 != 0
}

func (this *FLV) Video() bool {
	return len(this.Head) >= 9 && this.Head[4]&1 != 0
}

func (this *FLV) readHead(r io.Reader) (err error) {
	var (
		meta [9]byte
		data []byte
	)
	if _, err = r.Read(meta[:]); err != nil {
		return uncomplete
	}
	if string(meta[:4]) != "FLV\x01" {
		return wrongcodec
	}
	l := bytes_int32(meta[5:])
	data = make([]byte, l+4)
	if _, err = r.Read(data[9:]); err != nil {
		return err
	}
	copy(data, meta[:])
	this.Head = data[:l]
	return nil
}

func (this *FLV) readTag(r io.Reader) (err error) {
	var (
		meta []byte
		data []byte
	)
	meta = make([]byte, 11)
	if _, err = r.Read(meta); err != nil {
		return err
	}
	switch meta[0] & 31 {
	case 8, 9, 18:
		l := bytes_int24(meta[1:])
		data = make([]byte, l+4)
		if i, _ := r.Read(data); i != int(l+4) {
			return uncomplete
		}
		this.Tags = append(this.Tags, Tag{meta, data[:l]})
	default:
		return wrongcodec
	}
	return nil
}

func (this *FLV) ReadFrom(r io.Reader) (err error) {
	if err = this.readHead(r); err != nil {
		return err
	}
	if err = this.readTag(r); err != nil {
		return err
	}
	if string(this.Tags[0].Data[3:13]) != "onMetaData" {
		return findnometa
	}
	d := AMF.NewDecoder(bytes.NewReader(this.Tags[0].Data))
	if err = d.Decode(nil); err != nil {
		return err
	}
	if err = d.Decode(&this.Meta); err != nil {
		return err
	}
	for ; err == nil; err = this.readTag(r) {
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func (this *FLV) WriteTo(w io.Writer) (err error) {
	if _, err = w.Write(this.Head); err != nil {
		return err
	}
	if _, err = w.Write([]byte{0, 0, 0, 0}); err != nil {
		return err
	}
	for _, tag := range this.Tags {
		if _, err = w.Write(tag.Head); err != nil {
			return err
		}
		if _, err = w.Write(tag.Data); err != nil {
			return err
		}
		if _, err = w.Write(int32_bytes(int(tag.Size() + 11))); err != nil {
			return err
		}
	}
	return nil
}

package flv

import (
	"github.com/hydra13142/amf"
	"bytes"
	"fmt"
	"io/ioutil"
	"io"
)

type Flv struct {
	Head []byte
	Tags []Tag
}

func New() *Flv {
	return &Flv{[]byte{}, []Tag{}}
}

func (this *Flv) Audio() bool {
	return len(this.Head) >= 9 && this.Head[4]&4 != 0
}

func (this *Flv) Video() bool {
	return len(this.Head) >= 9 && this.Head[4]&1 != 0
}

func (this *Flv) ReadHeadFrom(r io.Reader) error {
	data := make([]byte, 9)
	i, err := r.Read(data)
	if i != 9 || err != nil {
		return fmt.Errorf("Flv Header uncomplete")
	}
	if data[0] != 'F' || data[1] != 'L' || data[2] != 'V' || data[3] != 1 {
		return fmt.Errorf("Not a Flv Header")
	}
	l := amf.Int32(data[5:]) - 5
	rest := make([]byte, l)
	i, err = r.Read(rest)
	if i != int(l) {
		return fmt.Errorf("Flv Header uncomplete / without previous tag size #0")
	}
	if err != nil {
		return fmt.Errorf("Without Flv Tags")
	}
	this.Head = make([]byte, 5+l)
	copy(this.Head, data)
	copy(this.Head[9:], rest[:l-4])
	return nil
}

func (this *Flv) ReadTagFrom(r io.Reader) error {
	data := make([]byte, 11)
	i, err := r.Read(data)
	if i != 11 {
		if err == io.EOF {
			return err
		}
		return fmt.Errorf("Flv Tag Header uncomplete")
	}
	l := amf.Int24(data[1:])
	body := make([]byte, l+4)
	i, err = r.Read(body)
	if i != int(l+4) {
		return fmt.Errorf("Flv Tag Body uncomplete")
	}
	switch data[0] & 31 {
	case 8, 9, 18:
		this.Tags = append(this.Tags, Tag{data, body[:l]})
	default:
		return fmt.Errorf("Unrecognised Tag type")
	}
	return err
}

func (this *Flv) ReadFromFile(r io.Reader) (err error) {
	err = this.ReadHeadFrom(r)
	if err != nil {
		return err
	}
	err = this.ReadTagFrom(r)
	if err != nil {
		return err
	}
	if string(this.Tags[0].Data[3:13]) != "onMetaData" {
		return fmt.Errorf("without media information")
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	r = bytes.NewReader(data)
	for err != io.EOF {
		err = this.ReadTagFrom(r)
		if err != nil && err != io.EOF {
			return err
		}
	}
	defer func() {
		if pnc := recover(); pnc != nil {
			err = fmt.Errorf("%v", pnc)
		}
	}()
	return nil
}

func (this *Flv) WriteTo(w io.Writer) error {
	_, err := w.Write(this.Head)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0, 0, 0, 0})
	if err != nil {
		return err
	}
	for _, tag := range this.Tags {
		_, err = w.Write(tag.Head)
		if err != nil {
			return err
		}
		_, err = w.Write(tag.Data)
		if err != nil {
			return err
		}
		_, err = w.Write(amf.ToInt32(int32(tag.Size() + 11)))
		if err != nil {
			return err
		}
	}
	return nil
}

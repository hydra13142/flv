package main

import (
	"flag"
	"fmt"
	"github.com/hydra13142/flv"
	"os"
)

func main() {
	name := flag.String("o", "./merge.flv", "merged flv name")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.PrintDefaults()
		return
	}
	flvs := make([]*flv.FLV, flag.NArg())
	for i := 0; i < len(flvs); i++ {
		flvs[i] = flv.New()
	}
	for i, name := range flag.Args() {
		r, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = flvs[i].ReadFrom(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		r.Close()
	}
	for i := 1; i < len(flvs); i++ {
		err := flvs[0].Append(flvs[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	w, err := os.Create(*name)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = flvs[0].WriteTo(w)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Close()
}

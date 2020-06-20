package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"path"
	"replays"
)

var replayPath = flag.String("file", "", "Path to the Overwatch Replay MP4")

func init() {
	flag.Parse()
}
func main() {
	if *replayPath == "" || path.Ext(*replayPath) != ".mp4"{
		flag.PrintDefaults()
		return
	}


	f, _ := ioutil.ReadFile(*replayPath)
	// buf := bytes.NewBuffer(f)

	r, err := replays.Parse(f)
	if err != nil {
		panic(err)
	}

	spew.Dump(r)
}
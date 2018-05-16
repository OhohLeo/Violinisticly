package main

import (
	"engo.io/engo"

	"github.com/ohohleo/violin/jack/graphs"
)

const (
	DEFAULT_WIDTH  = 400
	DEFAULT_HEIGHT = 400
)

func main() {

	engo.Run(engo.RunOptions{
		Title:  "Graph",
		Width:  DEFAULT_WIDTH,
		Height: DEFAULT_HEIGHT,
	}, graphs.NewScene(DEFAULT_WIDTH, DEFAULT_HEIGHT))
}

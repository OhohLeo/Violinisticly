package graphs

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo/common"
)

type Scene struct {
	values    map[string]chan []float32
	fftValues map[string]chan []float32
	height    float32
	width     float32
	font      *common.Font
}

func NewScene(defaultFontPath string,
	values map[string]chan []float32,
	fftValues map[string]chan []float32,
	width, height float32) (scene *Scene, err error) {
	scene = &Scene{
		values:    values,
		fftValues: fftValues,
		height:    height,
		width:     width,
		font: &common.Font{
			URL:  "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			FG:   color.Black,
			Size: 64,
		},
	}

	err = scene.font.CreatePreloaded()
	return
}

func (*Scene) Type() string { return "Graph" }

func (*Scene) Preload() {}

func (s *Scene) Setup(world *ecs.World) {

	// Set white background
	common.SetBackground(color.White)

	renderSystem := &common.RenderSystem{}

	nbValues := len(s.values) + len(s.fftValues)*4

	var idx int
	height := s.height / float32(nbValues)

	for name, values := range s.values {
		graph := NewGraph(name, 0, height*float32(idx), s.width, height, 2,
			height*float32(idx)+height/2, 0.6, 512, s.font)
		go graph.GraphRaw(values)
		graph.AddToSystem(renderSystem)

		world.AddSystem(graph)

		idx++
	}

	for name, values := range s.fftValues {

		// Init graph
		graph := NewGraph(name, 0, height*float32(idx), s.width, height,
			2, height*float32(idx)+height, 0.6, 512, s.font)
		graph.AddToSystem(renderSystem)
		world.AddSystem(graph)

		// Init Spectrogram
		spectrogram := NewSpectrogram(name, 0, height*float32(idx+1), s.width, height)
		spectrogram.AddToSystem(renderSystem)
		world.AddSystem(spectrogram)

		// Display everything
		go func() {
			for {
				samples, ok := <-values
				if ok == false {
					return
				}

				graph.GraphFFT(samples)
				spectrogram.AddSamples(samples)
			}
		}()

		idx++
	}

	world.AddSystem(renderSystem)
}

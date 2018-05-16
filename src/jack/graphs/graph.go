package graphs

import (
	"fmt"
	"image/color"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Graph struct {
	Name   string
	Factor float32

	horizontalBar *Bar
	bars          []*Bar

	font  *common.Font
	texts []*Text

	lastIdx int

	offset float32
}

func NewGraph(name string, x, y, width, height, thickness, offset, maxValue float32, resolution int, font *common.Font) (g *Graph) {

	fmt.Printf("%s  x:%f y:%f\n", name, x, y)

	if thickness > 1 {
		offset -= thickness / 2
	}

	g = &Graph{
		Name:   name,
		Factor: height / maxValue,

		// Bar at the offset
		horizontalBar: NewBar(engo.Point{x, offset}, width, thickness, color.Black),

		font: font,

		offset: offset,
	}

	// Init bars
	barWidth := width / float32(resolution)

	fmt.Printf("width %f height %f barWidth %f\n", width, height, barWidth)

	for resolution > 0 {
		g.bars = append(g.bars, NewBar(engo.Point{x, offset},
			barWidth, 0, color.RGBA{255, 0, 0, 150}))
		x += barWidth
		resolution--
	}

	g.lastIdx = len(g.bars) - 1

	// Init texts
	// g.texts = append(g.texts, NewText(
	// 	engo.Point{0, 0}, width, height, g.font, color.RGBA{255, 0, 0, 150}))

	return
}

func (g *Graph) Update(dt float32) {}

func (*Graph) Remove(ecs.BasicEntity) {}

func (g *Graph) GraphRaw(values chan []float32) {

	for {
		samples, ok := <-values
		if ok == false {
			return
		}

		// Add all samples
		// for _, sample := range samples {
		// 	g.AddSample(sample * g.Factor)
		// }

		// Add only min/max
		var max, min float32
		for _, sample := range samples {
			if sample > max {
				max = sample
			}

			if sample < min {
				min = sample
			}
		}

		//fmt.Printf("%s: max %f\n", g.Name, max*g.Factor)

		g.Add(min*g.Factor, max*g.Factor)
	}
}

func (g *Graph) Add(min, max float32) {

	var position, height, previousPosition, previousHeight float32

	lastIdx := g.lastIdx
	position = g.offset - max
	height = max - min

	for lastIdx >= 0 {

		// Get from the last bar
		bar := g.bars[lastIdx]

		// Get previous value
		previousPosition, previousHeight = bar.Get()

		// Set new value
		bar.Set(position, height)

		// Replace new value by old one
		position = previousPosition
		height = previousHeight

		lastIdx--
	}
}

func (g *Graph) AddSample(sample float32) {

	var previous float32
	lastIdx := g.lastIdx

	for lastIdx >= 0 {

		// Get from the last bar
		bar := g.bars[lastIdx]

		// Get previous value
		previous = bar.GetValue()

		// Set new value
		bar.SetValue(sample)

		// Replace new value by old one
		sample = previous

		lastIdx--
	}
}

func (g *Graph) GraphFFT(samples []float32) {

	for idx, sample := range samples {

		if idx > g.lastIdx {
			continue
		}

		g.bars[idx].SetValue(sample)
	}

}

func (g *Graph) GraphRandom() {

	for _, bar := range g.bars {
		bar.Random()
	}
}

func (g *Graph) AddToSystem(system *common.RenderSystem) {

	if g.horizontalBar != nil {
		g.horizontalBar.AddToSystem(system)
	}

	for _, text := range g.texts {
		text.AddToSystem(system)
	}

	for _, bar := range g.bars {
		bar.AddToSystem(system)
	}

}

type Bar struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func NewBar(position engo.Point, width, height float32, c color.Color) (h *Bar) {

	h = &Bar{
		BasicEntity: ecs.NewBasic(),
	}

	h.SpaceComponent = common.SpaceComponent{Position: position, Width: width, Height: height}
	h.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: c}
	return
}

func (h *Bar) Get() (position, height float32) {
	position = h.SpaceComponent.Position.Y
	height = h.SpaceComponent.Height
	return
}

func (h *Bar) Set(position, height float32) {
	h.SpaceComponent.Position.Y = position
	h.SpaceComponent.Height = height
}

func (h *Bar) GetValue() float32 {
	return h.SpaceComponent.Height
}

func (h *Bar) SetValue(value float32) {
	h.SpaceComponent.Height = value
}

func (h *Bar) Random() {
	h.SpaceComponent.Height = -float32(rand.Intn(200))
}

func (h *Bar) AddToSystem(system *common.RenderSystem) {
	system.Add(&h.BasicEntity, &h.RenderComponent, &h.SpaceComponent)
}

type Text struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func NewText(position engo.Point, width, height float32, font *common.Font, c color.Color) (t *Text) {

	t = &Text{
		BasicEntity: ecs.NewBasic(),
	}

	t.SpaceComponent = common.SpaceComponent{Position: position, Width: width, Height: height}
	t.RenderComponent = common.RenderComponent{Drawable: common.Text{
		Font:          font,
		Text:          "Hello World",
		LineSpacing:   0.5,
		LetterSpacing: 0.15,
	}, Color: c}
	return
}

func (t *Text) Get() (position, height float32) {
	position = t.SpaceComponent.Position.Y
	height = t.SpaceComponent.Height
	return
}

func (t *Text) Set(position, height float32) {
	t.SpaceComponent.Position.Y = position
	t.SpaceComponent.Height = height
}

func (t *Text) AddToSystem(system *common.RenderSystem) {
	system.Add(&t.BasicEntity, &t.RenderComponent, &t.SpaceComponent)
}

package graphs

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/gl"
)

const (
	B_MASK = 255
	G_MASK = 255 << 8  //65280
	R_MASK = 255 << 16 //16711680
)

type Spectrogram struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	Name    string
	texture *SpectrogramTexture
}

func NewSpectrogram(name string, x, y, width, height float32) (s *Spectrogram) {

	fmt.Printf("%f %f %f %f\n", x, y, width, height)

	texture := NewSpectrogramTexture(int(width), int(height))

	s = &Spectrogram{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{x, y},
			Width:    width,
			Height:   height,
		},
		RenderComponent: common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{1, 1},
		},
		Name:    name,
		texture: texture,
	}

	return
}

func (s *Spectrogram) Update(dt float32) {
	s.texture.Update()
}

func (*Spectrogram) Remove(ecs.BasicEntity) {}

func (s *Spectrogram) AddToSystem(system *common.RenderSystem) {
	system.Add(&s.BasicEntity, &s.RenderComponent, &s.SpaceComponent)
}

func (s *Spectrogram) AddSamples(samples []float32) {

	colors := make([]color.NRGBA, len(samples))

	for idx, sample := range samples {

		sampleAbs := int(math.Abs(float64(sample)) * 2000)
		if sampleAbs > 16711680 {
			fmt.Printf("ABOVE %d\n", sampleAbs)
		}

		colors[idx] = color.NRGBA{
			uint8((sampleAbs & R_MASK) >> 16),
			uint8((sampleAbs & G_MASK) >> 8),
			uint8(sampleAbs & B_MASK),
			255,
		}
	}

	s.texture.AddColors(colors)
}

type SpectrogramTexture struct {
	texture  *gl.Texture
	image    *image.NRGBA
	object   *common.ImageObject
	width    float32
	height   float32
	viewport engo.AABB

	column int
}

func NewSpectrogramTexture(width, height int) (s *SpectrogramTexture) {

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	object := common.NewImageObject(img)

	return &SpectrogramTexture{
		image:    img,
		object:   object,
		width:    float32(width),
		height:   float32(height),
		texture:  common.UploadTexture(object),
		viewport: engo.AABB{Max: engo.Point{X: 1.0, Y: 1.0}},
	}
}

func (s *SpectrogramTexture) AddColors(colors []color.NRGBA) {

	if s.column >= int(s.width) {
		for i := 0; i < s.column; i++ {
			for j := 0; j < int(s.height); j++ {
				s.image.Set(i, j, s.image.At(i+1, j))
			}
		}

		s.column -= 1
	}

	for idx, color := range colors {
		s.image.Set(s.column, idx, color)
	}

	if s.column < int(s.width) {
		s.column++
	}

	s.object = common.NewImageObject(s.image)
}

func (s *SpectrogramTexture) Update() {

	if s.texture != nil {
		s.Close()
	}

	s.texture = common.UploadTexture(s.object)
}

func (s *SpectrogramTexture) Texture() *gl.Texture {
	return s.texture
}

func (s *SpectrogramTexture) Width() float32 {
	return s.width
}

func (s *SpectrogramTexture) Height() float32 {
	return s.height
}

func (s *SpectrogramTexture) View() (float32, float32, float32, float32) {
	return s.viewport.Min.X, s.viewport.Min.Y, s.viewport.Max.X, s.viewport.Max.Y
}

func (s *SpectrogramTexture) Close() {
	if !engo.Headless() {
		engo.Gl.DeleteTexture(s.texture)
	}
}

package opengl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

const windowWidth = 800
const windowHeight = 600

type Window struct {
	window       *glfw.Window
	previousTime time.Time
	fps          int
	objects      []*Object
}

func CreateWindow() (w *Window, err error) {

	// Initialisation de la fenêtre & OpenGL
	if err = glfw.Init(); err != nil {
		err = fmt.Errorf("failed to initialize glfw: %s", err)
		return
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Violin", nil, nil)
	if err != nil {
		return
	}

	w = &Window{
		window: window,
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(w.keyboard)

	// Initialisation de OpenGl
	if err = gl.Init(); err != nil {
		err = fmt.Errorf("Init OpenGl failed : %s", err)
		return
	}

	// Récupération de la version, du fabriquant & de la carte graphique
	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	return
}

func (w *Window) Start() {

	fps := 0
	counter := 0.0
	w.previousTime = time.Now()

	for !w.window.ShouldClose() {

		// Affichage du fps
		if time.Now().Unix() > w.previousTime.Unix() {
			w.previousTime = time.Now()
			fmt.Printf("\x1b[2K\x1b[G%d fps", fps)
			w.fps = fps
			fps = 0
		}

		// Nettoyage de l'écran
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Affichage des objects
		for _, object := range w.objects {
			object.Render()
		}

		// Rafraîchit la fenêtre
		w.window.SwapBuffers()
		glfw.PollEvents()

		counter += 0.01
		fps++
	}
}

func (w *Window) Stop() {
	glfw.Terminate()
}

func (w *Window) AddObject() *Object {

	object := CreateObject([]Vertex{
		{
			Position:     mgl32.Vec3{-0.5, -0.5, 0.0},
			TextureCoord: mgl32.Vec2{0.0, 0.0},
		},
		{
			Position:     mgl32.Vec3{0.0, 0.5, 0.0},
			TextureCoord: mgl32.Vec2{0.4, 1.0},
		},
		{
			Position:     mgl32.Vec3{0.5, -0.5, 0.0},
			TextureCoord: mgl32.Vec2{0.8, 0.0},
		},
	}, "./basicShader", "./bricks.jpg")

	w.objects = append(w.objects, object)

	return object
}

func (w *Window) keyboard(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifiers glfw.ModifierKey) {

	if len(w.objects) < 1 {
		return
	}

	transform := w.objects[0].GetTransform()

	switch key {
	case glfw.KeyUp:
		if modifiers&glfw.ModShift > 0 {
			transform.Scale(0.05, 0.05, 0.05)
		} else if modifiers&glfw.ModAlt > 0 {
			transform.Rotate(0.05, 0.0, 0.0)
		} else {
			transform.Move(0.0, 0.05, 0.0)
		}
	case glfw.KeyDown:
		if modifiers&glfw.ModShift > 0 {
			transform.Scale(-0.05, -0.05, -0.05)
		} else if modifiers&glfw.ModAlt > 0 {
			transform.Rotate(-0.05, 0.0, 0.0)
		} else {
			transform.Move(0.0, -0.05, 0.0)
		}
	case glfw.KeyLeft:
		if modifiers&glfw.ModControl > 0 {
			transform.Rotate(0.0, -0.05, 0.0)
		} else if modifiers&glfw.ModAlt > 0 {
			transform.Rotate(0.0, 0.0, -0.05)
		} else {
			transform.Move(-0.05, 0.0, 0.0)
		}
	case glfw.KeyRight:
		if modifiers&glfw.ModControl > 0 {
			transform.Rotate(0.0, 0.05, 0.0)
		} else if modifiers&glfw.ModAlt > 0 {
			transform.Rotate(0.0, 0.0, 0.05)
		} else {
			transform.Move(0.05, 0.0, 0.0)
		}
	}
}

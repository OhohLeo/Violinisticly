package opengl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

func CreateWindow() error {

	// Initialisation de la fenêtre & OpenGL
	window, err := initWindow()
	if err != nil {
		return err
	}
	defer window.Destroy()

	// Colorie le fond d'écran
	renderer, _ := sdl.CreateRenderer(window, -1, 0)
	renderer.SetDrawColor(20, 20, 20, 255)
	renderer.Clear()
	renderer.Present()

	vertices := []Vertex{
		{
			Position:     mgl32.Vec3{-0.5, -0.5, 0.0},
			TextureCoord: mgl32.Vec2{0.0, 0.0},
		},
		{
			Position:     mgl32.Vec3{0.0, 0.5, 0.0},
			TextureCoord: mgl32.Vec2{0.5, 1.0},
		},
		{
			Position:     mgl32.Vec3{0.5, -0.5, 0.0},
			TextureCoord: mgl32.Vec2{1.0, 0.0},
		},
	}

	mesh := CreateMesh(vertices)

	shader, err := CreateShader("./basicShader")
	if err != nil {
		panic(err)
	}

	texture, err := CreateTexture("./bricks.jpg")
	if err != nil {
		panic(err)
	}

	for {

		// Nettoyage de l'écran
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.Bind()
		texture.Bind(0)
		mesh.Draw()

		// Gestion des évèvements
		if handleEvents() {
			break
		}

		// Rafraîchit la fenêtre
		sdl.GL_SwapWindow(window)
	}

	sdl.Quit()

	return nil
}

func initWindow() (window *sdl.Window, err error) {

	// Initialisation de SDL
	sdl.Init(sdl.INIT_EVERYTHING)

	// Création de la fenêtre
	window, err = sdl.CreateWindow("Violin",
		sdl.WINDOWPOS_UNDEFINED, // Position x de la fenêtre
		sdl.WINDOWPOS_UNDEFINED, // Position y de la fenêtre
		800, // Largeur de la fenêtre
		600, // Hauteur de la fenêtre
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		err = fmt.Errorf("CreateWindow failed : %s", err)
		return
	}

	// Fixe la version d'OpenGl à utiliser
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)

	// Active le double buffuring avec OpenGl
	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GL_SetAttribute(sdl.GL_DEPTH_SIZE, 24)

	// Création du contexte graphique permettant d'openGl de fonctionner
	_, err = sdl.GL_CreateContext(window)
	if err != nil {
		panic(fmt.Errorf("CreateContext failed: %s", err))
	}

	// Initialisation de OpenGl
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("Init OpenGl failed : %s", err))
	}

	// Récupération de la version, du fabriquant & de la carte graphique
	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	return
}

func handleEvents() bool {

	// Récupération des évènements
	event := sdl.WaitEvent()

	// Gestion des évènements
	switch event.(type) {

	case *sdl.QuitEvent:
		return true
	}

	return false
}

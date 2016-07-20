package opengl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"unsafe"
)

const NUM_SHADERS = 2

type Shader struct {
	program uint32
	shaders []uint32
}

func CreateShader(filename string) (s *Shader, err error) {

	// Création du programme
	program := gl.CreateProgram()

	// Initialisation du tableau de shaders
	s = &Shader{
		program: program,
		shaders: make([]uint32, NUM_SHADERS),
	}

	// Charge le fichier contenant le VertexShader
	vertexShader, err := loadShader(filename + ".vs")
	if err != nil {
		return
	}

	// Charge le fichier contenant le FragmentShader
	fragmentShader, err := loadShader(filename + ".fs")
	if err != nil {
		return
	}

	s.shaders[0], err = createShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return
	}

	s.shaders[1], err = createShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return
	}

	for i := 0; i < NUM_SHADERS; i++ {
		gl.AttachShader(program, s.shaders[i])
	}

	name := []byte("position")
	gl.BindAttribLocation(program, 0, &name[0])

	// Vérification du programme
	gl.LinkProgram(program)
	if err = checkShader(program, gl.LINK_STATUS, true, "Program linking failed"); err != nil {
		return
	}

	gl.ValidateProgram(program)
	if err = checkShader(program, gl.VALIDATE_STATUS, true, "Program is invalid"); err != nil {
		return
	}

	return
}

func (s *Shader) Destroy() {

	for i := 0; i < NUM_SHADERS; i++ {
		gl.DetachShader(s.program, s.shaders[i])
		gl.DeleteShader(s.shaders[i])
	}

	gl.DeleteProgram(s.program)
}

func (s *Shader) Bind() {
	gl.UseProgram(s.program)
}

func createShader(data []byte, shaderType uint32) (shader uint32, err error) {

	shader = gl.CreateShader(shaderType)
	if shader == 0 {
		err = fmt.Errorf("createShader: Shader creation failed")
		return
	}

	cstrs, freeCb := gl.Strs(string(data))
	length := int32(len(data))
	gl.ShaderSource(shader, 1, cstrs, &length)
	freeCb()

	gl.CompileShader(shader)

	if err = checkShader(shader, gl.COMPILE_STATUS, false, "Shader compilation failed"); err != nil {
		return
	}

	return
}

func loadShader(filename string) (data []byte, err error) {

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("loadShader: %s", err.Error())
		return
	}

	return
}

func checkShader(shader uint32, flag uint32, isProgram bool, msg string) error {

	var success int32

	if isProgram {
		gl.GetProgramiv(shader, flag, &success)
	} else {
		gl.GetShaderiv(shader, flag, &success)
	}

	if success != gl.FALSE {
		return nil
	}

	errMsg := make([]byte, 1024)

	if isProgram {
		gl.GetProgramInfoLog(shader, int32(unsafe.Sizeof(errMsg)), nil, &errMsg[0])
	} else {
		gl.GetShaderInfoLog(shader, int32(unsafe.Sizeof(errMsg)), nil, &errMsg[0])
	}

	return fmt.Errorf("checkShader: %s [%s]", msg, errMsg)
}

package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type Vertex struct {
	Position mgl32.Vec3
}

const (
	NUM_BUFFERS = 1
)

type Mesh struct {
	vertexArrayObject  uint32
	vertexArrayBuffers uint32
	drawCount          int32
}

func CreateMesh(vertices []Vertex) *Mesh {
	numVertices := len(vertices)

	m := &Mesh{
		drawCount: int32(numVertices),
	}

	// Alloue le tableau de vertices
	gl.GenVertexArrays(1, &m.vertexArrayObject)
	gl.BindVertexArray(m.vertexArrayObject)

	// Alloue les buffers
	gl.GenBuffers(NUM_BUFFERS, &m.vertexArrayBuffers)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexArrayBuffers)

	// Récupère la taille du vertex
	vertexSize := int(unsafe.Sizeof(vertices[0]))

	// Transfère les données dans les buffers
	gl.BufferData(gl.ARRAY_BUFFER, numVertices*vertexSize,
		unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// Active la lecture des attributs
	gl.EnableVertexAttribArray(0)

	// Lis les attributs
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(0)

	return m
}

func (m *Mesh) Destroy() {
	gl.DeleteVertexArrays(1, &m.vertexArrayObject)
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vertexArrayObject)

	// Paramètres pour afficher l'object
	gl.DrawArrays(gl.TRIANGLES, 0, m.drawCount)

	gl.BindVertexArray(0)
}

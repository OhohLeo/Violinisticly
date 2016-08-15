package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type Vertex struct {
	Position     mgl32.Vec3
	TextureCoord mgl32.Vec2
}

const (
	POSITION_VB = iota
	TEXTURECOORD_VB
	NUM_BUFFERS
)

type Mesh struct {
	vertexArrayObject  uint32
	vertexArrayBuffers []uint32
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

	// Récupération des positions & des coordonnées des textures
	positions := make([]mgl32.Vec3, numVertices)
	textureCoords := make([]mgl32.Vec2, numVertices)

	for i := 0; i < numVertices; i++ {
		positions[i] = vertices[i].Position
		textureCoords[i] = vertices[i].TextureCoord
	}

	m.vertexArrayBuffers = make([]uint32, NUM_BUFFERS)

	// Création des buffers
	gl.GenBuffers(NUM_BUFFERS, &m.vertexArrayBuffers[POSITION_VB])

	// **********************
	// Gestion de la position

	// Alloue les buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexArrayBuffers[POSITION_VB])

	// Récupère la taille du vertex
	vertexSize := int(unsafe.Sizeof(positions[0]))

	// Transfère les données dans les buffers
	gl.BufferData(gl.ARRAY_BUFFER, numVertices*vertexSize,
		gl.Ptr(positions), gl.STATIC_DRAW)

	// Active la lecture des attributs
	gl.EnableVertexAttribArray(0)

	// Lis les attributs
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// *********************
	// Gestion de la texture

	// Alloue les buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexArrayBuffers[TEXTURECOORD_VB])

	// Récupère la taille du vertex
	vertexSize = int(unsafe.Sizeof(textureCoords[0]))

	// Transfère les données dans les buffers
	gl.BufferData(gl.ARRAY_BUFFER, numVertices*vertexSize,
		gl.Ptr(textureCoords), gl.STATIC_DRAW)

	// Active la lecture des attributs
	gl.EnableVertexAttribArray(1)

	// Lis les attributs
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

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

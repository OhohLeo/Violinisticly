package opengl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Quat
	scale    mgl32.Vec3
}

func SetCenter() *Transform {
	return &Transform{
		position: mgl32.Vec3{0.0, 0.0, 0.0},
		rotation: mgl32.QuatIdent(),
		scale:    mgl32.Vec3{1.0, 1.0, 1.0},
	}
}

func (t *Transform) GetModel() mgl32.Mat4 {
	positionMatrix := mgl32.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	rotateMatrix := t.rotation.Mat4()
	scaleMatrix := mgl32.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z())

	return positionMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))
}

func (t *Transform) Move(x float32, y float32, z float32) {
	t.position = t.position.Add(mgl32.Vec3{x, y, z})
}

func (t *Transform) Rotate(x float32, y float32, z float32) {
	// rotateX := mgl32.QuatRotate(x, mgl32.Vec3{0, 0, 1})
	// rotateY := mgl32.QuatRotate(y, mgl32.Vec3{0, 1, 0})
	// rotateZ := mgl32.QuatRotate(z, mgl32.Vec3{0, 0, 1})

	//t.rotation = t.rotation.Add(rotateX)
}

func (t *Transform) SetRotate(w float32, x float32, y float32, z float32) {
	t.rotation = mgl32.Quat{
		W: w,
		V: mgl32.Vec3{x, y, z},
	}
}

func (t *Transform) Scale(x float32, y float32, z float32) {
	t.scale = t.scale.Add(mgl32.Vec3{x, y, z})
}

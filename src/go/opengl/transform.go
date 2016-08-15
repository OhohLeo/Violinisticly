package opengl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    mgl32.Vec3
}

func SetCenter() *Transform {
	return &Transform{
		position: mgl32.Vec3{0.0, 0.0, 0.0},
		rotation: mgl32.Vec3{0.0, 0.0, 0.0},
		scale:    mgl32.Vec3{1.0, 1.0, 1.0},
	}
}

func (t *Transform) GetModel() mgl32.Mat4 {
	positionMatrix := mgl32.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	rotateXMatrix := mgl32.HomogRotate3D(t.rotation.X(), mgl32.Vec3{1.0, 0.0, 0.0})
	rotateYMatrix := mgl32.HomogRotate3D(t.rotation.Y(), mgl32.Vec3{0.0, 1.0, 0.0})
	rotateZMatrix := mgl32.HomogRotate3D(t.rotation.Z(), mgl32.Vec3{0.0, 0.0, 1.0})
	scaleMatrix := mgl32.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z())

	rotateMatrix := rotateZMatrix.Mul4(rotateYMatrix.Mul4(rotateXMatrix))

	return positionMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))
}

func (t *Transform) Move(x float32, y float32, z float32) {
	t.position = t.position.Add(mgl32.Vec3{x, y, z})
}

func (t *Transform) SetRotate(x float32, y float32, z float32) {
	t.rotation = mgl32.Vec3{x, y, z}
}

func (t *Transform) Rotate(x float32, y float32, z float32) {
	t.rotation = t.rotation.Add(mgl32.Vec3{x, y, z})
}

func (t *Transform) Scale(x float32, y float32, z float32) {
	t.scale = t.scale.Add(mgl32.Vec3{x, y, z})
}

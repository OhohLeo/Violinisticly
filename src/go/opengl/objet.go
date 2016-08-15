package opengl

type Object struct {
	transform *Transform
	mesh      *Mesh
	shader    *Shader
	texture   *Texture
}

func CreateObject(vertices []Vertex, shaderUrl string, textureUrl string) *Object {

	mesh := CreateMesh(vertices)

	shader, err := CreateShader(shaderUrl)
	if err != nil {
		panic(err)
	}

	texture, err := CreateTexture(textureUrl)
	if err != nil {
		panic(err)
	}

	return &Object{
		mesh:      mesh,
		shader:    shader,
		texture:   texture,
		transform: SetCenter(),
	}
}

func (o *Object) GetTransform() *Transform {
	return o.transform
}

func (o *Object) Render() {
	o.shader.Bind()
	o.shader.Update(o.transform)
	o.texture.Bind(0)
	o.mesh.Draw()
}

package opengl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Texture struct {
	id uint32
}

func CreateTexture(filename string) (t *Texture, err error) {

	// Récupération du fichier png, jpeg ou gif
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return
	}

	// Récupération de l'image à décoder
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}

	// Conversion de l'image au format RGBA
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		err = fmt.Errorf("CreateTexture: unsupported stride")
		return
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Création de la texture
	t = new(Texture)

	// Récupération de l'identifiant de la texture
	gl.GenTextures(1, &t.id)

	// Mise en place de la mémoire
	gl.BindTexture(gl.TEXTURE_2D, t.id)

	// Répétition en hauteur & largeur
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Comportement en cas de minification/agrandissement de la texture
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Envoie la texture au CPU
	gl.TexImage2D(gl.TEXTURE_2D,
		0,                         // level
		gl.RGBA,                   // internal format
		int32(rgba.Rect.Size().X), // width
		int32(rgba.Rect.Size().Y), // height
		0,       // level
		gl.RGBA, // format
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return
}

func (t *Texture) Destroy() {
	gl.DeleteTextures(1, &t.id)
}

func (t *Texture) Bind(unit uint32) error {
	if unit < 0 && unit > 31 {
		return fmt.Errorf("Bind Texture: invalid unit value '%d'", unit)
	}

	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, t.id)

	return nil
}

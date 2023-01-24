package animation

import (
	"image"

	_ "image/png"
)

type Animatable interface {
	Frame() image.Image
	Update()
	Create(filepath string)
}

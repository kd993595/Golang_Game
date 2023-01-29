package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LightSystem struct {
	lightmap *ebiten.Image
}

type GlowSystem struct {
	redGlow *ebiten.Image
}

func (g *GlowSystem) Initialize(rImg *ebiten.Image) {
	g.redGlow = rImg
}

func (l *LightSystem) Initialize(textureImg *ebiten.Image) {
	l.lightmap = textureImg
}

func (l *LightSystem) DrawLight(coord [2]int, lightImg *ebiten.Image, clr color.Color, radius float64) {
	w, h := lightImg.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	if clr != nil {
		r, g, b, a := clr.RGBA()
		op.ColorM.Translate(float64(r), float64(g), float64(b), float64(a))
	}
	op.GeoM.Translate(float64(coord[0]-(w/2*int(radius))), float64(coord[1]-h/2*int(radius)))
	op.CompositeMode = ebiten.CompositeModeLighter
	l.lightmap.DrawImage(lightImg, op)
}

func (l *LightSystem) GetLightMap() *ebiten.Image {
	return l.lightmap
}

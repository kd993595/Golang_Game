package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LightSystem struct {
	lightmap *ebiten.Image
}

func (l *LightSystem) Initialize(textureImg *ebiten.Image) {
	l.lightmap = textureImg
}

func (l *LightSystem) DrawLight(coord [2]float64, lightImg *ebiten.Image, clr color.Color, radius float64) {
	w, h := lightImg.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	if clr != nil {
		r, g, b, a := clr.RGBA()
		op.ColorM.Translate(float64(r), float64(g), float64(b), float64(a))
	}
	op.GeoM.Translate(coord[0]-float64(w/2)*radius, coord[1]-float64(h/2)*radius)
	op.CompositeMode = ebiten.CompositeModeLighter
	l.lightmap.DrawImage(lightImg, op)
}

func (l *LightSystem) GetLightMap() *ebiten.Image {
	return l.lightmap
}

type GlowSystem struct {
	glowmap *ebiten.Image
}

func (g *GlowSystem) Initialize(rImg *ebiten.Image) {
	g.glowmap = rImg
}

func (g *GlowSystem) DrawGlow(coord [2]float64, glowImg *ebiten.Image, clr color.Color, radius float64) {
	w, h := glowImg.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	if clr != nil {
		op.ColorM.Scale(0, 0, 0, 1)
		r, g, b, a := clr.RGBA()
		op.ColorM.Translate(float64(r), float64(g), float64(b), float64(a))

	}
	op.GeoM.Translate(coord[0]-float64(w/2)*radius, coord[1]-float64(h/2)*radius)
	op.CompositeMode = ebiten.CompositeModeLighter
	g.glowmap.DrawImage(glowImg, op)
}

func (g *GlowSystem) GetGlowMap() *ebiten.Image {
	return g.glowmap
}

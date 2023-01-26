package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	idle = iota
	running
)

type Player struct {
	Count  int
	PosX   int
	PosY   int
	FrameX int
	FrameY int
	Width  int
	Height int
	Img    *ebiten.Image
	Speed  int
	Idle   bool

	State int
}

func (p *Player) Frame() image.Image {
	p.FrameX = 0
	i := (p.Count / 10) % 4
	if p.Idle {
		p.FrameX = 112
	}
	p.FrameX += p.Width * i
	return p.Img.SubImage(image.Rect(p.FrameX, p.FrameY, p.FrameX+p.Width, p.Height))
}

func (p *Player) Update() {
	p.Idle = true
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.PosX += p.Speed
		p.Idle = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.PosY += p.Speed
		p.Idle = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.PosX -= p.Speed
		p.Idle = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.PosY -= p.Speed
		p.Idle = false
	}
	p.Count++
}

func (p *Player) Create(filepath string) {
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	w, h := origEbitenImage.Size()
	p.Img = ebiten.NewImage(w*2, h*2)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	p.Img.DrawImage(origEbitenImage, op)
}

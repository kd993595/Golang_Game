package player

import (
	"image"
)

//rectangle returned for image has values hard-coded into it so if player image were to change change values
type Player struct {
	Count      int
	Idle       bool
	FacingLeft bool
}

//160,0 for idle and 30*40 and 245,0 for running

func (p *Player) Frame() image.Rectangle {
	xPos := 169
	yPos := 0
	imgWidth := 18
	imgHeight := 39
	i := (p.Count / 10) % 6
	if !p.Idle {
		xPos = 250
	}
	yPos += imgHeight * i
	return image.Rect(xPos, yPos, xPos+imgWidth, yPos+imgHeight)
}

func (p *Player) Update() {
	p.Count++
	if p.Count == 60 {
		p.Count = 0
	}
}

/*func (p *Player) Create(filepath string) {
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
}*/

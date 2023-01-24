package card

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Deck struct {
	CardDeck [8]Card
}

type Card struct {
	Img  *ebiten.Image
	Text string
	Icon *ebiten.Image
}

func (c *Card) Create(filepath string) {
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	w, h := origEbitenImage.Size()
	c.Img = ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}

	c.Img.DrawImage(origEbitenImage, op)
}

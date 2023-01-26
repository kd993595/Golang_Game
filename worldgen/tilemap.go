package worldgen

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	grass = iota //304,272
	tree         //272,144
  
)

type Tilemap struct {
	GameMap [Width][Height]int
	Randgen *rand.Rand
}

func (t *Tilemap) Initialize(r *rand.Rand) {
	t.GameMap = [Width][Height]int{}
	t.Randgen = r
}

func (t *Tilemap) ProcessMap(m [Width][Height]int) {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if m[x][y] == 1 {
				t.GameMap[x][y] = tree
			} else {
				t.GameMap[x][y] = grass
			}
		}
	}
}

func (t *Tilemap) DrawWorld(e *ebiten.Image, tImg1 *ebiten.Image,tImg2 *ebiten.Image) {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))
			e.DrawImage(tImg1.SubImage(image.Rect(304, 272, 320, 288)).(*ebiten.Image), &op) //grass drawn all over map
      
			if t.GameMap[x][y] == tree {
				e.DrawImage(tImg1.SubImage(image.Rect(272, 144, 288, 160)).(*ebiten.Image), &op)
			}

      //probabilty gotta figure out for decorations
      probability := t.Randgen.Intn(100)
      if(probability<10){
        
      }
		}
	}
}

func isBoundary(x int, y int) bool {
	return x == 0 || x == Width-1 || y == 0 || y == Height-1
}

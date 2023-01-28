package components

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

//all images here are 16x16 pixels
const (
	grass                  = iota //304,272 tileset1
	tree                          //221,155 tileset1
	grassWithFlower1              //478,156 tileset2
	grassWithFlower2              //494,156 tileset2
	grassWithFlower3              //510,156 tileset2
	grassWithFlower4              //528,156 tileset2
	grassWithFlower5              //49,240 tileset1
	grassWithSmallGrass1          //0,80 tileset1
	grassWithSmallGrass2          //0,130 tileset1
	grassWithBrownMushroom        //816,38 tileset2
	grassWithRedMushroom          //816,73 tileset2
	grassWithSpotMushroom         //816,107 tileset2
	blank
	bigtree
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
				probability := t.Randgen.Intn(100)
				if probability < 20 {
					option := t.Randgen.Intn(2)
					switch option {
					case 0:
						t.GameMap[x][y] = grassWithSmallGrass1
					case 1:
						t.GameMap[x][y] = grassWithSmallGrass2
					default:
						t.GameMap[x][y] = grass
					}
				} else if probability < 30 {
					option := t.Randgen.Intn(5)
					switch option {
					case 0:
						t.GameMap[x][y] = grassWithFlower1
					case 1:
						t.GameMap[x][y] = grassWithFlower2
					case 2:
						t.GameMap[x][y] = grassWithFlower3
					case 3:
						t.GameMap[x][y] = grassWithFlower4
					case 4:
						t.GameMap[x][y] = grassWithFlower5
					default:
						t.GameMap[x][y] = grass
					}
				} else if probability < 40 {
					option := t.Randgen.Intn(3)
					switch option {
					case 0:
						t.GameMap[x][y] = grassWithRedMushroom
					case 1:
						t.GameMap[x][y] = grassWithBrownMushroom
					case 2:
						t.GameMap[x][y] = grassWithSpotMushroom
					default:
						t.GameMap[x][y] = grass
					}
				}
			}
		}
	}
}

func (t *Tilemap) DrawWorld(e *ebiten.Image, tImg1 *ebiten.Image, tImg2 *ebiten.Image) {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))
			e.DrawImage(tImg1.SubImage(image.Rect(304, 272, 320, 288)).(*ebiten.Image), &op) //grass drawn all over map

			if t.GameMap[x][y] == tree && !isBoundary(x, y) {
				if t.GameMap[x+1][y] == tree && t.GameMap[x][y+1] == tree && t.GameMap[x+1][y+1] == tree {
					t.GameMap[x+1][y] = blank
					t.GameMap[x+1][y+1] = blank
					t.GameMap[x][y+1] = blank
					t.GameMap[x][y] = bigtree
					continue
				}
				e.DrawImage(tImg1.SubImage(image.Rect(208, 129, 224, 145)).(*ebiten.Image), &op)

			}

		}
	}

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {

			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))
			if isBoundary(x, y) && x%2 == 0 && y%2 == 0 {
				e.DrawImage(tImg1.SubImage(image.Rect(96, 160, 128, 192)).(*ebiten.Image), &op)
				//e.DrawImage(tImg1.SubImage(image.Rect(0, 160, 32, 192)).(*ebiten.Image), &op)
				continue
			}

			if t.GameMap[x][y] == bigtree {
				e.DrawImage(tImg1.SubImage(image.Rect(96, 160, 128, 192)).(*ebiten.Image), &op)
			}

			randomoffset := [2]int{t.Randgen.Intn(5), t.Randgen.Intn(5)}
			op.GeoM.Reset()
			op.GeoM.Scale(.75, .75)
			op.GeoM.Translate(float64(x*16)+float64(randomoffset[0]), float64(y*16)+float64(randomoffset[1]))
			switch t.GameMap[x][y] {
			case grassWithFlower1:
				e.DrawImage(tImg2.SubImage(image.Rect(478, 156, 494, 172)).(*ebiten.Image), &op)
			case grassWithFlower2:
				e.DrawImage(tImg2.SubImage(image.Rect(494, 156, 510, 172)).(*ebiten.Image), &op)
			case grassWithFlower3:
				e.DrawImage(tImg2.SubImage(image.Rect(510, 156, 526, 172)).(*ebiten.Image), &op)
			case grassWithFlower4:
				e.DrawImage(tImg2.SubImage(image.Rect(528, 156, 544, 172)).(*ebiten.Image), &op)
			case grassWithFlower5:
				e.DrawImage(tImg1.SubImage(image.Rect(49, 240, 65, 256)).(*ebiten.Image), &op)
			case grassWithBrownMushroom:
				e.DrawImage(tImg2.SubImage(image.Rect(816, 38, 832, 54)).(*ebiten.Image), &op)
			case grassWithRedMushroom:
				e.DrawImage(tImg2.SubImage(image.Rect(816, 73, 832, 89)).(*ebiten.Image), &op)
			case grassWithSpotMushroom:
				e.DrawImage(tImg2.SubImage(image.Rect(816, 107, 832, 123)).(*ebiten.Image), &op)
			case grassWithSmallGrass1:
				e.DrawImage(tImg1.SubImage(image.Rect(0, 80, 16, 94)).(*ebiten.Image), &op)
			case grassWithSmallGrass2:
				e.DrawImage(tImg1.SubImage(image.Rect(0, 130, 16, 144)).(*ebiten.Image), &op)
			}

		}
	}
}

func isBoundary(x int, y int) bool {
	return x <= 3 || x >= Width-4 || y <= 3 || y >= Height-4
}

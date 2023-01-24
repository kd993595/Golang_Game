package worldgen

import "math/rand"

const (
	grass = iota
	lgrass
	rgrass
	ugrass
	dgrass
	water
)

type tilemap struct {
	GameMap [Width][Height]int
	Randgen *rand.Rand
}

func (t *tilemap) Initialize(r *rand.Rand) {
	t.GameMap = [Width][Height]int{}
	t.Randgen = r
}

func (t *tilemap) ProcessMap(m [Width][Height]int) {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if m[x][y] == 1 {
				t.GameMap[x][y] = grass
			} else {
				t.GameMap[x][y] = water
			}
		}
	}
}

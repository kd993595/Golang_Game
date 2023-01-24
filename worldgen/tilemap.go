package worldgen

import "math/rand"

const (
	grass = iota
	lwater
	rwater
	uwater
	dwater
  crdwater
  cldwater
  cluwater
  cruwater
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
        if isBoundary(x,y){
          continue
        }
        if m[x][y+1]==1{
          t.GameMap[x][y] = dwater
        }else if m[x][y-1] == 1{
          t.GameMap[x][y] = uwater
        }else if m[x+1][y] == 1{
          t.GameMap[x][y] = rwater
        }else if m[x-1][y] == 1{
          t.GameMap[x][y] = lwater
        }
			}
		}
	}
}

func isBoundary(x int,y int) bool{
  return x == 0 || x==Width-1 || y==0 || y==Height-1
}
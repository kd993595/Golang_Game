package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 16 //must match size of tile images in actuual game
)

type World struct {
	tilemap       []RectAABB
	BoundingBoxes []Quad
}

type Quad struct {
	x             int
	y             int
	width         int
	height        int
	StaticBodies  []RectAABB
	dynamicBodies []RigidBody
}

func (w *World) Initialize(m [Width][Height]int) {
	worldSizeX := Width * tileSize
	worldSizeY := Width * tileSize

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			newQuad := Quad{}
			newQuad.x = worldSizeX / 4 * i
			newQuad.y = worldSizeY / 4 * j
			newQuad.width = worldSizeX / 4
			newQuad.height = worldSizeY / 4
			w.BoundingBoxes = append(w.BoundingBoxes, newQuad)
		}
	}

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if !isWorldBoundary(x, y) && m[x][y] == 1 {
				if m[x+1][y] == 0 || m[x-1][y] == 0 || m[x][y+1] == 0 || m[x][y-1] == 0 {
					newBody := RectAABB{PosX: x * tileSize, PosY: y * tileSize, Width: tileSize, Height: tileSize}
					for i := 0; i < len(w.BoundingBoxes)-1; i++ {
						if w.BoundingBoxes[i].QuadCollide(&newBody) {
							w.BoundingBoxes[i].StaticBodies = append(w.BoundingBoxes[i].StaticBodies, newBody)
						}
					}
				}
			}
		}
	}

}

func (r *Quad) QuadCollide(o *RectAABB) bool {
	if r.x < o.PosX+o.Width && r.x+r.width > o.PosX && r.y+r.height > o.PosY && r.y < o.PosY+o.Height {
		return true
	}

	return false
}

func (q *Quad) DrawCollisionStatic(screen *ebiten.Image, cam [2]int) {
	if len(q.StaticBodies) <= 0 {
		fmt.Println("no rects")
		return
	}

	for _, v := range q.StaticBodies {
		ebitenutil.DrawRect(screen, float64(v.PosX)-float64(cam[0]), float64(v.PosY)-float64(cam[1]), float64(v.Width), float64(v.Height), color.RGBA{100, 10, 10, 255})
	}
}

func (w *World) MovePlayer(r *RectAABB, vel [2]int) {
	for i := 0; i < len(w.BoundingBoxes); i++ {
		if w.BoundingBoxes[i].QuadCollide(r) {
			r.PosX += vel[0]
			r.PosY += vel[1]
			for j := 0; j < len(w.BoundingBoxes[i].StaticBodies); j++ {
				if w.BoundingBoxes[i].StaticBodies[j].CollideWithAABB(r) {
					/*if vel[0] != 0 {
						r.PosX -= vel[0]
					}

					if vel[1] != 0 {
						r.PosY -= vel[1]
					}*/

					if vel[0] == 0 || vel[1] == 0 {
						if vel[0] != 0 {
							r.PosX -= vel[0]
						}

						if vel[1] != 0 {
							r.PosY -= vel[1]
						}
						continue
					}

					r.PosX -= vel[0]
					r.PosY -= vel[1]

					r.PosX += vel[0]
					if w.BoundingBoxes[i].StaticBodies[j].CollideWithAABB(r) {
						r.PosX -= vel[0]
					}
					r.PosY += vel[1]
					if w.BoundingBoxes[i].StaticBodies[j].CollideWithAABB(r) {
						r.PosY -= vel[1]
					}

				}
			}
		}
	}
}

func isWorldBoundary(x, y int) bool {
	return x == 0 || x == Width-1 || y == 0 || y == Height-1
}

type RigidBody struct {
	Body RectAABB
}

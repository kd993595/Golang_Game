package physic

import "github.com/KevinD/LogicAndNightmares/player"

type World struct {
	tilemap       []RectAABB
	BoundingBoxes []Quad
}

type Quad struct {
	x             int
	y             int
	width         int
	height        int
	staticBodies  []RigidBody
	dynamicBodies []RigidBody
}

type RigidBody struct {
	Body RectAABB
}

func (w *World) MovePlayer(p *player.Player, dir Vec2) {
	p.PosX += dir.X
	p.PosY += dir.Y

	for i := 0; i < len(w.tilemap); i++ {

	}
}

package physic

import "github.com/KevinD/LogicAndNightmares/player"

type World struct{
  tilemap []RectAABB
}

type RigidBody struct{
  body RectAABB
}

func (w *World) MovePlayer(p *player.Player,dir Vec2){
  p.PosX += dir.X
  p.PosY += dir.Y

  for i:=0;i<len(w.tilemap);i++{
    
  }
}
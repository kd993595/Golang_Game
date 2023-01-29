package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RectAABB struct {
	PosX   int
	PosY   int
	Width  int
	Height int
}

type Polygon struct {
	Points []Vec2
}

type Vec2 struct {
	X int
	Y int
}

type Vec3 struct {
	X int
	Y int
	Z int
}

/*
Look up vector graphics to render the shapes onto screen must make image to render onto screen where you draw vector graphics
Color variable must have values between 0-255
*/
func (p *Polygon) DebugDrawShape(screen *ebiten.Image, subImg *ebiten.Image, clr Vec3) {
	var path vector.Path

	path.MoveTo(float32(p.Points[0].X), float32(p.Points[0].Y))
	for i := 1; i < len(p.Points); i++ {
		path.LineTo(float32(p.Points[i].X), float32(p.Points[i].Y))
	}

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(clr.X) / float32(0xff)
		vs[i].ColorG = float32(clr.Y) / float32(0xff)
		vs[i].ColorB = float32(clr.Z) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, subImg, op)

}

func (r *RectAABB) DebugDrawRect(scr *ebiten.Image, offset [2]int) {
	ebitenutil.DrawRect(scr, float64(r.PosX-offset[0]), float64(r.PosY-offset[1]), float64(r.Width), float64(r.Height), color.RGBA{100, 255, 100, 255})
}

func (r *RectAABB) MoveRect(dir Vec2) {
	r.PosX += dir.X
	r.PosY += dir.Y
}

func (p *Polygon) MovePolygon(dir Vec2) {
	for i := 0; i < len(p.Points); i++ {
		p.Points[i].X += dir.X
		p.Points[i].Y += dir.Y
	}
}

func (r *RectAABB) CollideWithAABB(o *RectAABB) bool {
	if r.PosX < o.PosX+o.Width && r.PosX+r.Width > o.PosX && r.PosY+r.Height > o.PosY && r.PosY < o.PosY+o.Height {
		return true
	}

	return false
}

func (p *Polygon) CollideWithPolygon(o *Polygon) bool {
	var axis map[Vec2]struct{} = map[Vec2]struct{}{}

	for i := 0; i < len(p.Points)-1; i++ {
		dx := p.Points[i+1].X - p.Points[i].X
		dy := p.Points[i+1].Y - p.Points[i].Y
		normaldx := -dy
		normaldy := dx
		n_axis := Vec2{normaldx, normaldy}
		axis[n_axis] = struct{}{}
	}
	dx := p.Points[len(p.Points)-1].X - p.Points[0].X
	dy := p.Points[len(p.Points)-1].Y - p.Points[0].Y
	normaldx := -dy
	normaldy := dx
	n_axis := Vec2{normaldx, normaldy}
	axis[n_axis] = struct{}{}
	for i := 0; i < len(o.Points)-1; i++ {
		dx := o.Points[i+1].X - o.Points[i].X
		dy := o.Points[i+1].Y - o.Points[i].Y
		normaldx := -dy
		normaldy := dx
		n_axis := Vec2{normaldx, normaldy}
		axis[n_axis] = struct{}{}
	}
	dx = o.Points[len(o.Points)-1].X - o.Points[0].X
	dy = o.Points[len(o.Points)-1].Y - o.Points[0].Y
	normaldx = -dy
	normaldy = dx
	n_axis = Vec2{normaldx, normaldy}
	axis[n_axis] = struct{}{}
	//We have all axis normals now
	for a := range axis {
		pmin := DotProduct(p.Points[0].X, p.Points[0].Y, a.X, a.Y)
		pmax := pmin
		for i := 1; i < len(p.Points); i++ {
			n := DotProduct(p.Points[i].X, p.Points[i].Y, a.X, a.Y)
			if n < pmin {
				pmin = n
			}
			if n > pmax {
				pmax = n
			}
		}
		omin := DotProduct(o.Points[0].X, o.Points[0].Y, a.X, a.Y)
		omax := omin
		for i := 1; i < len(o.Points); i++ {
			n := DotProduct(o.Points[i].X, o.Points[i].Y, a.X, a.Y)
			if n < omin {
				omin = n
			}
			if n > omax {
				omax = n
			}
		}
		if pmin > omax || pmax < omin {

			return false
		}
	}
	return true

}

func DotProduct(x1, y1, x2, y2 int) int {
	return (x1 * x2) + (y1 + y2)
}

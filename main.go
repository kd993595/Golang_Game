package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"

	"math/rand"

	"github.com/KevinD/LogicAndNightmares/components"
	"github.com/KevinD/LogicAndNightmares/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	MapHeight    = 100
	seed         = 47
)

var (
	random       *rand.Rand
	tilemapImg1  *ebiten.Image
	tilemapImg2  *ebiten.Image
	worldImg     *ebiten.Image
	maskImg      *ebiten.Image
	lightTexture *ebiten.Image
	lightPoint   *ebiten.Image
	lightPoint2  *ebiten.Image
	lightPoint3  *ebiten.Image
	playerImg    *ebiten.Image

	//go:embed lightshader.kage
	lighting_go []byte
)

var shaderSrcs = [][]byte{
	lighting_go,
}

func init() {
	random = rand.New(rand.NewSource(seed))
	img, _, err := ebitenutil.NewImageFromFile("./tilemap/tileset.png")
	if err != nil {
		log.Fatal(err)
	}
	tilemapImg1 = ebiten.NewImageFromImage(img)

	img, _, err = ebitenutil.NewImageFromFile("./tilemap/roguelikeSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	tilemapImg2 = ebiten.NewImageFromImage(img)

	img, _, err = ebitenutil.NewImageFromFile("lightPoint.png")
	if err != nil {
		log.Fatal(err)
	}
	lightPoint = ebiten.NewImageFromImage(img)

	img, _, err = ebitenutil.NewImageFromFile("lightPoint2.png")
	if err != nil {
		log.Fatal(err)
	}
	lightPoint2 = ebiten.NewImageFromImage(img)

	img, _, err = ebitenutil.NewImageFromFile("lightPoint3.png")
	if err != nil {
		log.Fatal(err)
	}
	lightPoint3 = ebiten.NewImageFromImage(img)

	worldImg = ebiten.NewImage(components.Width*16, components.Height*16)
	worldImg.Fill(color.RGBA{255, 0, 0, 255})

	maskImg = ebiten.NewImage(screenWidth, screenHeight)
	maskImg.Fill(color.Black)

	lightTexture = ebiten.NewImage(screenWidth, screenHeight)
	lightTexture.Fill(color.RGBA{50, 50, 50, 200})

	img, _, err = ebitenutil.NewImageFromFile("./player/spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}
	playerImg = ebiten.NewImageFromImage(img)

}

type Game struct {
	cam     [2]int
	p       PlayerSystem
	time    int
	shaders map[int]*ebiten.Shader
	space   *components.World
	components.LightSystem
}

type PlayerSystem struct {
	player.Player
	moveBox components.RectAABB
	vel     [2]int
}

func (g *Game) clampCam() {
	if g.cam[0] < 0 {
		g.cam[0] = 0
	}
	if g.cam[1] < 0 {
		g.cam[1] = 0
	}
	if g.cam[0] > components.Width*16-screenWidth {
		g.cam[0] = components.Width*16 - screenWidth
	}
	if g.cam[1] > components.Height*16-screenHeight {
		g.cam[1] = components.Height*16 - screenHeight
	}
}

func (g *Game) Update() error {
	g.time++

	//Use inpututil.IsKeyJustPressed(ebiten.KeyE) to detect one key press not multiple within seconds

	g.p.Idle = true
	g.p.vel = [2]int{0, 0}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.p.vel[0] = 1
		g.p.Idle = false
		g.p.FacingLeft = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.p.vel[1] = 1
		g.p.Idle = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.p.vel[0] = -1
		g.p.Idle = false
		g.p.FacingLeft = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.p.vel[1] = -1
		g.p.Idle = false
	}
	g.cam[0] += int((int(g.p.moveBox.PosX) - g.cam[0] - screenWidth/2) / 20)
	g.cam[1] += int((int(g.p.moveBox.PosY) - g.cam[1] - screenHeight/2) / 20)
	g.clampCam()

	g.p.Update()

	g.space.MovePlayer(&g.p.moveBox, g.p.vel)

	/*if g.shaders == nil {
		g.shaders = map[int]*ebiten.Shader{}
	}
	if _, ok := g.shaders[0]; !ok {
		s, err := ebiten.NewShader([]byte(shaderSrcs[0]))
		if err != nil {
			return err
		}
		g.shaders[0] = s
	}*/

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	lightTexture.Fill(color.RGBA{75, 75, 75, 255})

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(worldImg.SubImage(image.Rect(g.cam[0], g.cam[1], g.cam[0]+screenWidth, g.cam[1]+screenHeight)).(*ebiten.Image), op)

	op.GeoM.Reset()
	if g.p.FacingLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(g.p.moveBox.PosX-g.cam[0]+14), float64(g.p.moveBox.PosY-g.cam[1]-30))
	} else {
		op.GeoM.Translate(float64(g.p.moveBox.PosX-g.cam[0]-6), float64(g.p.moveBox.PosY-g.cam[1]-30))
	}
	screen.DrawImage(playerImg.SubImage(g.p.Frame()).(*ebiten.Image), op)

	//shader part
	/*w, h := screen.Size()
	cx, cy := ebiten.CursorPosition()

	opShader := &ebiten.DrawRectShaderOptions{}
	opShader.Uniforms = map[string]interface{}{
		"Time":       float32(g.time) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"ScreenSize": []float32{float32(w), float32(h)},
	}
	opShader.Images[0] = maskImg
	opShader.Images[1] = lightTexture
	screen.DrawRectShader(w, h, g.shaders[0], opShader)*/

	//lightmap rendering

	cx, cy := ebiten.CursorPosition()
	g.LightSystem.DrawLight([2]int{cx, cy}, lightPoint, color.RGBA{255, 10, 10, 0}, 1)

	//g.LightSystem.DrawLight([2]int{g.p.moveBox.PosX - g.cam[0] - 3, g.p.moveBox.PosY - g.cam[1] - 25}, lightPoint2, nil, playerrad)

	op.GeoM.Reset()
	op.CompositeMode = ebiten.CompositeModeSourceOver
	op.GeoM.Translate(50, 50)
	screen.DrawImage(lightPoint3, op)

	//blending images together using multiply
	op = &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeMultiply
	screen.DrawImage(g.LightSystem.GetLightMap(), op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	var t_p PlayerSystem = PlayerSystem{Player: player.Player{Count: 1, Idle: true, FacingLeft: false}, moveBox: components.RectAABB{PosX: 100, PosY: 100, Width: 10, Height: 5}, vel: [2]int{0, 0}}

	//shader creation
	s, err := ebiten.NewShader([]byte(shaderSrcs[0]))
	if err != nil {
		log.Fatal(err)
	}
	var myShaders map[int]*ebiten.Shader = map[int]*ebiten.Shader{1: s}

	//generate 1 and 0 bitmaps
	newmap := components.WorldGenerator{}
	newmap.Initialize(random)
	newmap.GenerateBitMap()

	//draw all tile image onto image to render
	mymap := components.Tilemap{}
	mymap.Initialize(random)
	mymap.ProcessMap(newmap.GameMap)
	mymap.DrawWorld(worldImg, tilemapImg1, tilemapImg2)

	//create physic world where all objects with hitbox resides
	physicWorld := components.World{}
	physicWorld.Initialize(newmap.GameMap)

	//creating light system
	myLights := components.LightSystem{}
	myLights.Initialize(lightTexture)

	return &Game{p: t_p, cam: [2]int{0, 0}, shaders: myShaders, space: &physicWorld, LightSystem: myLights}
}

func main() {

	/*newmap := worldgen.WorldGenerator{}
	newmap.Initialize(random)
	newmap.GenerateBitMap()

	var levelstring string = ""
	for x := 0; x < worldgen.Width; x++ {
		for y := 0; y < worldgen.Height; y++ {
			levelstring += fmt.Sprintf("%v", newmap.GameMap[x][y])
		}
		levelstring += "\n"
	}

	f, err := os.Create("mapdata.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(levelstring)
	if err2 != nil {
		log.Fatal(err2)
	}*/

	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("Hello, World!")
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

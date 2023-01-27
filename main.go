package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/KevinD/LogicAndNightmares/card"
	"github.com/KevinD/LogicAndNightmares/player"
	"github.com/KevinD/LogicAndNightmares/worldgen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	MapHeight    = 100
	seed         = 47
)

var (
	random      *rand.Rand
	tilemapImg1 *ebiten.Image
	tilemapImg2 *ebiten.Image
	worldImg    *ebiten.Image
	maskImg     *ebiten.Image

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

	worldImg = ebiten.NewImage(worldgen.Width*16, worldgen.Height*16)
	worldImg.Fill(color.RGBA{255, 0, 0, 255})

	maskImg = ebiten.NewImage(screenWidth, screenHeight)
	maskImg.Fill(color.Black)
}

type Game struct {
	cam            vec2
	p              player.Player
	d              card.Deck
	tiles          worldgen.Tilemap
	selectionPhase bool
	selected       int
	time           int
	shaders        map[int]*ebiten.Shader
}

type vec2 struct {
	x int
	y int
}

func (g *Game) clampCam() {
	if g.cam.x < 0 {
		g.cam.x = 0
	}
	if g.cam.y < 0 {
		g.cam.y = 0
	}
	if g.cam.x > worldgen.Width*16-screenWidth {
		g.cam.x = worldgen.Width*16 - screenWidth
	}
	if g.cam.y > worldgen.Height*16-screenHeight {
		g.cam.y = worldgen.Height*16 - screenHeight
	}
}

func (g *Game) Update() error {
	g.time++

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.selectionPhase = !g.selectionPhase
	}
	if g.selectionPhase { //card being displayed
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			g.selected += 1
			if g.selected >= len(g.d.CardDeck)-1 {
				g.selected = len(g.d.CardDeck) - 1
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.selected -= 1
			if g.selected < 0 {
				g.selected = 0
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cam.x += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cam.y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cam.x -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cam.y -= 5
	}
	g.clampCam()

	//g.p.Update()

	if g.shaders == nil {
		g.shaders = map[int]*ebiten.Shader{}
	}
	if _, ok := g.shaders[0]; !ok {
		s, err := ebiten.NewShader([]byte(shaderSrcs[0]))
		if err != nil {
			return err
		}
		g.shaders[0] = s
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.selectionPhase {
		for i := 0; i < len(g.d.CardDeck); i++ {
			if i == g.selected {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(1.2, 1.2)
				op.GeoM.Translate((float64)(90+(58*i)), 200)

				screen.DrawImage(g.d.CardDeck[i].Img, op)
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate((float64)(90+(58*i)), 375)
			screen.DrawImage(g.d.CardDeck[i].Img, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(worldImg.SubImage(image.Rect(g.cam.x, g.cam.y, g.cam.x+screenWidth, g.cam.y+screenHeight)).(*ebiten.Image), op)

	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.p.PosX), float64(g.p.PosY))
	screen.DrawImage(g.p.Frame().(*ebiten.Image), op)*/

	//shader part
	w, h := screen.Size()
	cx, cy := ebiten.CursorPosition()

	opShader := &ebiten.DrawRectShaderOptions{}
	opShader.Uniforms = map[string]interface{}{
		"Time":       float32(g.time) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"ScreenSize": []float32{float32(w), float32(h)},
	}
	opShader.Images[0] = maskImg
	screen.DrawRectShader(w, h, g.shaders[0], opShader)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	var t_p player.Player = player.Player{Count: 0, PosX: 0, PosY: 0, FrameX: 0, FrameY: 0, Width: 28, Height: 34, Img: nil, Speed: 5, Idle: true}

	img, _, err := ebitenutil.NewImageFromFile("./card/card.png")
	if err != nil {
		log.Fatal(err)
	}
	origCardImage := ebiten.NewImageFromImage(img)
	var t_cs [8]card.Card
	for i := 0; i < len(t_cs); i++ {
		t_cs[i] = card.Card{Img: origCardImage}
	}

	t_d := card.Deck{CardDeck: t_cs}

	t_p.Create("./player/player.png")

	//shader creation
	s, err := ebiten.NewShader([]byte(shaderSrcs[0]))
	if err != nil {
		log.Fatal(err)
	}
	var myShaders map[int]*ebiten.Shader = map[int]*ebiten.Shader{1: s}

	newmap := worldgen.WorldGenerator{}
	newmap.Initialize(random)
	newmap.GenerateBitMap()

	mymap := worldgen.Tilemap{}
	mymap.Initialize(random)
	mymap.ProcessMap(newmap.GameMap)
	mymap.DrawWorld(worldImg, tilemapImg1, tilemapImg2)

	return &Game{p: t_p, d: t_d, selectionPhase: false, selected: 0, cam: vec2{0, 0}, tiles: mymap, shaders: myShaders}
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
	//fmt.Print(levelstring)

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

package components

import (
	"math/rand"
	"reflect"
	"sort"
)

const (
	Width             = 200
	Height            = 200
	fillpercent       = 50
	wallThresholdSize = 50
	roomThresholdSize = 30
	passageRadius     = 1
)

type WorldGenerator struct {
	GameMap    [Width][Height]int
	Randgen    *rand.Rand
	worldRooms []*Room
}
type Coord struct {
	tileX int
	tileY int
}
type Room struct {
	tiles                    []Coord
	edgetiles                []Coord
	connectedrooms           []*Room
	roomsize                 int
	isAccessibleFromMainRoom bool
	isMainRoom               bool
}

func (w *WorldGenerator) Initialize(r *rand.Rand) {
	w.GameMap = [Width][Height]int{}
	w.Randgen = r
}

func (thisgame *WorldGenerator) GenerateBitMap() {
	thisgame.populatemap(fillpercent)
	for i := 0; i < 3; i++ {
		thisgame.smoothmap()
	}
	thisgame.createBorders(4)
	thisgame.ProcessMap()

}

func (thisgame *WorldGenerator) populatemap(fillpercent int) {
	//randomly fills map with on and off bits
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if x == 0 || x == Width-1 || y == 0 || y == Height-1 {
				thisgame.GameMap[x][y] = 1
				continue
			}
			filled := thisgame.Randgen.Intn(100) < fillpercent
			if filled {
				thisgame.GameMap[x][y] = 1
			} else {
				thisgame.GameMap[x][y] = 0
			}

		}
	}
}

func (thisgame *WorldGenerator) smoothmap() {
	//applies cellular automata to smooth map
	parallelmap := thisgame.GameMap
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			neighborWallTiles := thisgame.getSurroundingWallCount(x, y)
			if neighborWallTiles > 4 {
				parallelmap[x][y] = 1
			} else if neighborWallTiles < 4 {
				parallelmap[x][y] = 0
			}
		}
	}
	thisgame.GameMap = parallelmap
}

func (thisgame *WorldGenerator) getSurroundingWallCount(gridx int, gridy int) int {
	//returns number of walls surrounding cell
	wallcount := 0
	for neighborX := gridx - 1; neighborX <= gridx+1; neighborX++ {
		for neighborY := gridy - 1; neighborY <= gridy+1; neighborY++ {
			if isInMapRange(neighborX, neighborY) {
				if neighborX != gridx || neighborY != gridy {
					wallcount += thisgame.GameMap[neighborX][neighborY]
				}
			} else {
				wallcount++
			}
		}
	}
	return wallcount
}

func (thisgame *WorldGenerator) getRegionTiles(startX, startY int) []Coord {
	//returns the tile in region bounded
	tiles := make([]Coord, 0, 40)
	mapFlags := [Width][Height]int{}
	tileType := thisgame.GameMap[startX][startY]
	queue := NewQueue()
	queue.Enqueue(Coord{tileX: startX, tileY: startY})
	mapFlags[startX][startY] = 1

	for queue.Count() > 0 {
		tile := queue.Dequeue()
		tiles = append(tiles, tile)

		for x := tile.tileX - 1; x <= tile.tileX+1; x++ {
			for y := tile.tileY - 1; y <= tile.tileY+1; y++ {
				if isInMapRange(x, y) && (y == tile.tileY || x == tile.tileX) {
					if mapFlags[x][y] == 0 && thisgame.GameMap[x][y] == tileType {
						mapFlags[x][y] = 1
						queue.Enqueue(Coord{tileX: x, tileY: y})
					}
				}
			}
		}
	}
	return tiles
}

func (thisgame *WorldGenerator) GetRegions(tiletype int) [][]Coord {
	//returns regions in game depending on type if wall or space
	regions := make([][]Coord, 0, 10)
	mapFlags := [Width][Height]int{}
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if mapFlags[x][y] == 0 && thisgame.GameMap[x][y] == tiletype {
				newRegion := thisgame.getRegionTiles(x, y)
				regions = append(regions, newRegion)

				for _, tile := range newRegion {
					mapFlags[tile.tileX][tile.tileY] = 1
				}
			}
		}
	}
	return regions
}

func (thisgame *WorldGenerator) ProcessMap() {
	//obtains regions bounded by walls and then empty space regions and finally connects empty space regions asigning largest room as the main room
	wallRegions := thisgame.GetRegions(1)
	for _, wallregion := range wallRegions {
		if len(wallregion) < wallThresholdSize {
			for _, tile := range wallregion {
				thisgame.GameMap[tile.tileX][tile.tileY] = 0
			}
		}
	}

	roomRegions := thisgame.GetRegions(0)
	survivingRooms := make([]*Room, 0, 10)

	for _, roomregion := range roomRegions {
		if len(roomregion) < roomThresholdSize {
			for _, tile := range roomregion {
				thisgame.GameMap[tile.tileX][tile.tileY] = 1
			}
		} else {
			survivingRooms = append(survivingRooms, NewRoom(roomregion, thisgame.GameMap))
		}
	}
	sort.Slice(survivingRooms, func(i, j int) bool { return survivingRooms[i].roomsize > survivingRooms[j].roomsize })
	survivingRooms[0].isMainRoom = true
	survivingRooms[0].isAccessibleFromMainRoom = true
	thisgame.worldRooms = survivingRooms
	thisgame.connectClosestRoom(survivingRooms, false)
}

func (thisgame *WorldGenerator) connectClosestRoom(allRooms []*Room, forceAccessibilityFromMainRoom bool) {
	//connects regions together
	roomListA := make([]*Room, 0, 4)
	roomListB := make([]*Room, 0, 4)

	if forceAccessibilityFromMainRoom {
		for _, room := range allRooms {
			if room.isAccessibleFromMainRoom {
				roomListB = append(roomListB, room)
			} else {
				roomListA = append(roomListA, room)
			}
		}
	} else {
		roomListA = allRooms
		roomListB = allRooms
	}

	var bestDistance int = 0
	bestTileA := Coord{}
	bestTileB := Coord{}
	bestRoomA := &Room{}
	bestRoomB := &Room{}
	possibleConnectionFound := false

	for _, roomA := range roomListA { //foreach room
		if !forceAccessibilityFromMainRoom {
			possibleConnectionFound = false
			if len(roomA.connectedrooms) > 0 {
				continue
			}
		}

		for _, roomB := range roomListB { //foreach room
			if reflect.DeepEqual(roomA, roomB) || roomA.isConnected(roomB) {
				continue
			}

			for tileIndexA := 0; tileIndexA < len(roomA.edgetiles); tileIndexA++ {
				for tileIndexB := 0; tileIndexB < len(roomB.edgetiles); tileIndexB++ {
					tileA := roomA.edgetiles[tileIndexA]
					tileB := roomB.edgetiles[tileIndexB]
					distanceBetweenRooms := IntPow(tileA.tileX-tileB.tileX, 2) + IntPow(tileA.tileY-tileB.tileY, 2)

					if distanceBetweenRooms < bestDistance || !possibleConnectionFound {
						bestDistance = distanceBetweenRooms
						possibleConnectionFound = true
						bestTileA = tileA
						bestTileB = tileB
						bestRoomA = roomA
						bestRoomB = roomB
					}
				}
			}
		}
		if possibleConnectionFound && !forceAccessibilityFromMainRoom {
			thisgame.CreatePassage(bestRoomA, bestRoomB, bestTileA, bestTileB)
		}
	}

	if possibleConnectionFound && forceAccessibilityFromMainRoom {
		thisgame.CreatePassage(bestRoomA, bestRoomB, bestTileA, bestTileB)
		thisgame.connectClosestRoom(allRooms, true)
	}
	if !forceAccessibilityFromMainRoom {
		thisgame.connectClosestRoom(allRooms, true)
	}
}

func (thisgame *WorldGenerator) CreatePassage(roomA, roomB *Room, tileA, tileB Coord) {
	//creates the passage carving out the actual space
	connectRooms(roomA, roomB)
	line := getLine(tileA, tileB)
	for _, c := range line {
		thisgame.drawCircle(c, passageRadius)
	}
}

func (thisgame *WorldGenerator) drawCircle(c Coord, r int) {
	//carves out circle around tile with r radius
	for x := -r; x <= r; x++ {
		for y := -r; y <= r; y++ {
			if x*x+y*y <= r*r {
				drawX := c.tileX + x
				drawY := c.tileY + y
				if isInMapRange(drawX, drawY) {
					thisgame.GameMap[drawX][drawY] = 0
				}
			}
		}
	}
}

func getLine(from, to Coord) []Coord {
	//bresenham line algorithm to draw line
	line := make([]Coord, 0, 5)
	x := from.tileX
	y := from.tileY

	dx := to.tileX - from.tileX
	dy := to.tileY - from.tileY

	inverted := false
	step := getSign(dx)
	gradientStep := getSign(dy)

	longest := Abs(dx)
	shortest := Abs(dy)

	if longest < shortest {
		inverted = true
		longest = Abs(dy)
		shortest = Abs(dx)
		step = getSign(dy)
		gradientStep = getSign(dx)
	}

	gradientAccumulation := longest / 2

	for i := 0; i < longest; i++ {
		line = append(line, Coord{tileX: x, tileY: y})

		if inverted {
			y += step
		} else {
			x += step
		}

		gradientAccumulation += shortest
		if gradientAccumulation >= longest {
			if inverted {
				x += gradientStep
			} else {
				y += gradientStep
			}
			gradientAccumulation -= longest
		}
	}
	return line
}

func (thisgame *WorldGenerator) createBorders(bordersize int) {
	//creates borded around map depending on size
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if x < bordersize || x >= Width-bordersize || y < bordersize || y >= Height-bordersize {
				thisgame.GameMap[x][y] = 1
			}
		}
	}
}

func connectRooms(roomA, roomB *Room) {
	//sets room connected to each other
	if roomA.isAccessibleFromMainRoom {
		roomB.setAccessibleFromMainRoom()
	} else if roomB.isAccessibleFromMainRoom {
		roomA.setAccessibleFromMainRoom()
	}
	roomA.connectedrooms = append(roomA.connectedrooms, roomB)
	roomB.connectedrooms = append(roomB.connectedrooms, roomA)
}

func (thisroom *Room) isConnected(otherroom *Room) bool {
	//checks if room is in connected rooms
	for _, ele := range thisroom.connectedrooms {
		if reflect.DeepEqual(otherroom, ele) {
			return true
		}
	}
	return false
}

func (thisroom *Room) setAccessibleFromMainRoom() {
	//sets this room and all connected to be accessible from main
	if !thisroom.isAccessibleFromMainRoom {
		thisroom.isAccessibleFromMainRoom = true
		for _, connectedRoom := range thisroom.connectedrooms {
			connectedRoom.setAccessibleFromMainRoom()
		}
	}
}

func NewRoom(roomtiles []Coord, GameMap [Width][Height]int) *Room {
	//returns new room object
	edgetiles := make([]Coord, 0, 10)
	for _, tile := range roomtiles {
		for x := tile.tileX - 1; x <= tile.tileX+1; x++ {
			for y := tile.tileY; y <= tile.tileY+1; y++ {
				if x == tile.tileX || y == tile.tileY {
					if GameMap[x][y] == 1 {
						edgetiles = append(edgetiles, tile)
					}
				}
			}
		}
	}
	return &Room{tiles: roomtiles, roomsize: len(roomtiles), connectedrooms: make([]*Room, 0, 2), edgetiles: edgetiles}
}

func isInMapRange(x, y int) bool {
	return x >= 0 && x < Width && y >= 0 && y < Height
}

//#region
type Queue struct {
	items []Coord
}

func (qself *Queue) Enqueue(x Coord) {
	qself.items = append(qself.items, x)
}
func (qself *Queue) Dequeue() Coord {
	h := qself.items
	var el Coord
	l := len(h)
	el, qself.items = h[0], h[1:l]
	// Or use this instead for a Stack
	// el, *qself = h[l-1], h[0:l-1]
	return el
}
func (qself *Queue) Count() int {
	h := qself.items
	return len(h)
}
func NewQueue() *Queue {
	return &Queue{items: make([]Coord, 0, 10)}
}
func IntPow(n, m int) int {
	//power function takes int and exponentiates to power m
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func getSign(num int) int {
	switch {
	case num < 0:
		return -1
	case num > 0:
		return 1
	default:
		return 0
	}
}

//#endregion

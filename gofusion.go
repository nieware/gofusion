package main

import (
	"fmt"
	"gopkg.in/v0/qml"
	"gopkg.in/v0/qml/gl"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

/*

TODO:
- swipe gesture support
- high score webservice (hack proof, DOS proof)
- "real" 3D (perspective projection, all tiles in same scene graph -> Qt3D ?)

*/

const tileSize int = 150
const gridSize int = 150
const boardSize int = 4
const maxTileValue int = 11

var board Board
var ctrl Control

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

// ### BOARD ###

type Board struct {
	tiles [16]*Tile

	width  int
	height int

	// has a tile actually moved during the last move?
	moved bool
}

func (b *Board) freeSpaces(t *Tile) int {
	cnt := 0
	for _, t := range b.tiles {
		if t == nil {
			cnt++
		}
	}
	return cnt
}

// insertTile puts a given tile on the board
// returns false if there is no more space
func (b *Board) insertTile(t *Tile) bool {
	for i, ct := range b.tiles {
		if ct == nil {
			//fmt.Println("found space", i)
			b.tiles[i] = t
			//fmt.Println(b)
			return true
		}
	}
	return false
}

// removeTile remove the given tile from the board
// returns false if the tile was not found
func (b *Board) removeTile(t *Tile) bool {
	for i, ct := range b.tiles {
		if ct == t {
			//fmt.Println("removing tile", i)
			b.tiles[i] = nil
			return true
		}
	}
	return false
}

// tileAt returns the tile a position x, y.
// (if there are several tiles, the first one is returned)
func (b *Board) tileAt(x, y int) *Tile {
	for _, t := range b.tiles {
		if t != nil && t.x == x && t.y == y {
			return t
		}
	}
	return nil
}

// addRandomTile generates a random tile and
// puts it on the board
func (b *Board) addRandomTile(maxValue int) {
	v := randGen.Intn(maxValue) + 1
	x, y := 0, 0
	// this is a very simple-minded approach,
	// but it'll have to do for the moment...
	// TODO check for full board! or game over detection
	for {
		x, y = randGen.Intn(b.width), randGen.Intn(b.height)
		if b.tileAt(x, y) == nil {
			break
		}
	}
	b.addTileAt(x, y, v)
}

// addTileAt adds a tile with the specified value at the specified position
func (b *Board) addTileAt(x, y, v int) bool {
	t := ctrl.createTile(v, x, y)
	return b.insertTile(t)
}

func (b *Board) createMergeTest() {
	b.clear()

	b.addTileAt(0, 0, 4)
	b.addTileAt(1, 0, 4)
	b.addTileAt(2, 0, 5)
	b.addTileAt(3, 0, 5)

	b.addTileAt(0, 1, 6)
	b.addTileAt(1, 1, 6)
	b.addTileAt(2, 1, 7)
	b.addTileAt(3, 1, 7)

	b.addTileAt(0, 2, 8)
	b.addTileAt(1, 2, 8)
	b.addTileAt(2, 2, 9)
	b.addTileAt(3, 2, 9)

	b.addTileAt(0, 3, 10)
	b.addTileAt(1, 3, 10)
	b.addTileAt(2, 3, 11)
	b.addTileAt(3, 3, 11)
}

func (b *Board) createGameOverTest() {
	b.clear()

	b.addTileAt(0, 0, 3)
	b.addTileAt(1, 0, 4)
	b.addTileAt(2, 0, 3)
	b.addTileAt(3, 0, 4)

	b.addTileAt(0, 1, 5)
	b.addTileAt(1, 1, 6)
	b.addTileAt(2, 1, 5)
	b.addTileAt(3, 1, 6)

	b.addTileAt(0, 2, 7)
	b.addTileAt(1, 2, 8)
	b.addTileAt(2, 2, 7)
	b.addTileAt(3, 2, 8)

	b.addTileAt(0, 3, 9)
	b.addTileAt(1, 3, 11)
	b.addTileAt(2, 3, 9)
}

// gameOverCheck returns "done" if
// - the board is full and no more moves are possible
// - a 2048 tile is present
// in case of a 2048 tile, "won" is true as well
func (b *Board) gameOverCheck() (done bool, won bool) {
	done = false
	won = false

	// return false if free space
	boardFull := true
	for _, tile := range b.tiles {
		if tile == nil {
			boardFull = false
		} else {
			if tile.Value() == maxTileValue {
				won = true
				return
			}
		}
	}
	if !boardFull {
		return
	}

	// try all possible moves of all possible tiles
	for _, tile := range b.tiles {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx != 0 && dy != 0 || dx == 0 && dy == 0 {
					continue
				}
				newx, newy, _ := b.getMoveTarget(tile, dx, dy)
				if newx != tile.x || newy != tile.y {
					//fmt.Println("tile @", tile.x, tile.y, "can move to", newx, newy)
					return
				}
			}
		}
	}

	// if we get here, no more moves are possible
	done = true
	return
}

func (b *Board) getMoveTarget(tile *Tile, dx, dy int) (x, y int, otherTile *Tile) {
	x, y = tile.x, tile.y
	curx, cury := x, y
	for {
		curx, cury = curx+dx, cury+dy
		if curx < 0 || curx > boardSize-1 || cury < 0 || cury > boardSize-1 {
			break
		}
		candidate := b.tileAt(curx, cury)
		if candidate == nil || candidate.Value() == tile.Value() && tile.Value() < maxTileValue && candidate.NextValue == 0 {
			//fmt.Println("setting x, y", curx, cury)
			x, y = curx, cury
			otherTile = candidate
		}
		if candidate != nil {
			return
		}
	}
	return
}

func (b *Board) clear() {
	for _, t := range b.tiles {
		if t != nil {
			// remove tiles
			b.removeTile(t)
			t.Object.Destroy()
		}
	}
}

func (b *Board) newGame() {
	b.clear()

	board.addRandomTile(2)
	board.addRandomTile(2)
}

func (b *Board) doMove(dx, dy int, next enumStrategy) {
	// get starting point
	x, y, done := next(-1, -1)
	b.moved = false
	for !done {
		t := b.tileAt(x, y)
		if t != nil {
			newx, newy, otherTile := b.getMoveTarget(t, dx, dy)
			if newx != x || newy != y {
				b.moved = true
				//fmt.Println("moved", x, y, newx, newy)
				t.SetPos(newx, newy)
			}
			if otherTile != nil {
				// mark tiles for merging
				t.NextValue = t.Value() + 1
				otherTile.NextValue = -1
				ctrl.enableMerge = true
			}
		}
		x, y, done = next(x, y)
	}
}

func (b *Board) doMerge() {
	for _, t := range b.tiles {
		if t != nil && t.NextValue != 0 {
			if t.NextValue > 0 {
				// marked for promotion
				t.SetValue(t.NextValue)
				t.Call("update")
				ctrl.SetScore(ctrl.score + 1<<uint(t.NextValue))
				t.NextValue = 0
			} else if t.NextValue == -1 {
				// marked for deletion
				// go out in a blaze of glory
				ctrl.Emit(gridSize*t.x+gridSize/2, gridSize*t.y+2*gridSize/2, t.Value())
				b.removeTile(t)
				t.Object.Destroy()
			}
		}
	}
}

func (b *Board) setBounceAnim() {
	for _, t := range board.tiles {
		if t != nil {
			t.SetBounce(true)
		}
	}
}

func (b *Board) setFallAnim() {
	for _, t := range board.tiles {
		if t != nil {
			t.SetFall(true)
		}
	}
}

// ### CONTROL ###

type Control struct {
	Root        qml.Object
	Score       qml.Object
	Message     qml.Object
	SubMessage  qml.Object
	score       int
	hiscore     int
	enableMerge bool

	Running  bool
	settings *GlobalSettings
}

func (ctrl *Control) showScore() {
	ctrl.Score.Set("text", "Score: "+strconv.Itoa(ctrl.score)+" Hi: "+strconv.Itoa(ctrl.hiscore))
}

func (ctrl *Control) SetScore(v int) {
	ctrl.score = v
	ctrl.showScore()
}

func (ctrl *Control) SetRunning(v bool) {
	ctrl.Running = v
	if v {
		ctrl.SetMessage("", "")
	}
}

func (ctrl *Control) SetHiScore(v int) {
	ctrl.hiscore = v
	ctrl.showScore()
	if ctrl.settings != nil {
		ctrl.settings.SetHiScore(uint32(ctrl.hiscore))
	}
}

func (ctrl *Control) SetMessage(m1, m2 string) {
	ctrl.Message.Set("text", m1)
	ctrl.SubMessage.Set("text", m2)
}

func (ctrl *Control) Emit(x, y, level int) {
	component := ctrl.Root.Object("emitterComponent")
	for i := 0; i <= level*2; i++ {
		emitter := component.Create(nil)
		emitter.Set("x", x)
		emitter.Set("y", y)
		emitter.Set("targetX", rand.Intn(240)-120+x)
		emitter.Set("targetY", rand.Intn(240)-120+y)
		emitter.Set("life", rand.Intn(200*level)+400)
		emitter.Set("emitRate", rand.Intn(5*level)+20)
		emitter.ObjectByName("xAnim").Call("start")
		emitter.ObjectByName("yAnim").Call("start")
		emitter.Set("enabled", true)
	}
}

func (ctrl *Control) Done(emitter qml.Object) {
	emitter.Destroy()
}

func (ctrl *Control) HandleKey(key int) {
	if !ctrl.Running {
		ctrl.SetRunning(true)
	}
	switch key {
	case 16777234:
		board.doMove(-1, 0, enumFromLeft)
	case 16777235:
		board.doMove(0, -1, enumFromTop)
	case 16777236:
		board.doMove(1, 0, enumFromRight)
	case 16777237:
		board.doMove(0, 1, enumFromBottom)
		/*default:
		fmt.Println(key)*/
	}
}

func (ctrl *Control) HandleMoveAnimationDone() {
	if ctrl.enableMerge {
		board.doMerge()
		ctrl.enableMerge = false
	}
	if board.moved {
		board.addRandomTile(2)
		done, won := board.gameOverCheck()
		if done {
			if ctrl.score >= ctrl.hiscore {
				ctrl.SetMessage("New High Score!", "click 'Restart'")
				ctrl.SetHiScore(ctrl.score)
				board.setBounceAnim()
			} else {
				ctrl.SetMessage("Game Over!", "click 'Restart'")
				board.setFallAnim()
			}
		}
		if won {
			ctrl.SetMessage("Congratulations, you have done it!", "click 'Restart'")
			if ctrl.score >= ctrl.hiscore {
				ctrl.SetHiScore(ctrl.score)
			}
			board.setBounceAnim()
		}
		board.moved = false
	}
}

func (ctrl *Control) HandleRestartButton() {
	board.newGame()
	ctrl.SetScore(0)
	ctrl.SetMessage("", "")
}

func (ctrl *Control) createTile(value, x, y int) (t *Tile) {
	t = &Tile{}

	component := ctrl.Root.Object("tileComponent")
	tile := component.Create(nil)
	//fmt.Println("t", t, "tile", tile)
	parent := ctrl.Root.ObjectByName("gameCanvas")
	//fmt.Println(parent)
	tile.Set("parent", parent)
	tile.Set("width", tileSize)
	tile.Set("height", tileSize)
	t.Object = tile
	t.SetPos(x, y)
	t.SetValue(value)

	return
}

// ### TILE ###

type Tile struct {
	qml.Object

	models [12]map[string]*Object

	//Value    int
	Rotation  int
	NextValue int

	x int
	y int
}

func (t *Tile) SetPos(x, y int) {
	t.Set("x", gridSize*x)
	t.Set("y", gridSize*y+gridSize/2)
	t.Set("z", y) // for animations, lower lines on screen are in front of higher lines
	t.x = x
	t.y = y
}

func (t *Tile) SetBounce(enabled bool) {
	y0 := t.Int("y")
	y1 := t.Int("y") - (12-t.Value())*5
	//fmt.Println(t.Value(), y0, y1)
	if enabled {
		t.Set("bounceY0", y0)
		t.Set("bounceY1", y1)
		t.Set("bounceDuration", (12-t.Value())*70)
		t.Set("bounceEnable", true)
	} else {
		t.Set("bounceEnable", false)
	}
}

func (t *Tile) SetFall(enabled bool) {
	if enabled {
		t.Set("fallEnable", true)
	} else {
		t.Set("fallEnable", false)
	}
}

func (t *Tile) SetValue(v int) {
	t.Set("nvalue", v)
}

func (t *Tile) Value() int {
	return t.Int("nvalue")
}

func (t *Tile) SetRotation(rotation int) {
	t.Rotation = rotation
	t.Call("update")
}

func (t *Tile) Paint(p *qml.Painter) {
	width := gl.Float(t.Int("width"))

	// TODO: find out how to use perspective projection ?!
	/*gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//gl.Frustum(-100.0, 100.0, -100.0, 100.0, 3.0, 10.0)
	gl.Frustum(-1.0, 1.0, -1.0, 1.0, 3.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)*/

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ShadeModel(gl.SMOOTH)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.Enable(gl.NORMALIZE)

	gl.Clear(gl.DEPTH_BUFFER_BIT)

	gl.Scalef(width/3, width/3, width/3)

	lks := []gl.Float{0.3, 0.3, 0.3, 1.0}
	switch t.Value() {
	/*case 1: // 2
		lks = []gl.Float{0.3, 0.3, 0.1, 1.0}
	case 2: // 4
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 3: // 8
		lks = []gl.Float{0.1, 0.3, 0.3, 1.0}
	case 4: // 16
		lks = []gl.Float{0.1, 0.1, 0.3, 1.0}
	case 5: // 32
		lks = []gl.Float{0.3, 0.1, 0.3, 1.0}
	case 6: // 64
		lks = []gl.Float{0.3, 0.1, 0.1, 1.0}
	case 7: // 128
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 8: // 256
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 9: // 512
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 10: // 1024
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 11: // 2048
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}*/
	/*case 1: // 2
		lks = []gl.Float{0.1, 0.1, 0.6, 1.0}
	case 2: // 4
		lks = []gl.Float{0.2, 0.1, 0.4, 1.0}
	case 3: // 8
		lks = []gl.Float{0.4, 0.1, 0.2, 1.0}
	case 4: // 16
		lks = []gl.Float{0.6, 0.1, 0.1, 1.0}
	case 5: // 32
		lks = []gl.Float{0.4, 0.2, 0.1, 1.0}
	case 6: // 64
		lks = []gl.Float{0.2, 0.4, 0.1, 1.0}
	case 7: // 128
		lks = []gl.Float{0.1, 0.6, 0.1, 1.0}
	case 8: // 256
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 9: // 512
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 10: // 1024
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	case 11: // 2048
		lks = []gl.Float{0.1, 0.3, 0.1, 1.0}
	}*/

	case 1: // 2
		lks = []gl.Float{0.1, 0.1, 0.5, 1.0}
	case 2: // 4
		lks = []gl.Float{0.1, 0.2, 0.3, 1.0}
	case 3: // 8
		lks = []gl.Float{0.1, 0.3, 0.2, 1.0}
	case 4: // 16
		lks = []gl.Float{0.1, 0.5, 0.1, 1.0}
	case 5: // 32
		lks = []gl.Float{0.2, 0.3, 0.1, 1.0}
	case 6: // 64
		lks = []gl.Float{0.3, 0.2, 0.1, 1.0}
	case 7: // 128
		lks = []gl.Float{0.5, 0.1, 0.1, 1.0}
	case 8: // 256
		lks = []gl.Float{0.3, 0.1, 0.2, 1.0}
	case 9: // 512
		lks = []gl.Float{0.2, 0.1, 0.3, 1.0}
	case 10: // 1024
		lks = []gl.Float{0.1, 0.1, 0.7, 1.0}
	case 11: // 2048
		lks = []gl.Float{0.7, 0.3, 0.3, 1.0}
	}

	//lka := []gl.Float{0.3, 0.3, 0.3, 1.0}
	//lkd := []gl.Float{0.3, 0.3, 0.3, 0.0}

	lpos := []gl.Float{-2, 3, 3, 1.0}
	//lpos := []gl.Float{-100, 100, 100, 1.0}

	gl.Enable(gl.LIGHTING)
	//gl.Lightfv(gl.LIGHT0, gl.AMBIENT, lka)
	//gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, lkd)
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, lks)
	gl.Lightfv(gl.LIGHT0, gl.POSITION, lpos)
	gl.Enable(gl.LIGHT0)

	gl.EnableClientState(gl.NORMAL_ARRAY)
	gl.EnableClientState(gl.VERTEX_ARRAY)

	gl.Translatef(1.5, 1.5, 0)
	if t.Value() != 11 {
		gl.Rotatef(-90, 0, 0, 1)
	} else {
		gl.Translatef(0.48, -0.45, 0)
	}
	gl.Rotatef(gl.Float(90+((36000+t.Rotation)%360)), 1, 0, 0)

	gl.Disable(gl.COLOR_MATERIAL)

	model := t.models[t.Value()]
	//fmt.Println("painting", &t, t.Value())
	for _, obj := range model {
		for _, group := range obj.Groups {
			gl.Materialfv(gl.FRONT, gl.AMBIENT, group.Material.Ambient)
			gl.Materialfv(gl.FRONT, gl.DIFFUSE, group.Material.Diffuse)
			gl.Materialfv(gl.FRONT, gl.SPECULAR, group.Material.Specular)
			gl.Materialf(gl.FRONT, gl.SHININESS, group.Material.Shininess)
			gl.VertexPointer(3, gl.FLOAT, 0, group.Vertexes)
			gl.NormalPointer(gl.FLOAT, 0, group.Normals)
			gl.DrawArrays(gl.TRIANGLES, 0, gl.Sizei(len(group.Vertexes)/3))
		}
	}

	gl.Enable(gl.COLOR_MATERIAL)

	//gl.DisableClientState(gl.NORMAL_ARRAY)
	//gl.DisableClientState(gl.VERTEX_ARRAY)
}

// ### INIT / RUN ###
func run(filename string) error {
	qml.Init(nil)
	engine := qml.NewEngine()

	initTiles()

	component, err := engine.LoadFile(filename)
	if err != nil {
		return err
	}

	// init window
	win := component.CreateWindow(nil)
	win.Set("x", 600)
	win.Set("y", 675)

	// init control object (used for communicating with the QML code)
	// and pass it to the QML code.
	ctrl = Control{}
	ctrl.Root = win.Root()
	context := engine.Context()
	context.SetVar("ctrl", &ctrl)

	ctrl.Score = ctrl.Root.ObjectByName("score")
	ctrl.Message = ctrl.Root.ObjectByName("message")
	ctrl.SubMessage = ctrl.Root.ObjectByName("submessage")

	u, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		hiScoreFile := filepath.Join(u.HomeDir, ".gofusion")
		ctrl.settings = NewGlobalSettings(hiScoreFile)
		if ctrl.settings != nil {
			ctrl.hiscore = int(ctrl.settings.GetHiScore())
		}
	}

	board = Board{width: boardSize, height: boardSize}

	board.newGame()
	//board.createGameOverTest()

	win.Show()
	win.Wait()

	return nil
}

func initTiles() error {
	var err error
	var models [12]map[string]*Object

	for i := range models {
		if i == 0 {
			continue
		}

		models[i], err = Read(fmt.Sprintf("model/tile_%04d.obj", 1<<uint(i)))
		if err != nil {
			return err
		}
	}
	//fmt.Println(models)

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{
		{
			Init: func(g *Tile, obj qml.Object) {
				g.Object = obj
				g.models = models
			},
		},
	})

	return nil
}

func main() {
	filename := "gofusion.qml"
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}
	if err := run(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// ### ENUM STRATEGIES ###

func enumFromLeft(cx, cy int) (x, y int, done bool) {
	if cx == -1 {
		return 0, 0, false
	}
	cy++
	if cy >= boardSize {
		cx++
		cy = 0
	}
	return cx, cy, (cx >= boardSize)
}

func enumFromRight(cx, cy int) (x, y int, done bool) {
	if cx == -1 {
		return boardSize - 1, 0, false
	}
	cy++
	if cy >= boardSize {
		cx--
		cy = 0
	}
	return cx, cy, (cx < 0)
}

func enumFromTop(cx, cy int) (x, y int, done bool) {
	if cx == -1 {
		return 0, 0, false
	}
	cx++
	if cx >= boardSize {
		cy++
		cx = 0
	}
	return cx, cy, (cy >= boardSize)
}

func enumFromBottom(cx, cy int) (x, y int, done bool) {
	if cx == -1 {
		return 0, boardSize - 1, false
	}
	cx++
	if cx >= boardSize {
		cy--
		cx = 0
	}
	return cx, cy, (cy < 0)
}

type enumStrategy func(int, int) (int, int, bool)

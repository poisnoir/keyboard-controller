package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/poisnoir/spine-go/client-go"
)

var node *spine.Node
var pub *spine.Publisher[[4][4]float64]

type Game struct{}

// Update runs 60 times per second. Put your input logic here.
func (g *Game) Update() error {
	var goal [4][4]float64

	goal[0][0] = 1
	goal[1][1] = 1
	goal[2][2] = 1
	goal[3][3] = 1

	// 1. Check for single key presses (e.g., pressing Escape to quit)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination // Quits the program safely
	}

	// 2. Check for SIMULTANEOUS key holds
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		goal[1][3] += 0.001
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		goal[1][3] -= 0.001
		fmt.Print("Moving DOWN ↓ ")
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		goal[2][3] += 0.001
		fmt.Print("Moving LEFT ← ")
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		goal[2][3] -= 0.001
		fmt.Print("Moving RIGHT → ")
	}

	// If any keys were pressed, print a new line so it's readable
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
		fmt.Println()
	}

	pub.Publish(goal)

	return nil
}

// Draw is required by Ebiten, but we can leave it empty since we are using the console
func (g *Game) Draw(screen *ebiten.Image) {}

// Layout defines the window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {

	namespace := flag.String("namespace", "rime", "spine namespace to join")
	name := flag.String("name", "r1-change", "publisher name")
	key := flag.String("key", "ppap", "spine namespace key")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var err error

	node, err = spine.CreateNode(*namespace, *key, context.Background(), logger)
	if err != nil {
		panic(err)
	}

	pub, err = spine.NewPublisher[[4][4]float64](node, *name)
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowTitle("Keyboard Tracker")
	ebiten.SetWindowSize(320, 240)

	// Run the game loop
	if err := ebiten.RunGame(&Game{}); err != nil && err != ebiten.Termination {
		panic(err)
	}
}

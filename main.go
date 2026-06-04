package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/poisnoir/spine-go"
)

func main() {
	namespace := flag.String("namespace", "rime", "spine namespace to join")
	name := flag.String("name", "r1-change", "publisher name")
	key := flag.String("key", "ppap", "spine namespace key")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ns, err := spine.JointNamespace(*namespace, *key, logger)
	if err != nil {
		panic(err)
	}

	pub, err := spine.NewPublisher[[4][4]float64](ns, *name)
	if err != nil {
		panic(err)
	}

	var goal [4][4]float64
	for {

		// don't change orientation
		goal[0][0] = 1
		goal[1][1] = 1
		goal[2][2] = 1
		goal[3][3] = 1

		// just z for now
		goal[1][3] = 0.0001

		pub.Publish(goal)

		// make it slow :(
		time.Sleep(50 * time.Microsecond)
	}

}

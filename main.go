package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	pacman "PacmanGo/src"
)

func main() {
	var (
		configFile = flag.String("config-file", "config.json", "path to custom configuration file")
		mazeFile   = flag.String("maze-file", "maze01.txt", "path to a custom maze file")
	)

	var ghostsStatusMx sync.RWMutex
	var pillMx sync.Mutex

	var cfg pacman.Config
	var player pacman.Sprite
	var ghosts []*pacman.Ghost
	var maze []string
	var score int
	var numDots int
	var lives = 3

	var pillTimer *time.Timer

	flag.Parse() // parse the command line arguments

	var ghostNum int
	if len(os.Args) != 2 {
		log.Println("No ghost number provided or too many arguments. Correct usage: go run main.go [ghost number]")
		return
	}

	ghostNum, _ = strconv.Atoi(os.Args[1])

	if ghostNum < 1 || ghostNum > 12 {
		log.Println("Invalid ghost number. It must be between 1 and 12")
		return
	}

	pacman.Initialise()
	defer pacman.Cleanup()

	err := pacman.LoadResources(*mazeFile, *configFile, &maze, &ghosts, &player, &numDots, &cfg, ghostNum)
	if err != nil {
		log.Fatal(err)
		return
	}

	// run the game
	pacman.Run(&player, &maze, &numDots, &score, &lives, &pillMx, &ghostsStatusMx, &ghosts, pillTimer, &cfg)
}

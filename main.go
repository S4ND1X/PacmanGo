package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	pacman "PacmanGo/src"
)

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

func main() {
	flag.Parse() // parse the command line arguments

	// initialize game
	pacman.Initialise()
	defer pacman.Cleanup()

	// load resources
	err := pacman.LoadMaze(*mazeFile, &maze, &ghosts, &player, &numDots)
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}

	err = pacman.LoadConfig(*configFile, &cfg)
	if err != nil {
		log.Println("failed to load configuration:", err)
		return
	}

	// process input (async)
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := pacman.ReadInput()
			if err != nil {
				log.Print("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	// game loop
	for {
		// process movement
		select {
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			pacman.MovePlayer(inp, &player, &maze, &numDots, &score, &pillMx, &ghostsStatusMx, &ghosts, pillTimer, &cfg)
		default:
		}

		pacman.MoveGhosts(&ghosts, &maze)

		// process collisions
		for _, g := range ghosts {
			if player.Row == g.Position.Row && player.Col == g.Position.Col {
				ghostsStatusMx.RLock()
				if g.Status == pacman.GhostStatusNormal {
					lives = lives - 1
					if lives != 0 {
						pacman.MoveCursor(player.Row, player.Col, &cfg)
						fmt.Print(cfg.Death)
						pacman.MoveCursor(len(maze)+2, 0, &cfg)
						ghostsStatusMx.RUnlock()
						go pacman.UpdateGhosts(&ghosts, pacman.GhostStatusNormal, &ghostsStatusMx)
						time.Sleep(1000 * time.Millisecond) //dramatic pause before reseting player position
						player.Row, player.Col = player.StartRow, player.StartCol
					}
				} else if g.Status == pacman.GhostStatusBlue {
					ghostsStatusMx.RUnlock()
					go pacman.UpdateGhosts(&[]*pacman.Ghost{g}, pacman.GhostStatusNormal, &ghostsStatusMx)
					g.Position.Row, g.Position.Col = g.Position.StartRow, g.Position.StartCol
				}
			}
		}

		// update screen
		pacman.PrintScreen(&cfg, &maze, &player, &ghosts, &numDots, &score, &lives, &pillMx, &ghostsStatusMx)

		// check game over
		if numDots == 0 || lives <= 0 {
			if lives == 0 {
				pacman.MoveCursor(player.Row, player.Col, &cfg)
				fmt.Print(cfg.Death)
				pacman.MoveCursor(player.StartRow, player.StartCol-1, &cfg)
				fmt.Print("GAME OVER")
				pacman.MoveCursor(len(maze)+2, 0, &cfg)
			}
			break
		}

		// repeat
		time.Sleep(200 * time.Millisecond)
	}

}

package src

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func processCollisions(player *Sprite, maze *[]string, lives *int, ghostsStatusMx *sync.RWMutex, ghosts *[]*Ghost, cfg *Config) {
	for _, g := range *ghosts {
		if player.Row == g.Position.Row && player.Col == g.Position.Col {
			ghostsStatusMx.RLock()
			if g.Status == GhostStatusNormal {
				*lives = *lives - 1
				if *lives != 0 {
					MoveCursor(player.Row, player.Col, cfg)
					fmt.Print(cfg.Death)
					MoveCursor(len(*maze)+2, 0, cfg)
					ghostsStatusMx.RUnlock()
					go UpdateGhosts(ghosts, GhostStatusNormal, ghostsStatusMx)
					time.Sleep(1000 * time.Millisecond) //dramatic pause before reseting player position
					player.Row, player.Col = player.StartRow, player.StartCol
				}
			} else if g.Status == GhostStatusBlue {
				ghostsStatusMx.RUnlock()
				go UpdateGhosts(&[]*Ghost{g}, GhostStatusNormal, ghostsStatusMx)
				g.Position.Row, g.Position.Col = g.Position.StartRow, g.Position.StartCol
			}
		}
	}
}

func Run(player *Sprite, maze *[]string, numDots, score, lives *int, pillMx *sync.Mutex, ghostsStatusMx *sync.RWMutex, ghosts *[]*Ghost, pillTimer *time.Timer, cfg *Config) {
	// process input with a goroutine
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := ReadInput()
			if err != nil {
				log.Print("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	// while true
	for {
		// process movement
		select {
		case inp := <-input:
			if inp == "ESC" {
				*lives = 0
			}
			MovePlayer(inp, player, maze, numDots, score, pillMx, ghostsStatusMx, ghosts, pillTimer, cfg)
		default:
		}

		MoveGhosts(ghosts, maze)

		// process collisions
		processCollisions(player, maze, lives, ghostsStatusMx, ghosts, cfg)

		// update screen
		PrintScreen(cfg, maze, player, ghosts, numDots, score, lives, pillMx, ghostsStatusMx)

		// check game over
		if *numDots == 0 || *lives <= 0 {
			if *lives == 0 {
				MoveCursor(player.Row, player.Col, cfg)
				fmt.Print(cfg.Death)
				MoveCursor(player.StartRow, player.StartCol-1, cfg)
				fmt.Print("GAME OVER")
				MoveCursor(len(*maze)+2, 0, cfg)
			}
			break
		}

		// repeat
		time.Sleep(100 * time.Millisecond)
	}
}

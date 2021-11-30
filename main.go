package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	pacman "PacmanGo/src"

	"github.com/danicat/simpleansi"
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

func moveCursor(row, col int) {
	if cfg.UseEmoji {
		simpleansi.MoveCursor(row, col*2)
	} else {
		simpleansi.MoveCursor(row, col)
	}
}

func printScreen() {
	simpleansi.ClearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Print(simpleansi.WithBlueBackground(cfg.Wall))
			case '.':
				fmt.Print(cfg.Dot)
			case 'X':
				fmt.Print(cfg.Pill)
			default:
				fmt.Print(cfg.Space)
			}
		}
		fmt.Println()
	}

	moveCursor(player.Row, player.Col)
	fmt.Print(cfg.Player)

	ghostsStatusMx.RLock()
	for _, g := range ghosts {
		moveCursor(g.Position.Row, g.Position.Col)
		if g.Status == pacman.GhostStatusNormal {
			fmt.Printf(cfg.Ghost)
		} else if g.Status == pacman.GhostStatusBlue {
			fmt.Printf(cfg.GhostBlue)
		}
	}
	ghostsStatusMx.RUnlock()

	moveCursor(len(maze)+1, 0)

	livesRemaining := strconv.Itoa(lives) //converts lives int to a string
	if cfg.UseEmoji {
		livesRemaining = getLivesAsEmoji()
	}

	fmt.Println("Score:", score, "\tLives:", livesRemaining)
}

//concatenate the correct number of player emojis based on lives
func getLivesAsEmoji() string {
	buf := bytes.Buffer{}
	for i := lives; i > 0; i-- {
		buf.WriteString(cfg.Player + " ") //concatenate player emoji
	}
	return buf.String()
}

var pillTimer *time.Timer

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func moveGhosts() {
	for _, g := range ghosts {
		dir := drawDirection()
		g.Position.Row, g.Position.Col = pacman.MakeMove(g.Position.Row, g.Position.Col, dir, maze)
	}
}

func main() {
	flag.Parse()

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
			pacman.MovePlayer(inp, &player, maze, &numDots, &score, &pillMx, &ghostsStatusMx, ghosts, pillTimer, &cfg)
		default:
		}

		moveGhosts()

		// process collisions
		for _, g := range ghosts {
			if player.Row == g.Position.Row && player.Col == g.Position.Col {
				ghostsStatusMx.RLock()
				if g.Status == pacman.GhostStatusNormal {
					lives = lives - 1
					if lives != 0 {
						moveCursor(player.Row, player.Col)
						fmt.Print(cfg.Death)
						moveCursor(len(maze)+2, 0)
						ghostsStatusMx.RUnlock()
						go pacman.UpdateGhosts(ghosts, pacman.GhostStatusNormal, &ghostsStatusMx)
						time.Sleep(1000 * time.Millisecond) //dramatic pause before reseting player position
						player.Row, player.Col = player.StartRow, player.StartCol
					}
				} else if g.Status == pacman.GhostStatusBlue {
					ghostsStatusMx.RUnlock()
					go pacman.UpdateGhosts([]*pacman.Ghost{g}, pacman.GhostStatusNormal, &ghostsStatusMx)
					g.Position.Row, g.Position.Col = g.Position.StartRow, g.Position.StartCol
				}
			}
		}

		// update screen
		printScreen()

		// check game over
		if numDots == 0 || lives <= 0 {
			if lives == 0 {
				moveCursor(player.Row, player.Col)
				fmt.Print(cfg.Death)
				moveCursor(player.StartRow, player.StartCol-1)
				fmt.Print("GAME OVER")
				moveCursor(len(maze)+2, 0)
			}
			break
		}

		// repeat
		time.Sleep(200 * time.Millisecond)
	}

}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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







const (
	GhostStatusNormal pacman.GhostStatus = "Normal"
	GhostStatusBlue   pacman.GhostStatus = "Blue"
)

var ghostsStatusMx sync.RWMutex
var pillMx sync.Mutex

type config struct {
	Player           string        `json:"player"`
	Ghost            string        `json:"ghost"`
	GhostBlue        string        `json:"ghost_blue"`
	Wall             string        `json:"wall"`
	Dot              string        `json:"dot"`
	Pill             string        `json:"pill"`
	Death            string        `json:"death"`
	Space            string        `json:"space"`
	UseEmoji         bool          `json:"use_emoji"`
	PillDurationSecs time.Duration `json:"pill_duration_secs"`
}

var cfg config
var player pacman.Sprite
var ghosts []*pacman.Ghost
var maze []string
var score int
var numDots int
var lives = 3

func loadConfig(file string) error {
	f, err := os.Open(file) // open file
	if err != nil {
		return err 
	}

	defer f.Close() // close file on return

	decoder := json.NewDecoder(f) // create json decoder
	err = decoder.Decode(&cfg) // decode file into config struct
	if err != nil {
		return err
	}

	return nil
}

func loadMaze(file string) error {
	f, err := os.Open(file) // open file
	if err != nil { 
		return err
	}
	defer f.Close() // close file on return

	scanner := bufio.NewScanner(f) // create scanner
	for scanner.Scan() {  
		line := scanner.Text() // get next line
		maze = append(maze, line) // append line to maze
	}

	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P':
				player = pacman.Sprite{row, col, row, col}
			case 'G':
				ghosts = append(ghosts, &pacman.Ghost{pacman.Sprite{row, col, row, col}, GhostStatusNormal})
			case '.':
				numDots++
			}
		}
	}

	return nil
}

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
		if g.Status == GhostStatusNormal {
			fmt.Printf(cfg.Ghost)
		} else if g.Status == GhostStatusBlue {
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
		buf.WriteString(cfg.Player)
	}
	return buf.String()
}



func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze)-1 {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}

	if maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func movePlayer(dir string) {
	player.Row, player.Col = makeMove(player.Row, player.Col, dir)

	removeDot := func(row, col int) {
		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
	}

	switch maze[player.Row][player.Col] {
	case '.':
		numDots--
		score++
		removeDot(player.Row, player.Col)
	case 'X':
		score += 10
		removeDot(player.Row, player.Col)
		go processPill()
	}
}

func updateGhosts(ghosts []*pacman.Ghost, ghostStatus pacman.GhostStatus) {
	ghostsStatusMx.Lock()
	defer ghostsStatusMx.Unlock()
	for _, g := range ghosts {
		g.Status = ghostStatus
	}
}

var pillTimer *time.Timer

func processPill() {
	pillMx.Lock()
	go updateGhosts(ghosts, GhostStatusBlue)
	if pillTimer != nil {
		pillTimer.Stop()
	}
	pillTimer = time.NewTimer(time.Second * cfg.PillDurationSecs)
	pillMx.Unlock()
	<-pillTimer.C
	pillMx.Lock()
	pillTimer.Stop()
	go updateGhosts(ghosts, GhostStatusNormal)
	pillMx.Unlock()
}

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
		g.Position.Row, g.Position.Col = makeMove(g.Position.Row, g.Position.Col, dir)
	}
}



func main() {
	flag.Parse()

	// initialize game
	pacman.Initialise()
	defer pacman.Cleanup()

	// load resources
	err := loadMaze(*mazeFile)
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}

	err = loadConfig(*configFile)
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
			movePlayer(inp)
		default:
		}

		moveGhosts()

		// process collisions
		for _, g := range ghosts {
			if player.Row == g.Position.Row && player.Col == g.Position.Col {
				ghostsStatusMx.RLock()
				if g.Status == GhostStatusNormal {
					lives = lives - 1
					if lives != 0 {
						moveCursor(player.Row, player.Col)
						fmt.Print(cfg.Death)
						moveCursor(len(maze)+2, 0)
						ghostsStatusMx.RUnlock()
						go updateGhosts(ghosts, GhostStatusNormal)
						time.Sleep(1000 * time.Millisecond) //dramatic pause before reseting player position
						player.Row, player.Col = player.StartRow, player.StartCol
					}
				} else if g.Status == GhostStatusBlue {
					ghostsStatusMx.RUnlock()
					go updateGhosts([]*pacman.Ghost{g}, GhostStatusNormal)
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

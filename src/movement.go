package src

import (
	"sync"
	"time"
)

func MakeMove(oldRow, oldCol int, dir string, maze []string) (newRow, newCol int) {
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

func MovePlayer(dir string, player *Sprite, maze []string, numDots, score *int, pillMx *sync.Mutex, ghostsStatusMx *sync.RWMutex, ghosts []*Ghost, pillTimer *time.Timer, cfg *Config) {
	player.Row, player.Col = MakeMove(player.Row, player.Col, dir, maze)

	removeDot := func(row, col int) {
		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
	}

	switch maze[player.Row][player.Col] {
	case '.':
		*numDots--
		*score++
		removeDot(player.Row, player.Col)
	case 'X':
		*score += 10
		removeDot(player.Row, player.Col)
		go ProcessPill(pillMx, ghostsStatusMx, ghosts, pillTimer, cfg)
	}
}

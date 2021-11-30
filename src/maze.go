package src

import (
	"bufio"
	"os"
)

func LoadMaze(file string, maze *[]string, ghosts *[]*Ghost, player *Sprite, numDots *int, ghostsNum int) error {
	f, err := os.Open(file) // open file
	if err != nil {
		return err
	}
	defer f.Close() // close file on return

	scanner := bufio.NewScanner(f) // create scanner
	for scanner.Scan() {
		line := scanner.Text()      // get next line
		*maze = append(*maze, line) // append line to maze
	}

	for row, line := range *maze {
		for col, char := range line {
			switch char {
			case 'P':
				*player = Sprite{Row: row, Col: col, StartRow: row, StartCol: col}
			case 'G':
				if len(*ghosts) < ghostsNum {
					*ghosts = append(*ghosts, &Ghost{Position: Sprite{Row: row, Col: col, StartRow: row, StartCol: col}, Status: GhostStatusNormal})
				}
			case '.':
				*numDots++
			}
		}
	}

	return nil
}

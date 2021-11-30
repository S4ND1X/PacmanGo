package src

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"github.com/danicat/simpleansi"
)

func PrintScreen(cfg *Config, maze *[]string, player *Sprite, ghosts *[]*Ghost, numDots, score, lives *int, pillMx *sync.Mutex, ghostsStatusMx *sync.RWMutex) {
	simpleansi.ClearScreen()
	for _, line := range *maze {
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

	MoveCursor(player.Row, player.Col, cfg)
	fmt.Print(cfg.Player)

	ghostsStatusMx.RLock()
	for _, g := range *ghosts {
		MoveCursor(g.Position.Row, g.Position.Col, cfg)
		if g.Status == GhostStatusNormal {
			fmt.Printf(cfg.Ghost)
		} else if g.Status == GhostStatusBlue {
			fmt.Printf(cfg.GhostBlue)
		}
	}
	ghostsStatusMx.RUnlock()

	MoveCursor(len(*maze)+1, 0, cfg)

	livesRemaining := strconv.Itoa(*lives) //converts lives int to a string
	if cfg.UseEmoji {
		livesRemaining = getLivesAsEmoji(cfg, lives)
	}

	fmt.Println("Score:", *score, "\tLives:", livesRemaining)
}

func MoveCursor(row, col int, cfg *Config) {
	if cfg.UseEmoji {
		simpleansi.MoveCursor(row, col*2)
	} else {
		simpleansi.MoveCursor(row, col)
	}
}

//concatenate the correct number of player emojis based on lives
func getLivesAsEmoji(cfg *Config, lives *int) string {
	buf := bytes.Buffer{}
	for i := *lives; i > 0; i-- {
		buf.WriteString(cfg.Player + " ") //concatenate player emoji
	}
	return buf.String()
}

package src

import (
	"log"
)

func LoadResources(mazefile, configFile string, maze *[]string, ghosts *[]*Ghost, player *Sprite, numDots *int, cfg *Config) error {
	err := LoadMaze(mazefile, maze, ghosts, player, numDots)
	if err != nil {
		log.Println("failed to load maze:", err)
		return err
	}

	err = LoadConfig(configFile, cfg)
	if err != nil {
		log.Println("failed to load configuration:", err)
		return err
	}

	return nil
}

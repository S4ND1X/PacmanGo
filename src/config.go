package src

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
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

func LoadConfig(file string, cfg *Config) error {
	f, err := os.Open(file) // open file
	if err != nil {
		return err
	}

	defer f.Close() // close file on return

	decoder := json.NewDecoder(f) // create json decoder
	err = decoder.Decode(&cfg)    // decode file into config struct
	if err != nil {
		return err
	}

	return nil
}

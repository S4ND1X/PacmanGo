package src

import (
)

type Ghost struct {
	Position Sprite
	Status   GhostStatus
}

type GhostStatus string
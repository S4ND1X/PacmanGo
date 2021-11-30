package src

import (
	"sync"
)

type Ghost struct {
	Position Sprite
	Status   GhostStatus
}

type GhostStatus string

const (
	GhostStatusNormal GhostStatus = "Normal"
	GhostStatusBlue   GhostStatus = "Blue"
)

func UpdateGhosts(ghosts []*Ghost, ghostStatus GhostStatus, ghostsStatusMx *sync.RWMutex) {
	ghostsStatusMx.Lock()
	defer ghostsStatusMx.Unlock()
	for _, g := range ghosts {
		g.Status = ghostStatus
	}
}

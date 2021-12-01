package src

import (
	"sync"
	"time"
)

func ProcessPill(pillMx *sync.Mutex, ghostsStatusMx *sync.RWMutex, ghosts *[]*Ghost, pillTimer *time.Timer, cfg *Config) {
	pillMx.Lock()
	go UpdateGhosts(ghosts, GhostStatusBlue, ghostsStatusMx)
	if pillTimer != nil {
		pillTimer.Stop()
	}
	pillTimer = time.NewTimer(time.Second * cfg.PillDurationSecs)
	pillMx.Unlock()
	<-pillTimer.C
	pillMx.Lock()
	pillTimer.Stop()
	go UpdateGhosts(ghosts, GhostStatusNormal, ghostsStatusMx)
	pillMx.Unlock()

}

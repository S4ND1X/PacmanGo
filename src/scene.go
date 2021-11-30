package pacman

import (
	"github.com/hajimehoshi/ebiten"
	pacimages "github.com/kgosse/pacmanresources/images"
)

type scene struct {
	matrix           [][]elem
	wallSurface      *ebiten.Image
	images           map[elem]*ebiten.Image
	stage            *stage
	dotManager       *dotManager
	bigDotManager    *bigDotManager
	player           *player
	ghostManager     *ghostManager
	textManager      *textManager
	fruitManager     *fruitManager
	lives            int
	pointManager     *pointManager
	explosionManager *explosionManager
	sounds           *sounds
	over             bool
}

func newScene(st *stage) *scene {
	s := &scene{}
	s.stage = st
	if s.stage == nil {
		s.stage = defaultStage
	}
	s.lives = s.stage.lives
	s.images = make(map[elem]*ebiten.Image)
	s.dotManager = newDotManager()
	s.explosionManager = newExplosionManager()
	s.bigDotManager = newBigDotManager()
	s.ghostManager = newGhostManager()
	s.pointManager = newPointManager()
	s.sounds = newSounds()
	h := len(s.stage.matrix)
	w := len(s.stage.matrix[0])
	s.textManager = newTextManager(w*stageBlocSize, h*stageBlocSize)
	s.loadImages()
	s.createStage()
	s.buildWallSurface()
	s.textManager.entranceAnim(true)
	s.sounds.playEntrance()
	return s
}

func (s *scene) createStage() {
	h := len(s.stage.matrix)
	w := len(s.stage.matrix[0])
	s.matrix = make([][]elem, h)
	for i := 0; i < h; i++ {
		s.matrix[i] = make([]elem, w)
		for j := 0; j < w; j++ {
			c := s.stage.matrix[i][j] - '0'
			if c <= 9 {
				s.matrix[i][j] = elem(c)
			} else {
				s.matrix[i][j] = elem(s.stage.matrix[i][j] - 'a' + 10)
			}

			switch s.matrix[i][j] {
			case dotElem:
				s.dotManager.add(i, j)
			case bigDotElem:
				s.bigDotManager.add(i, j)
			case playerElem:
				s.player = newPlayer(i, j)
			case fruitElem:
				s.fruitManager = newFruitManager(float64(j*stageBlocSize), float64(i*stageBlocSize))
			case blinkyElem:
				s.ghostManager.addGhost(i, j, blinkyElem)
			case inkyElem:
				s.ghostManager.addGhost(i, j, inkyElem)
			case pinkyElem:
				s.ghostManager.addGhost(i, j, pinkyElem)
			case clydeElem:
				s.ghostManager.addGhost(i, j, clydeElem)
			}
		}
	}
}

func (s *scene) screenWidth() int {
	w := len(s.stage.matrix[0])
	return w * stageBlocSize
}

func (s *scene) screenHeight() int {
	h := len(s.stage.matrix)
	sizeH := ((h*stageBlocSize)/backgroundImageSize + 2) * backgroundImageSize
	return sizeH
}

func (s *scene) buildWallSurface() {
	h := len(s.stage.matrix)
	w := len(s.stage.matrix[0])

	sizeW := ((w*stageBlocSize)/backgroundImageSize + 1) * backgroundImageSize
	sizeH := ((h*stageBlocSize)/backgroundImageSize + 2) * backgroundImageSize
	s.wallSurface, _ = ebiten.NewImage(sizeW, sizeH, ebiten.FilterDefault)

	for i := 0; i < sizeH/backgroundImageSize; i++ {
		y := float64(i * backgroundImageSize)
		for j := 0; j < sizeW/backgroundImageSize; j++ {
			op := &ebiten.DrawImageOptions{}
			x := float64(j * backgroundImageSize)
			op.GeoM.Translate(x, y)
			s.wallSurface.DrawImage(s.images[backgroundElem], op)
		}
	}

	for i := 0; i < h; i++ {
		y := float64(i * stageBlocSize)
		for j := 0; j < w; j++ {
			if !isWall(s.matrix[i][j]) {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			x := float64(j * stageBlocSize)
			op.GeoM.Translate(x, y)
			s.wallSurface.DrawImage(s.images[s.matrix[i][j]], op)
		}
	}
}

func (s *scene) loadImages() {
	for i := w0; i <= w24; i++ {
		s.images[i] = loadImage(pacimages.WallImages[i])
	}
	s.images[backgroundElem] = loadImage(pacimages.Background_png)
}

func (s *scene) move(in input) {
	s.explosionManager.move()
	if !s.over {
		s.ghostManager.move(s.matrix, s.player.curPos)
	}
	if s.lives > 0 {
		s.player.move(s.matrix, in, s.afterPacmanExplosion)
	}
}

func (s *scene) detectCollision() {
	y, x := s.player.screenPos()

	s.dotManager.detectCollision(s.matrix, s.player.curPos, s.afterPacmanDotCollision)
	s.fruitManager.detectCollision(y, x, s.afterPacmanFruitCollision)
	s.bigDotManager.detectCollision(s.matrix, s.player.curPos, s.afterPacmanBigDotCollision)
	s.ghostManager.detectCollision(y, x, s.afterPacmanGhostCollision)
}

func (s *scene) afterPacmanDotCollision() {
	s.player.score += 10
	s.dotManager.delete(s.player.curPos)
	s.matrix[s.player.curPos.y][s.player.curPos.x] = empty
	if !s.over && s.won() {
		s.victory()
	}
}

func (s *scene) afterPacmanFruitCollision() {
	y, x := s.player.screenPos()
	s.player.score += 100
	s.pointManager.show(0, x, y)
	s.sounds.playEeatFruit()
	s.lives++
	if s.lives > s.stage.maxLives {
		s.lives = s.stage.maxLives
	}
}

func (s *scene) afterPacmanBigDotCollision() {
	s.player.score += 50
	s.bigDotManager.delete(s.player.curPos)
	s.matrix[s.player.curPos.y][s.player.curPos.x] = empty
	if !s.over && s.won() {
		s.victory()
		return
	}
	s.ghostManager.makeVulnerable()
	s.sounds.playWail()

}

func (s *scene) afterPacmanGhostCollision(vulnerable bool, y, x float64) {
	if vulnerable {
		s.explosionManager.addExplosion(pacimages.GhostParticle_png, x, y)
		s.sounds.playEeatGhost()
		eaten := s.ghostManager.eaten
		if eaten == 1 {
			s.player.score += 200
		} else if eaten == 2 {
			s.player.score += 400
		} else if eaten == 3 {
			s.player.score += 800
		} else {
			s.player.score += 1600
		}
		s.pointManager.show(eaten, x, y)
	} else {
		s.sounds.playDeath()
		s.player.explode()
	}
}

func (s *scene) afterPacmanExplosion() {
	s.ghostManager.reset(s.explosionManager, false)
	x, y := s.textManager.livesPos(s.lives)
	s.explosionManager.addExplosion(pacimages.PacParticle_png, x-16, y+16)
	s.lives--
	if s.lives == 0 {
		s.player.gameover()
	}
}

func (s *scene) reinit() {
	s.dotManager.reinit(s.matrix)
	s.bigDotManager.reinit(s.matrix)
	s.player.reinit()
	s.textManager.reinit()
	s.lives = s.stage.maxLives - 1
	s.over = false
	s.ghostManager.reinit()
	s.explosionManager.reinit()
	s.textManager.entranceAnim(true)
	s.sounds.pause()
	s.sounds.playEntrance()
}

func (s *scene) won() bool {
	if s.dotManager.empty() && s.bigDotManager.empty() {
		return true
	}
	return false
}

func (s *scene) victory() {
	s.over = true
	s.ghostManager.reset(s.explosionManager, true)
	s.sounds.pause()
	s.sounds.playApplause()
}

func (s *scene) update(screen *ebiten.Image, in input) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if in == sKey {
		s.sounds.toggleSound()
	}

	if in == rKey {
		s.reinit()
	} else if !s.textManager.entrance {
		s.move(in)
		s.detectCollision()
	}

	screen.Clear()
	screen.DrawImage(s.wallSurface, nil)
	s.dotManager.draw(screen)
	s.bigDotManager.draw(screen)
	s.ghostManager.draw(screen)
	s.player.draw(screen)
	s.fruitManager.draw(screen)
	s.pointManager.draw(screen)
	s.explosionManager.draw(screen)
	s.textManager.draw(screen, s.player.score, s.lives, s.player.images[1], s.won(), s.sounds.status())
	s.sounds.playSiren(s.won())
	return nil
}

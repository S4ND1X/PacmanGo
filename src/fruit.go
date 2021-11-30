package pacman

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	pacimages "github.com/kgosse/pacmanresources/images"
)

type fruitManager struct {
	fruits  [3]*ebiten.Image
	x, y    float64
	count   int
	curr    int
	alpha   float64
	show    bool
	visible bool
}

func newFruitManager(x, y float64) *fruitManager {
	fm := &fruitManager{
		x:     x,
		y:     y,
		alpha: 0,
		count: -400,
		show:  true,
	}
	fm.loadImages()
	return fm
}

func (fm *fruitManager) loadImages() {
	copy(fm.fruits[:], loadImages(pacimages.FruitImages[:]))
}

func (fm *fruitManager) draw(screen *ebiten.Image) {
	fm.update()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(fm.x, fm.y)
	op.ColorM.Scale(1, 1, 1, fm.alpha)
	screen.DrawImage(fm.fruits[fm.curr], op)
}

func (fm *fruitManager) update() {
	if fm.show {
		fm.count++
		if fm.count >= 70 {
			fm.alpha += 0.01
			if fm.alpha > 1 {
				fm.alpha = 1
			}
		}
	} else {
		fm.count--
		fm.alpha -= 0.01
		if fm.alpha < 0 {
			fm.alpha = 0
		}
	}

	if fm.alpha >= 0.1 {
		fm.visible = true
	} else {
		fm.visible = false
	}

	if fm.count == 400 {
		fm.show = false
	} else if fm.count <= -500 && !fm.show {
		fm.show = true
		fm.curr = (fm.curr + 1) % len(fm.fruits)
	}
}

func (fm *fruitManager) detectCollision(pY, pX float64, cb func()) {
	if !fm.visible {
		return
	}
	if math.Abs(pY-fm.y) < 32 && math.Abs(pX-fm.x) < 32 {
		fm.show = true
		fm.count = -900
		fm.visible = false
		fm.alpha = 0
		fm.curr = (fm.curr + 1) % len(fm.fruits)
		cb()
	}
}

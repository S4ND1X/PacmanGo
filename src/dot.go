package pacman

import (
	"container/list"

	"github.com/hajimehoshi/ebiten"
	pacimages "github.com/kgosse/pacmanresources/images"
)

type dotManager struct {
	dots  *list.List
	gc    *list.List
	image *ebiten.Image
}

func newDotManager() *dotManager {
	d := &dotManager{}
	d.dots = list.New()
	d.gc = list.New()
	d.loadImage()
	return d
}

func (d *dotManager) loadImage() {
	d.image = loadImage(pacimages.Dot_png)
}

func (d *dotManager) add(y, x int) {
	d.dots.PushBack(pos{y, x})
}

func (d *dotManager) draw(sc *ebiten.Image) {
	for e := d.dots.Front(); e != nil; e = e.Next() {
		v := e.Value.(pos)
		x := float64(v.x * stageBlocSize)
		y := float64(v.y * stageBlocSize)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		sc.DrawImage(d.image, op)
	}
}

func (d *dotManager) delete(p pos) {
	for e := d.dots.Front(); e != nil; e = e.Next() {
		v := e.Value.(pos)
		if v.x == p.x && v.y == p.y {
			d.gc.PushBack(d.dots.Remove(e).(pos))
			return
		}
	}
}

func (d *dotManager) detectCollision(m [][]elem, p pos, cb func()) {
	if m[p.y][p.x] == dotElem {
		cb()
	}
}

func (d *dotManager) reinit(m [][]elem) {
	e := d.gc.Front()
	for {
		if e == nil {
			break
		}
		v := e.Value.(pos)
		cur := e
		e = e.Next()
		d.dots.PushBack(d.gc.Remove(cur))
		m[v.y][v.x] = dotElem
	}
}

func (d *dotManager) empty() bool {
	if d.dots.Len() == 0 {
		return true
	}
	return false
}

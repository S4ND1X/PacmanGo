package pacman

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
)

func canMove(m [][]elem, p pos) bool {
	return !isWall(m[p.y][p.x])
}

func isWall(e elem) bool {
	if w0 <= e && e <= w24 {
		return true
	}
	return false
}

func addPosDir(d input, p pos) pos {
	r := pos{p.y, p.x}

	switch d {
	case up:
		r.y--
	case right:
		r.x++
	case down:
		r.y++
	case left:
		r.x--
	}

	if r.x < 0 {
		r.x = 0
	}
	if r.y < 0 {
		r.y = 0
	}

	return r
}

func oppDir(d input) input {
	switch d {
	case up:
		return down
	case right:
		return left
	case down:
		return up
	case left:
		return right
	default:
		return 0
	}
}

func loadImage(b []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(b))
	handleError(err)
	ebImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	handleError(err)
	return ebImg
}

func loadImages(images [][]byte) []*ebiten.Image {
	var res []*ebiten.Image
	size := len(images)
	for i := 0; i < size; i++ {
		res = append(res, loadImage(images[i]))
	}
	return res
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

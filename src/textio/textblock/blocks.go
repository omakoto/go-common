package textblock

import (
	"bytes"
	"github.com/omakoto/go-common/src/common"
	"regexp"
)

type TextBlock struct {
	lines [][]byte
}

var lf = []byte("\n")

func copyLines(lines ...[]byte) [][]byte {
	ret := make([][]byte, len(lines), len(lines))

	for i, v := range lines {
		ret[i] = make([]byte, len(v), len(v))
		copy(ret[i], v)
	}

	return ret
}

func NewFromBuffer(buf []byte) *TextBlock {
	lines := bytes.Split(buf, lf)
	return &TextBlock{lines: lines}
}

func NewFromLines(lines [][]byte) *TextBlock {
	return &TextBlock{lines: copyLines(lines...)}
}

func (b *TextBlock) Lines() [][]byte {
	return copyLines(b.lines...)
}

func (b *TextBlock) Bytes() []byte {
	return bytes.Join(b.lines, lf)
}

func (b *TextBlock) Append(target ...*TextBlock) *TextBlock {
	for _, t := range target {
		b.lines = append(b.lines, copyLines(t.lines...)...)
	}
	return b
}

func (b *TextBlock) Slice(start, end int) *TextBlock {
	return NewFromLines(b.lines[start:end])
}

func (b *TextBlock) copyRegion(startX, startY, endX, endY int, doCut bool) *TextBlock {
	if startX < 0 {
		common.Panicf("Invalid startX: %i", startX)
	}
	if startY < 0 {
		common.Panicf("Invalid startY: %i", startY)
	}
	if endX < startX {
		common.Panicf("endX %i must not be smaller than startX %i", endX, startX)
	}
	if endY < startY {
		common.Panicf("endY %i must not be smaller than startY %i", endY, startY)
	}
	cut := TextBlock{lines: make([][]byte, endY-startY+1)}

	toY := 0
	for y := startY; y < endY; y++ {
		if startY >= len(b.lines) {
			continue
		}
		lineLen := len(b.lines[y])
		if startX >= lineLen {
			continue
		}
		realEndX := endX
		if realEndX > lineLen  {
			realEndX = lineLen
		}

		cut.lines[toY] = make([]byte, realEndX-startX+1, realEndX-startX+1)
		copy(cut.lines[toY], b.lines[y][startX:realEndX])
		toY++

		if !doCut {
			continue
		}

		panic("not tested")

		fromX := endX
		for x := startX; x < endX; x++ {
			if fromX >= len(b.lines[y]) {
				break
			}
			b.lines[y][x] = b.lines[y][fromX]
		}
	}
	return &cut
}

func (b *TextBlock) Copy(startX, startY, endX, endY int) *TextBlock {
	return b.copyRegion(startX, startY, endX, endY, false)
}

func (b *TextBlock) Cut(startX, startY, endX, endY int) *TextBlock {
	return b.copyRegion(startX, startY, endX, endY, true)
}

func (b *TextBlock) FindFirst(regexp *regexp.Regexp) (int, int) {
	for y, line := range b.lines {
		loc := regexp.FindIndex(line)
		if loc != nil {
			return loc[0], y
		}
	}
	return -1, -1
}

func (b *TextBlock) Size() int {
	return len(b.lines)
}

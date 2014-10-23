package gomez

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
	"os"
)

var (
	White = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	Black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	Red   = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	Green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	Blue  = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
	Pink  = color.RGBA{0xF9, 0x60, 0x87, 0xFF}
	Mint  = color.RGBA{0xCC, 0xFF, 0x99, 0xFF}
	Teal  = color.RGBA{0x33, 0xD5, 0xCC, 0xFF}
)

type Point struct {
	X, Y int
}

type Maze struct {
	image.Paletted
	start, end, cur Point
	// important colors
	white, black, red, green, blue uint8
	// highlighting colors
	pink, mint, teal uint8
}

func (m *Maze) atPoint(p Point) uint8 {
	return m.ColorIndexAt(p.X, p.Y)
}

func (m *Maze) setPoint(p Point, c uint8) {
	m.SetColorIndex(p.X, p.Y, c)
}

func (m *Maze) Solve() bool {
	return m.recSolve(m.options(m.start)[0])
}

func (m *Maze) recSolve(p Point) bool {
	// color current point
	m.setPoint(p, m.red)

	if m.atEnd(p) {
		return true
	}

	// determine where to go next
	for _, route := range m.options(p) {
		if m.recSolve(route) {
			return true
		}
	}
	m.setPoint(p, m.mint)
	return false
}

// search for dead ends, then follow them back to their first branching point.
// this is less efficient than recSolve, but arguably more fun to watch.
func (m *Maze) deadEndSolve() bool {
	for i := range m.Pix {
		p := Point{i % m.Stride, i / m.Stride}
		if m.atPoint(p) == m.white {
			for len(m.options(p)) == 1 && !m.atStart(p) && !m.atEnd(p) {
				m.setPoint(p, m.mint)
				p = m.options(p)[0]
			}
		}
	}
	// the remaining path is the solution
	for i := range m.Pix {
		if m.Pix[i] == m.white {
			m.Pix[i] = m.red
		}
	}
	return true
}

// brute-force the maze, changing colors at each branch
func (m *Maze) ColorRoutes() {
	m.recColorRoute(m.options(m.start)[0], []uint8{m.pink, m.mint, m.teal})
}

func (m *Maze) recColorRoute(p Point, colors []uint8) {
	// color current point
	m.setPoint(p, colors[0])

	// determine where to go next
	for _, route := range m.options(p) {
		m.recColorRoute(route, colors)
		colors = append(colors[1:], colors[0])
	}
}

// return a list of untraveled (white) directions
func (m *Maze) options(p Point) (ps []Point) {
	adj := []Point{
		{p.X + 0, p.Y - 1},
		{p.X - 1, p.Y + 0}, {p.X + 1, p.Y + 0},
		{p.X + 0, p.Y + 1},
	}
	for _, a := range adj {
		if a.X >= 0 && a.Y >= 0 && m.atPoint(a) == m.white {
			ps = append(ps, a)
		}
	}
	return
}

// return true if we are adjacent to the start point
func (m *Maze) atStart(p Point) bool {
	adj := []Point{
		{p.X + 0, p.Y - 1},
		{p.X - 1, p.Y + 0}, {p.X + 1, p.Y + 0},
		{p.X + 0, p.Y + 1},
	}
	for _, a := range adj {
		if m.atPoint(a) == m.green {
			return true
		}
	}
	return false
}

// return true if we are adjacent to the end point
func (m *Maze) atEnd(p Point) bool {
	adj := []Point{
		{p.X + 0, p.Y - 1},
		{p.X - 1, p.Y + 0}, {p.X + 1, p.Y + 0},
		{p.X + 0, p.Y + 1},
	}
	for _, a := range adj {
		if m.atPoint(a) == m.blue {
			return true
		}
	}
	return false
}

func (m *Maze) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return gif.Encode(f, &m.Paletted, nil) // hopefully...
}

func New(filename string) (m Maze, err error) {
	// read file
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	rawImage, err := gif.Decode(f)
	if err != nil {
		return
	}
	// underlying representation is image.Paletted, which is easier to work with
	m.Paletted = *rawImage.(*image.Paletted)

	// identify colors
	m.white = uint8(m.Palette.Index(White))
	m.black = uint8(m.Palette.Index(Black))
	m.red = uint8(m.Palette.Index(Red))
	m.green = uint8(m.Palette.Index(Green))
	m.blue = uint8(m.Palette.Index(Blue))
	m.pink = uint8(m.Palette.Index(Pink))
	m.mint = uint8(m.Palette.Index(Mint))
	m.teal = uint8(m.Palette.Index(Teal))

	// locate start/end
	for i := range m.Pix {
		switch m.Pix[i] {
		case m.green:
			m.start = Point{i % m.Stride, i / m.Stride}
		case m.blue:
			m.end = Point{i % m.Stride, i / m.Stride}
		}
	}
	if m.start == (Point{0, 0}) {
		err = errors.New("maze is missing start pixel")
	} else if m.end == (Point{0, 0}) {
		err = errors.New("maze is missing end pixel")
	}

	return
}

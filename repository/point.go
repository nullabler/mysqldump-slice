package repository

type PointInterface interface {
	Next(int) (string, int)
	Key(string, int) string
	Count() int
	Current() int
}

type Point struct {
	row     []int
	length  []int
	count   int
	rowLen  int
	match   map[string]int
	current int

	keys map[string][]string
}

func NewPoint(keys map[string][]string) *Point {
	p := &Point{
		match: make(map[string]int),
		keys:  keys,
	}

	n := 0
	for col, list := range p.keys {
		p.match[col] = n
		p.row = append(p.row, 0)

		length := len(list)
		p.length = append(p.length, length)
		if p.count == 0 {
			p.count = length
		} else {
			p.count += length
		}
		n++
	}
	p.rowLen = len(p.row)

	return p
}

func (p *Point) Next(n int) (string, int) {
	line := p.current
	p.current++
	if p.current >= p.rowLen {
		p.current = 0
	}

	for r, c := range p.row {
		if r == line {
			if n+1 < p.count && line == p.rowLen-1 {
				p.up(line)
			}

			return p.col(r), c
		}
	}
	return "", 0
}

func (p *Point) Key(k string, i int) string {
	return p.keys[k][i]
}

func (p *Point) Count() int {
	return p.count
}

func (p *Point) Current() int {
	return p.current
}

func (p *Point) up(line int) {
	if p.row[line] == p.length[line]-1 {
		p.row[line] = 0
		if line == p.rowLen-1 {
			p.up(0)
		} else {
			p.up(line + 1)
		}
	} else {
		p.row[line]++
	}
}

func (p *Point) col(line int) string {
	for c, l := range p.match {
		if l == line {
			return c
		}
	}

	return ""
}

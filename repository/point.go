package repository

type Point struct {
	row     []int
	length  []int
	Count   int
	rowLen  int
	match   map[string]int
	Current int

	Keys map[string][]string
}

func NewPoint(keys map[string][]string) *Point {
	p := &Point{
		match: make(map[string]int),
		Keys:  keys,
	}

	n := 0
	for col, list := range p.Keys {
		p.match[col] = n
		p.row = append(p.row, 0)

		length := len(list)
		p.length = append(p.length, length)
		if p.Count == 0 {
			p.Count = length
		} else {
			p.Count += length
		}
		n++
	}
	p.rowLen = len(p.row)

	return p
}

func (p *Point) Next(n int) (string, int) {
	line := p.Current
	p.Current++
	if p.Current >= p.rowLen {
		p.Current = 0
	}

	for r, c := range p.row {
		if r == line {
			if n+1 < p.Count && line == p.rowLen - 1 {
				p.up(line)
			}

			return p.col(r), c
		}
	}
	return "", 0
}

func (p *Point) up(line int) {
	if p.row[line] == p.length[line]-1 {
		p.row[line] = 0
		if line == p.rowLen - 1 {
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

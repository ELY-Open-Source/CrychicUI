package code

import (
	"fmt"
	"strconv"
)

// eval evaluates a small arithmetic expression.
func eval(src string) (float64, error) {
	p := &arith{s: src}
	v, err := p.expr()
	if err != nil {
		return 0, err
	}
	p.space()
	if p.i < len(p.s) {
		return 0, fmt.Errorf("unexpected %q", p.s[p.i:])
	}
	return v, nil
}

type arith struct {
	s string
	i int
}

func (p *arith) space() {
	for p.i < len(p.s) && (p.s[p.i] == ' ' || p.s[p.i] == '\t') {
		p.i++
	}
}

func (p *arith) expr() (float64, error) {
	v, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		p.space()
		if p.i >= len(p.s) || (p.s[p.i] != '+' && p.s[p.i] != '-') {
			return v, nil
		}
		op := p.s[p.i]
		p.i++
		r, err := p.term()
		if err != nil {
			return 0, err
		}
		if op == '+' {
			v += r
		} else {
			v -= r
		}
	}
}

func (p *arith) term() (float64, error) {
	v, err := p.factor()
	if err != nil {
		return 0, err
	}
	for {
		p.space()
		if p.i >= len(p.s) || (p.s[p.i] != '*' && p.s[p.i] != '/') {
			return v, nil
		}
		op := p.s[p.i]
		p.i++
		r, err := p.factor()
		if err != nil {
			return 0, err
		}
		if op == '*' {
			v *= r
		} else {
			if r == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			v /= r
		}
	}
}

func (p *arith) factor() (float64, error) {
	p.space()
	if p.i >= len(p.s) {
		return 0, fmt.Errorf("missing a value")
	}
	switch p.s[p.i] {
	case '-':
		p.i++
		v, err := p.factor()
		return -v, err
	case '+':
		p.i++
		return p.factor()
	case '(':
		p.i++
		v, err := p.expr()
		if err != nil {
			return 0, err
		}
		p.space()
		if p.i >= len(p.s) || p.s[p.i] != ')' {
			return 0, fmt.Errorf("missing a closing paren")
		}
		p.i++
		return v, nil
	}
	start := p.i
	for p.i < len(p.s) && (p.s[p.i] >= '0' && p.s[p.i] <= '9' || p.s[p.i] == '.') {
		p.i++
	}
	if start == p.i {
		return 0, fmt.Errorf("missing a value")
	}
	v, err := strconv.ParseFloat(p.s[start:p.i], 64)
	if err != nil {
		return 0, fmt.Errorf("bad number %q", p.s[start:p.i])
	}
	return v, nil
}

// evalText formats an arithmetic result for the REPL.
func evalText(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

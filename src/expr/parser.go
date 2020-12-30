package main

import (
	"log"
)

type Parser struct {
	src string
	idx int
}

type Ast struct {
	lhs  uint64
	oper string
	rhs  uint64
}

const (
	TokenEOF     = "EOF"
	TokenInvalid = "INVALID"
	TokenNumber  = "NUMBER"
)

type Token struct {
	kind string
	ival uint64
}

func (p *Parser) EOF() bool {
	return p.idx == len(p.src)
}

func (p *Parser) Current() byte {
	if p.idx < len(p.src) {
		return p.src[p.idx]
	} else {
		return 0
	}
}

func (p *Parser) Next() byte {
	if p.idx < len(p.src) {
		p.idx++
	}

	return p.Current()
}

func (p *Parser) ReadToken() Token {
	// skip whitespace
	s := p.Current()
	for s <= ' ' && !p.EOF() {
		s = p.Next()
	}

	if p.EOF() {
		return Token{kind: TokenEOF}
	}

	is := func(a, b byte) bool { return a == b }
	in := func(a, b, c byte) bool { return a >= b && a <= c }

	switch {
	case is(s, '+'):
		p.Next()

		return Token{kind: "+"}
	case is(s, '-'):
		p.Next()

		return Token{kind: "-"}
	case in(s, '0', '9'):
		t := Token{kind: TokenNumber}

		// collect digits
		for in(s, '0', '9') && !p.EOF() {
			t.ival = t.ival*10 + uint64(s-'0')

			s = p.Next()
		}

		return t
	}

	return Token{kind: TokenInvalid}
}

func (p *Parser) ParseExpression() Ast {
	// default values, so we can always return *something*
	ast := Ast{
		lhs:  0,
		oper: "+",
		rhs:  0,
	}

	lhs := p.ReadToken()
	if lhs.kind != TokenNumber {
		log.Printf("expected number, got %v\n", lhs.kind)

		return ast
	}

	ast.lhs = lhs.ival

	op := p.ReadToken()
	if op.kind != "+" && op.kind != "-" {
		log.Printf("expected operator after %d, got %v\n", lhs.ival, op.kind)

		return ast
	}

	ast.oper = op.kind

	rhs := p.ReadToken()
	if rhs.kind != TokenNumber {
		log.Printf("expected number after %v, got %v\n", op.kind, rhs.kind)

		return ast
	}

	ast.rhs = rhs.ival

	return ast
}

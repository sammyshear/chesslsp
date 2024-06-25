package analysis

import "fmt"

type parser struct {
	lex       *lexer
	pos       int
	token     token //current token
	lastToken token //previous token
}

type ParseError struct {
	Line      int
	Character int
	Message   string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%d:%d:%s", e.Line, e.Character, e.Message)
}

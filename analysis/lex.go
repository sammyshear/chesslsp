package analysis

import (
	"fmt"
	"log"
	"regexp"
)

type TokenKind int

const (
	EOF TokenKind = iota
	STRING
	INTEGER
	PERIOD
	ASTERISK
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_ANGLE_BRACKET
	RIGHT_ANGLE_BRACKET
	LEFT_CURLY_BRACKET
	RIGHT_CURLY_BRACKET
	NAG
	SYMBOL
	RESULT
)

type token struct {
	Value string
	Kind  TokenKind
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	Tokens   []token
	patterns []regexPattern
	source   string
	pos      int
}

func newToken(kind TokenKind, value string) token {
	return token{Kind: kind, Value: value}
}

func Tokenize(source string) []token {
	lex := createLexer(source)

	for !lex.atEOF() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer error: unrecognized token near: %v\n", lex.remainder()))
		}
	}

	lex.push(newToken(EOF, ""))
	return lex.Tokens
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		source: source,
		Tokens: []token{},
		patterns: []regexPattern{
			{regexp.MustCompile(`(1-0|0-1|1/2-1/2)`), resultHandler},
			{regexp.MustCompile(`[0-9]+`), intHandler},
			{regexp.MustCompile(`^$[0-9]+`), nagHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9#+=-]+`), symbolHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\[`), defaultHandler(LEFT_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(RIGHT_BRACKET, "]")},
			{regexp.MustCompile("<"), defaultHandler(LEFT_ANGLE_BRACKET, "<")},
			{regexp.MustCompile(">"), defaultHandler(RIGHT_ANGLE_BRACKET, ">")},
			{regexp.MustCompile("{"), defaultHandler(LEFT_CURLY_BRACKET, "}")},
			{regexp.MustCompile("}"), defaultHandler(RIGHT_CURLY_BRACKET, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(LEFT_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(RIGHT_PAREN, ")")},
			{regexp.MustCompile(`\.`), defaultHandler(PERIOD, ".")},
			{regexp.MustCompile(`\*`), defaultHandler(ASTERISK, "*")},
		},
	}
}

func (lex *lexer) advance(n int) {
	lex.pos += n
}

func (lex *lexer) push(token token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) atEOF() bool {
	return lex.pos >= len(lex.source)
}

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) at() byte {
	return lex.source[lex.pos]
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(newToken(kind, value))
	}
}

func nagHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(NAG, match))
	lex.advance(len(match))
}

func intHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(INTEGER, match))
	lex.advance(len(match))
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	log.Print(lex.remainder())
	stringLiteral := lex.remainder()[match[0]+1 : match[1]-1]
	lex.push(newToken(STRING, stringLiteral))
	lex.advance(len(stringLiteral) + 2)
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(SYMBOL, match))
	lex.advance(len(match))
}

func resultHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(RESULT, match))
	lex.advance(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.advance(len(match))
}

package analysis

import (
	"reflect"
	"testing"
)

type lexTest struct {
	name   string
	input  string
	tokens []token
}

var (
	tEOF = token{"", EOF}
)

var lexTests = []lexTest{
	{"empty", "", []token{tEOF}},
	{"spaces", "\t\r", []token{tEOF}},
	{"tag", `[Event "test game"]`, []token{
		{"[", LEFT_BRACKET},
		{"Event", SYMBOL},
		{"test game", STRING},
		{"]", RIGHT_BRACKET},
		tEOF,
	}},
	{"moves", "12. O-O-O Bxe5+ (12...e8=Q)", []token{
		{"12", INTEGER},
		{".", PERIOD},
		{"O-O-O", SYMBOL},
		{"Bxe5+", SYMBOL},
		{"(", LEFT_PAREN},
		{"12", INTEGER},
		{".", PERIOD},
		{".", PERIOD},
		{".", PERIOD},
		{"e8=Q", SYMBOL},
		{")", RIGHT_PAREN},
		tEOF,
	}},
	{"results", "1-0 0-1 1/2-1/2", []token{
		{"1-0", RESULT},
		{"0-1", RESULT},
		{"1/2-1/2", RESULT},
		tEOF,
	}},
}

func collect(t *lexTest) []token {
	return Tokenize(t.input)
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		tokens := collect(&test)
		if !reflect.DeepEqual(tokens, test.tokens) {
			t.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}

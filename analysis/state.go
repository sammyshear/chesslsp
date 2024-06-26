package analysis

import (
	"bytes"
	"chesslsp/lsp"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/malbrecht/chess/pgn"
)

type State struct {
	Documents map[string][]byte
	lineStart []int //byte offset of each line
	lineOnce  sync.Once
	DB        *pgn.DB
}

func NewState() State {
	return State{Documents: map[string][]byte{}, DB: &pgn.DB{}}
}

func (s *State) initLines(uri string) {
	document := s.Documents[uri]
	s.lineOnce.Do(func() {
		nlines := bytes.Count(document, []byte("\n"))
		s.lineStart = make([]int, 1, nlines+1)
		for offset, b := range document {
			if b == '\n' {
				s.lineStart = append(s.lineStart, offset+1)
			}
		}
	})
}

func rangeOffset(r lsp.Range, uri string, s *State) (int, int, error) {
	start, err := positionOffset(r.Start, uri, s)
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding start offset: %s", err)
	}
	end, err := positionOffset(r.End, uri, s)
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding end offset: %s", err)
	}

	return start, end, nil
}

func positionOffset(p lsp.Position, uri string, s *State) (int, error) {
	s.initLines(uri)
	document := s.Documents[uri]
	if p.Line > len(s.lineStart) {
		return 0, fmt.Errorf("line number %d out of range of 0-%d", p.Line, len(s.lineStart))
	} else if p.Line == len(s.lineStart) {
		if p.Character == 0 {
			return len(document), nil
		}
		return 0, fmt.Errorf("column is beyond end of file")
	}

	offset := s.lineStart[p.Line]
	content := document[offset:]

	col8 := 0
	for col16 := 0; col16 < p.Character; col16++ {
		r, sz := utf8.DecodeRune(content)
		if sz == 0 {
			return 0, fmt.Errorf("column is beyond end of file")
		}
		if r == '\n' {
			return 0, fmt.Errorf("column is beyond end of line")
		}
		if sz == 1 && r == utf8.RuneError {
			return 0, fmt.Errorf("buffer contains invalid utf8 text")
		}
		content = content[sz:]
		if r >= 0x10000 {
			col16++

			if col16 == int(p.Character) {
				break
			}
		}
		col8 += sz
	}

	return offset + col8, nil
}

func getDiagnostics(text string, s *State) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	s.DB = &pgn.DB{}
	errs := s.DB.Parse(text)
	if errs != nil {
		for _, err := range errs {
			split := strings.Split(strings.Replace(err.Error(), " ", "", 1), ":")
			line, _ := strconv.Atoi(split[0])
			col, _ := strconv.Atoi(split[1])
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      line,
						Character: col,
					},
					End: lsp.Position{
						Line:      line,
						Character: col,
					},
				},
				Severity: 1,
				Source:   "chesslsp",
				Message:  split[2],
			})
		}
		return diagnostics
	}
	err := s.DB.ParseMoves(s.DB.Games[0])
	if err != nil {
		split := strings.Split(strings.Replace(err.Error(), " ", "", 1), ":")
		line, _ := strconv.Atoi(split[0])
		col, _ := strconv.Atoi(split[1])
		diagnostics = append(diagnostics, lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      line,
					Character: col,
				},
				End: lsp.Position{
					Line:      line,
					Character: col,
				},
			},
			Severity: 1,
			Source:   "chesslsp",
			Message:  split[2],
		})
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = []byte(text)

	return getDiagnostics(text, s)
}

func (s *State) UpdateDocument(uri string, contentChanges []lsp.TextDocumentChangeEvent) []lsp.Diagnostic {
	for _, change := range contentChanges {
		start, end, err := rangeOffset(change.Range, uri, s)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		buf.Write(s.Documents[uri][:start])
		buf.WriteString(change.Text)
		buf.Write(s.Documents[uri][end:])
		s.Documents[uri] = buf.Bytes()
	}

	return getDiagnostics(string(s.Documents[uri]), s)
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	document := s.Documents[uri]
	items := []lsp.CompletionItem{}
	for row, line := range strings.Split(string(document), "\n") {
		_ = row
		_ = line
	}

	_ = items
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: items,
	}

	return response

}

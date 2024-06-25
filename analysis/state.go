package analysis

import (
	"chesslsp/lsp"
	"strings"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnostics(text string, s *State) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	_ = s
	_ = text

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnostics(text, s)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnostics(text, s)
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	document := s.Documents[uri]
	items := []lsp.CompletionItem{}
	for row, line := range strings.Split(document, "\n") {
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

package main

import (
	"bufio"
	"chesslsp/analysis"
	"chesslsp/lsp"
	"chesslsp/rpc"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	w := os.Stdout
	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			log.Printf("Error decoding message: %s", err)
		}

		handleMessage(w, state, method, contents)
	}
}

func handleMessage(w io.Writer, state analysis.State, method string, contents []byte) {
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			return
		}
		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(w, msg)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			return
		}
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(w, lsp.PublishDiagnosticNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
	case "textDocument/didChange":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			return
		}
		diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(w, lsp.PublishDiagnosticNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
	}
}

func writeResponse(w io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	w.Write([]byte(reply))
}

package lsp

type ServerCapabilities struct {
	TextDocumentSync   int            `json:"textDocumentSync"`
	CompletionProvider map[string]any `json:"completionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   2,
				CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "chesslsp",
				Version: "0.1.0-beta1",
			},
		},
	}
}

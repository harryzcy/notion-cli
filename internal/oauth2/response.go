package oauth2

type codeResponse struct {
	code  string
	state string
	err   string
}

type tokenResponse struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceIcon string `json:"workspace_icon"`
	BotID         string `json:"bot_id"`
	Owner         owner  `json:"owner"`
}

type owner struct {
	Workspace bool `json:"workspace"`

	Object    string `json:"object"` // user
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`

	Bot bool `json:"bot"`
}

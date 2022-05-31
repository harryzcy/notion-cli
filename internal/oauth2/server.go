package oauth2

import (
	"net/http"
)

func getCallbackHandler(c chan codeResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := codeResponse{
			code:  r.URL.Query().Get("code"),
			state: r.URL.Query().Get("state"),
			err:   r.URL.Query().Get("error"),
		}
		c <- res
	}
}

func startServer(c chan codeResponse) {
	http.HandleFunc("/callback", getCallbackHandler(c))
	http.ListenAndServe(":"+port, nil)
}

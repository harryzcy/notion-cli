package oauth2

import (
	"fmt"
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

func startServer(c chan codeResponse) error {
	http.HandleFunc("/callback", getCallbackHandler(c))
	fmt.Println(21)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(23)
	return err
}

package oauth2

func buildAuthURL(state string) string {
	params := map[string]string{
		"client_id":     clientID,
		"redirect_uri":  "http://localhost:" + port + "/callback",
		"response_type": "code",
		"owner":         "user",
		"state":         state,
	}
	return baseURL + "?" + buildQuery(params)
}

func buildQuery(params map[string]string) string {
	var query string
	for key, value := range params {
		query += key + "=" + value + "&"
	}
	return query[:len(query)-1]
}

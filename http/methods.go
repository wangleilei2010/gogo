package http

func Post(url, body string) (*Response, error) {
	session := NewSession()
	return session.request(url, "POST", body)
}

func Get(url string) (*Response, error) {
	session := NewSession()
	return session.request(url, "GET", "")
}

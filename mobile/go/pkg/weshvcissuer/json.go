package weshvcissuer

type JSONError struct {
	Error string `json:"error"`
}

type JSONChallenge struct {
	Challenge string `json:"challenge"`
}

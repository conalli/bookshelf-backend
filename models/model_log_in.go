package models

// LogInRes represents the data sent back to the client when the user logs in.
type LogInRes struct {
	APIKey string `json:"apiKey"`
}

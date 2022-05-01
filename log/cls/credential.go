package cls

// Credential is the ticket for accessing the log service at QCloud
type Credential struct {
	SecretID  string `json:"SecretID"`
	SecretKey string `json:"SecretKey"`
}

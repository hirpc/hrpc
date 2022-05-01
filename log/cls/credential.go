package cls

// Credential is the ticket for accessing the log service at QCloud
type Credential struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

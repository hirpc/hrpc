package dgraph

type Option struct {
	Targets    []string   `json:"targets"`
	Credential credential `json:"credential"`
}

type credential struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	Namespace uint64 `json:"namespace"`
}

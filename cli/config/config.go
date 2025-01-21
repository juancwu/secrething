package config

var backendUrl string

func Init() {
	if backendUrl == "" {
		backendUrl = "http://localhost:3000/api/v1"
	}
}

func BackendUrl() string {
	return backendUrl
}

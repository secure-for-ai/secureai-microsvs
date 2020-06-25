package util

type AppInfo struct {
	Env string `json:"Env"` // example: local, dev, prod
}

const (
	AppEnvDev  = "dev"
	AppEnvTest = "test"
	AppEnvProd = "prod"
)

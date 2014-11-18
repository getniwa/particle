package spark

const (
	APIVersion        = "/v1"
	BaseUrl           = "https://api.spark.io"
	MaxVariableLen    = 12
	BasicAuthId       = "spark"
	BasicAuthPassword = "spark"
)

type APIUrl struct {
	BaseUrl    string
	APIVersion string
	Endpoint   string
}

func Endpoint(a *APIUrl) string {
	return a.BaseUrl + a.APIVersion + a.Endpoint
}

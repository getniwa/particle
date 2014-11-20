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

func (a APIUrl) String() string {
	return a.BaseUrl + a.APIVersion + a.Endpoint
}

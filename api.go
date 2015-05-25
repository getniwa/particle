package particle

const (
	APIVersion        = "/v1"
	BaseUrl           = "https://api.particle.io"
	MaxVariableLen    = 12
	BasicAuthId       = "particle"
	BasicAuthPassword = "particle"
)

type APIUrl struct {
	BaseUrl    string
	APIVersion string
	Endpoint   string
}

func (a APIUrl) String() string {
	return a.BaseUrl + a.APIVersion + a.Endpoint
}

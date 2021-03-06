package particle

type AuthToken interface {
	Valid() error
	String() string
}

type AuthTokenProvider interface {
	AuthToken() (AuthToken, error)
}

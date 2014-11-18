package spark

type AuthToken interface {
	Valid() error
	Token() (string, error)
}

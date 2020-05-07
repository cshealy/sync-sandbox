package data

type DAO struct {
	BearerToken string
}

type ExternalAPI interface {
	getToken() error
}

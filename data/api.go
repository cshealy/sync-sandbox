package data

import "time"

type DAO struct {
	BearerToken     string
	TokenExpiration time.Time
}

type ExternalAPI interface {
	getToken() error
}

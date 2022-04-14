package repository

import (
	"io/ioutil"
	"net/http"

	"github.com/swissvortex/cosmos-opt-api/logger"
)

type Repository interface {
	HttpGetBody(uri string) []byte
}

type repository struct {
	log logger.Logger
}

func NewRepository(log logger.Logger) Repository {
	return &repository{
		log: log,
	}
}

func (r *repository) HttpGetBody(uri string) []byte {
	response, err := http.Get(uri)
	if err != nil {
		r.log.ErrorWithContext(r.log.FileContext(), err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		r.log.ErrorWithContext(r.log.FileContext(), err)
	}
	return responseData
}

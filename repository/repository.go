package repository

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Colm3na/cosmos-opt-api/logger"
)

type Repository interface {
	HttpGetBody(uri string) []byte
}

type repository struct {
	logger *logger.Logger
}

func NewRepository(logger *logger.Logger) Repository {
	return &repository{
		logger: logger,
	}
}

func (r *repository) HttpGetBody(uri string) []byte {
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

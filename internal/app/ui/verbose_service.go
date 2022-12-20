package ui

import (
	"fmt"

	"github.com/kerelape/urlshortener/internal/app/model"
)

type VerboseService struct {
	Origin Service
	Name   string
	Log    model.Log
}

func NewVerboseService(origin Service, name string, log model.Log) *VerboseService {
	return &VerboseService{
		Origin: origin,
		Name:   name,
		Log:    log,
	}
}

func (service *VerboseService) Execute() error {
	service.Log.WriteInfo(fmt.Sprintf("Starting service: %s", service.Name))
	var err = service.Origin.Execute()
	if err != nil {
		service.Log.WriteFailure(fmt.Sprintf("Service %s failed", service.Name))
	} else {
		service.Log.WriteInfo(fmt.Sprintf("Service %s shut down", service.Name))
	}
	return err
}

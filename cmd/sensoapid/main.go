/*
 * Senso-Care
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"fmt"
	openapi "github.com/Senso-Care/SensoApi/internal/api"
	"github.com/Senso-Care/SensoApi/internal/config"
	"github.com/Senso-Care/SensoApi/internal/data"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
)

func main() {
	log.SetLevel(log.DebugLevel)

	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading configuration: %s\n", err)
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"DatabaseURI": configuration.Database.ConnectionUri,
	}).Info("Server address loaded from configuration")
	log.WithField("GOMAXPROCS", configuration.Cores).Info("Setting max number of cpus")
	log.WithField("PORT", configuration.Port).Info("Setting port to listen to")
	runtime.GOMAXPROCS(configuration.Cores)
	log.Printf("Server started")
	var service data.InfluxServicer
	if configuration.Mock {
		log.Println("Server is a MOCK, only false data generated")
		service = data.NewMockService(&configuration.Database)
	} else {
		service = data.NewInfluxService(&configuration.Database)

	}
	DefaultApiService := openapi.NewDefaultApiService(&service)
	DefaultApiController := openapi.NewDefaultApiController(DefaultApiService)

	router := openapi.NewRouter(DefaultApiController)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), router))
}

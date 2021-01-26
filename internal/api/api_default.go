/*
 * Senso-Care
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"encoding/json"
	"github.com/Senso-Care/SensoApi/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A DefaultApiController binds http requests to an api service and writes the service results to the http response
type DefaultApiController struct {
	service DefaultApiServicer
}

// NewDefaultApiController creates a default api controller
func NewDefaultApiController(s DefaultApiServicer) Router {
	return &DefaultApiController{ service: s }
}

// Routes returns all of the api route for the DefaultApiController
func (c *DefaultApiController) Routes() Routes {
	return Routes{
		{
			"GetLastMetrics",
			strings.ToUpper("Get"),
			"/metrics/{type}/last",
			c.GetLastMetrics,
		},
		{
			"GetMetrics",
			strings.ToUpper("Get"),
			"/metrics",
			c.GetMetrics,
		},
		{
			"GetMetricsFromSensor",
			strings.ToUpper("Get"),
			"/sensors/{name}",
			c.GetMetricsFromSensor,
		},
		{
			"GetMetricsFromType",
			strings.ToUpper("Get"),
			"/metrics/{type}",
			c.GetMetricsFromType,
		},
		{
			"GetSensors",
			strings.ToUpper("Get"),
			"/sensors",
			c.GetSensors,
		},
		{
			"PostMetricsFromType",
			strings.ToUpper("Post"),
			"/metrics/{type}",
			c.PostMetricsFromType,
		},
	}
}

// GetLastMetrics - Get last value of all sensors of a given metric
func (c *DefaultApiController) GetLastMetrics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	type_ := params["type"]
	range_ := query.Get("range")
	result, err := c.service.GetLastMetrics(r.Context(), type_, range_)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetMetrics - Get list of metrics types
func (c *DefaultApiController) GetMetrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	range_ := query.Get("range")
	result, err := c.service.GetMetrics(r.Context(), range_)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetMetricsFromSensor - Get data from sensor
func (c *DefaultApiController) GetMetricsFromSensor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	name := params["name"]
	range_ := query.Get("range")
	result, err := c.service.GetMetricsFromSensor(r.Context(), name, range_)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetMetricsFromType - Get data from type
func (c *DefaultApiController) GetMetricsFromType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	type_ := params["type"]
	range_ := query.Get("range")
	result, err := c.service.GetMetricsFromType(r.Context(), type_, range_)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetSensors - Get list of sensors
func (c *DefaultApiController) GetSensors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	range_ := query.Get("range")
	result, err := c.service.GetSensors(r.Context(), range_)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PostMetricsFromType - Get data from type
func (c *DefaultApiController) PostMetricsFromType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	type_ := params["type"]
	dataPoint := &models.DataPoint{}
	if err := json.NewDecoder(r.Body).Decode(&dataPoint); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := c.service.PostMetricsFromType(r.Context(), type_, *dataPoint)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
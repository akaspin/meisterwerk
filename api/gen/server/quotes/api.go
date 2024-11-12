// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Quotes API
 *
 * API for quotes
 *
 * API version: 1.0.0
 */

package quotes

import (
	"context"
	"net/http"
)

// DefaultAPIRouter defines the required methods for binding the api requests to a responses for the DefaultAPI
// The DefaultAPIRouter implementation should parse necessary information from the http request,
// pass the data to a DefaultAPIServicer to perform the required actions, then write the service results to the http response.
type DefaultAPIRouter interface {
	ListQuotes(http.ResponseWriter, *http.Request)
	CreateQuote(http.ResponseWriter, *http.Request)
	GetQuote(http.ResponseWriter, *http.Request)
	UpdateQuote(http.ResponseWriter, *http.Request)
	DeleteQuote(http.ResponseWriter, *http.Request)
	UpdateQuotes(http.ResponseWriter, *http.Request)
	CreateQuotes(http.ResponseWriter, *http.Request)
	DeleteQuotes(http.ResponseWriter, *http.Request)
}

// DefaultAPIServicer defines the api actions for the DefaultAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DefaultAPIServicer interface {
	ListQuotes(context.Context, []int32, []int32, int32, int32, string) (ImplResponse, error)
	CreateQuote(context.Context, Quote) (ImplResponse, error)
	GetQuote(context.Context, int32) (ImplResponse, error)
	UpdateQuote(context.Context, int32, Quote) (ImplResponse, error)
	DeleteQuote(context.Context, int32) (ImplResponse, error)
	UpdateQuotes(context.Context, []Quote) (ImplResponse, error)
	CreateQuotes(context.Context, []Quote) (ImplResponse, error)
	DeleteQuotes(context.Context, QuotesIds) (ImplResponse, error)
}

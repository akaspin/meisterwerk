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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// DefaultAPIController binds http requests to an api service and writes the service results to the http response
type DefaultAPIController struct {
	service      DefaultAPIServicer
	errorHandler ErrorHandler
}

// DefaultAPIOption for how the controller is set up.
type DefaultAPIOption func(*DefaultAPIController)

// WithDefaultAPIErrorHandler inject ErrorHandler into controller
func WithDefaultAPIErrorHandler(h ErrorHandler) DefaultAPIOption {
	return func(c *DefaultAPIController) {
		c.errorHandler = h
	}
}

// NewDefaultAPIController creates a default api controller
func NewDefaultAPIController(s DefaultAPIServicer, opts ...DefaultAPIOption) *DefaultAPIController {
	controller := &DefaultAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the DefaultAPIController
func (c *DefaultAPIController) Routes() Routes {
	return Routes{
		"ListQuotes": Route{
			strings.ToUpper("Get"),
			"/quotes",
			c.ListQuotes,
		},
		"CreateQuote": Route{
			strings.ToUpper("Post"),
			"/quotes",
			c.CreateQuote,
		},
		"GetQuote": Route{
			strings.ToUpper("Get"),
			"/quotes/{id}",
			c.GetQuote,
		},
		"UpdateQuote": Route{
			strings.ToUpper("Put"),
			"/quotes/{id}",
			c.UpdateQuote,
		},
		"DeleteQuote": Route{
			strings.ToUpper("Delete"),
			"/quotes/{id}",
			c.DeleteQuote,
		},
		"UpdateQuotes": Route{
			strings.ToUpper("Put"),
			"/bulk/quotes",
			c.UpdateQuotes,
		},
		"CreateQuotes": Route{
			strings.ToUpper("Post"),
			"/bulk/quotes",
			c.CreateQuotes,
		},
		"DeleteQuotes": Route{
			strings.ToUpper("Delete"),
			"/bulk/quotes",
			c.DeleteQuotes,
		},
	}
}

// ListQuotes - request multiple quotes
func (c *DefaultAPIController) ListQuotes(w http.ResponseWriter, r *http.Request) {
	query, err := parseQuery(r.URL.RawQuery)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	idParam, err := parseNumericArrayParameter[int32](
		query.Get("id"), ",", false,
		WithParse[int32](parseInt32),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Param: "id", Err: err}, nil)
		return
	}
	customerIdParam, err := parseNumericArrayParameter[int32](
		query.Get("customer_id"), ",", false,
		WithParse[int32](parseInt32),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Param: "customer_id", Err: err}, nil)
		return
	}
	var limitParam int32
	if query.Has("limit") {
		param, err := parseNumericParameter[int32](
			query.Get("limit"),
			WithParse[int32](parseInt32),
		)
		if err != nil {
			c.errorHandler(w, r, &ParsingError{Param: "limit", Err: err}, nil)
			return
		}

		limitParam = param
	} else {
	}
	var skipParam int32
	if query.Has("skip") {
		param, err := parseNumericParameter[int32](
			query.Get("skip"),
			WithParse[int32](parseInt32),
		)
		if err != nil {
			c.errorHandler(w, r, &ParsingError{Param: "skip", Err: err}, nil)
			return
		}

		skipParam = param
	} else {
	}
	var orderParam string
	if query.Has("order") {
		param := query.Get("order")

		orderParam = param
	} else {
	}
	result, err := c.service.ListQuotes(r.Context(), idParam, customerIdParam, limitParam, skipParam, orderParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateQuote - create quote
func (c *DefaultAPIController) CreateQuote(w http.ResponseWriter, r *http.Request) {
	quoteParam := Quote{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&quoteParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertQuoteRequired(quoteParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertQuoteConstraints(quoteParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateQuote(r.Context(), quoteParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetQuote -
func (c *DefaultAPIController) GetQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam, err := parseNumericParameter[int32](
		params["id"],
		WithRequire[int32](parseInt32),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Param: "id", Err: err}, nil)
		return
	}
	result, err := c.service.GetQuote(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateQuote - update quote
func (c *DefaultAPIController) UpdateQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam, err := parseNumericParameter[int32](
		params["id"],
		WithRequire[int32](parseInt32),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Param: "id", Err: err}, nil)
		return
	}
	quoteParam := Quote{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&quoteParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertQuoteRequired(quoteParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertQuoteConstraints(quoteParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateQuote(r.Context(), idParam, quoteParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteQuote - delete quote by id
func (c *DefaultAPIController) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam, err := parseNumericParameter[int32](
		params["id"],
		WithRequire[int32](parseInt32),
	)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Param: "id", Err: err}, nil)
		return
	}
	result, err := c.service.DeleteQuote(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateQuotes - update quotes
func (c *DefaultAPIController) UpdateQuotes(w http.ResponseWriter, r *http.Request) {
	quoteParam := []Quote{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&quoteParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	for _, el := range quoteParam {
		if err := AssertQuoteRequired(el); err != nil {
			c.errorHandler(w, r, err, nil)
			return
		}
	}
	result, err := c.service.UpdateQuotes(r.Context(), quoteParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateQuotes - create quotes
func (c *DefaultAPIController) CreateQuotes(w http.ResponseWriter, r *http.Request) {
	quoteParam := []Quote{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&quoteParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	for _, el := range quoteParam {
		if err := AssertQuoteRequired(el); err != nil {
			c.errorHandler(w, r, err, nil)
			return
		}
	}
	result, err := c.service.CreateQuotes(r.Context(), quoteParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// DeleteQuotes - delete quotes
func (c *DefaultAPIController) DeleteQuotes(w http.ResponseWriter, r *http.Request) {
	quotesIdsParam := QuotesIds{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&quotesIdsParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertQuotesIdsRequired(quotesIdsParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertQuotesIdsConstraints(quotesIdsParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.DeleteQuotes(r.Context(), quotesIdsParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

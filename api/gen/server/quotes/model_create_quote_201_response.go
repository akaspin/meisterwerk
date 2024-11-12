// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Quotes API
 *
 * API for quotes
 *
 * API version: 1.0.0
 */

package quotes

type CreateQuote201Response struct {

	// ID of created quote
	Id int64 `json:"id,omitempty"`
}

// AssertCreateQuote201ResponseRequired checks if the required fields are not zero-ed
func AssertCreateQuote201ResponseRequired(obj CreateQuote201Response) error {
	return nil
}

// AssertCreateQuote201ResponseConstraints checks if the values respects the defined constraints
func AssertCreateQuote201ResponseConstraints(obj CreateQuote201Response) error {
	return nil
}
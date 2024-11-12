/*
Quotes API

API for quotes

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package quotes

import (
	"encoding/json"
)

// checks if the CreateQuote201Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateQuote201Response{}

// CreateQuote201Response struct for CreateQuote201Response
type CreateQuote201Response struct {
	// ID of created quote
	Id *int64 `json:"id,omitempty"`
}

// NewCreateQuote201Response instantiates a new CreateQuote201Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateQuote201Response() *CreateQuote201Response {
	this := CreateQuote201Response{}
	return &this
}

// NewCreateQuote201ResponseWithDefaults instantiates a new CreateQuote201Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateQuote201ResponseWithDefaults() *CreateQuote201Response {
	this := CreateQuote201Response{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CreateQuote201Response) GetId() int64 {
	if o == nil || IsNil(o.Id) {
		var ret int64
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateQuote201Response) GetIdOk() (*int64, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CreateQuote201Response) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int64 and assigns it to the Id field.
func (o *CreateQuote201Response) SetId(v int64) {
	o.Id = &v
}

func (o CreateQuote201Response) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateQuote201Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	return toSerialize, nil
}

type NullableCreateQuote201Response struct {
	value *CreateQuote201Response
	isSet bool
}

func (v NullableCreateQuote201Response) Get() *CreateQuote201Response {
	return v.value
}

func (v *NullableCreateQuote201Response) Set(val *CreateQuote201Response) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateQuote201Response) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateQuote201Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateQuote201Response(val *CreateQuote201Response) *NullableCreateQuote201Response {
	return &NullableCreateQuote201Response{value: val, isSet: true}
}

func (v NullableCreateQuote201Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateQuote201Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
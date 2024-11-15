/*
Quotes API

API for quotes

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package quotes

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Item type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Item{}

// Item Quote item
type Item struct {
	Id      string   `json:"id"`
	Segment string   `json:"segment"`
	Price   *float32 `json:"price,omitempty"`
	Tax     *float32 `json:"tax,omitempty"`
}

type _Item Item

// NewItem instantiates a new Item object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewItem(id string, segment string) *Item {
	this := Item{}
	this.Id = id
	this.Segment = segment
	return &this
}

// NewItemWithDefaults instantiates a new Item object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewItemWithDefaults() *Item {
	this := Item{}
	return &this
}

// GetId returns the Id field value
func (o *Item) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Item) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Item) SetId(v string) {
	o.Id = v
}

// GetSegment returns the Segment field value
func (o *Item) GetSegment() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Segment
}

// GetSegmentOk returns a tuple with the Segment field value
// and a boolean to check if the value has been set.
func (o *Item) GetSegmentOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Segment, true
}

// SetSegment sets field value
func (o *Item) SetSegment(v string) {
	o.Segment = v
}

// GetPrice returns the Price field value if set, zero value otherwise.
func (o *Item) GetPrice() float32 {
	if o == nil || IsNil(o.Price) {
		var ret float32
		return ret
	}
	return *o.Price
}

// GetPriceOk returns a tuple with the Price field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Item) GetPriceOk() (*float32, bool) {
	if o == nil || IsNil(o.Price) {
		return nil, false
	}
	return o.Price, true
}

// HasPrice returns a boolean if a field has been set.
func (o *Item) HasPrice() bool {
	if o != nil && !IsNil(o.Price) {
		return true
	}

	return false
}

// SetPrice gets a reference to the given float32 and assigns it to the Price field.
func (o *Item) SetPrice(v float32) {
	o.Price = &v
}

// GetTax returns the Tax field value if set, zero value otherwise.
func (o *Item) GetTax() float32 {
	if o == nil || IsNil(o.Tax) {
		var ret float32
		return ret
	}
	return *o.Tax
}

// GetTaxOk returns a tuple with the Tax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Item) GetTaxOk() (*float32, bool) {
	if o == nil || IsNil(o.Tax) {
		return nil, false
	}
	return o.Tax, true
}

// HasTax returns a boolean if a field has been set.
func (o *Item) HasTax() bool {
	if o != nil && !IsNil(o.Tax) {
		return true
	}

	return false
}

// SetTax gets a reference to the given float32 and assigns it to the Tax field.
func (o *Item) SetTax(v float32) {
	o.Tax = &v
}

func (o Item) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Item) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["segment"] = o.Segment
	if !IsNil(o.Price) {
		toSerialize["price"] = o.Price
	}
	if !IsNil(o.Tax) {
		toSerialize["tax"] = o.Tax
	}
	return toSerialize, nil
}

func (o *Item) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"segment",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varItem := _Item{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varItem)

	if err != nil {
		return err
	}

	*o = Item(varItem)

	return err
}

type NullableItem struct {
	value *Item
	isSet bool
}

func (v NullableItem) Get() *Item {
	return v.value
}

func (v *NullableItem) Set(val *Item) {
	v.value = val
	v.isSet = true
}

func (v NullableItem) IsSet() bool {
	return v.isSet
}

func (v *NullableItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableItem(val *Item) *NullableItem {
	return &NullableItem{value: val, isSet: true}
}

func (v NullableItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

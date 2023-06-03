package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
)

type Validate struct {
	validate *validator.Validate
}
type ValidateI interface {
	ValidatePeer(peer models.Peer) []*ErrorResponse
	ValidateRestaurant(restaurant models.Restaurant) []*ErrorResponse
	ValidateEvent(event types.Event) []*ErrorResponse
	ValidateRestaurantData(data types.RestaurantData) []*ErrorResponse
}

func NewValidator(validate *validator.Validate) *Validate {
	return &Validate{validate}
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (v *Validate) getErrors(err error) []*ErrorResponse {
	var errors []*ErrorResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func (v *Validate) ValidatePeer(peer models.Peer) []*ErrorResponse {
	err := v.validate.Struct(peer)
	return v.getErrors(err)
}

func (v *Validate) ValidateRestaurant(restaurant models.Restaurant) []*ErrorResponse {
	err := v.validate.Struct(restaurant)
	return v.getErrors(err)
}

func (v *Validate) ValidateRestaurantData(data types.RestaurantData) []*ErrorResponse {
	err := v.validate.Struct(data)
	return v.getErrors(err)
}

func (v *Validate) ValidateEvent(event types.Event) []*ErrorResponse {
	err := v.validate.Struct(event)
	return v.getErrors(err)
}

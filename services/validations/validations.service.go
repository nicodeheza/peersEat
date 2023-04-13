package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/models"
)

type Validate struct {
	validate *validator.Validate
}
type ValidateI interface {
	ValidatePeer(peer models.Peer) []*ErrorResponse
	ValidateRestaurant(restaurant models.Restaurant) []*ErrorResponse
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
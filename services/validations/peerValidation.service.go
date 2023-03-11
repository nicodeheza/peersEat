package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodeheza/peersEat/models"
)

type ErrorResponse struct{
	FailedField string
	Tag string
	Value string
}

var validate = validator.New()

func ValidatePeer(peer models.Peer) []*ErrorResponse{
	var errors []*ErrorResponse
	err := validate.Struct(peer)

	if err != nil{
		for _,err := range err.(validator.ValidationErrors){
			var element ErrorResponse
			element.FailedField= err.StructNamespace()
			element.Tag= err.Tag()
			element.Value= err.Param()
			errors= append(errors, &element)
		}
	}

	return errors
}
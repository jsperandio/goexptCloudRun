package validator

import (
	"testing"

	playground "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type zipcodePayload struct {
	Zipcode string `validate:"required,len=8,number"`
}

func Test_Validator_Validate(t *testing.T) {
	t.Run("when the struct satisfies every tag, should return nil", func(t *testing.T) {
		payload := zipcodePayload{
			Zipcode: "01001000",
		}

		assert.NoError(t, New().Validate(payload))
	})

	t.Run("when a tag is not satisfied, should return the validation error", func(t *testing.T) {
		payload := zipcodePayload{
			Zipcode: "0100100",
		}

		var errs playground.ValidationErrors

		assert.ErrorAs(t, New().Validate(payload), &errs)
	})

	t.Run("when the value is not a struct, should return an invalid validation error", func(t *testing.T) {
		var invalid *playground.InvalidValidationError

		assert.ErrorAs(t, New().Validate("01001000"), &invalid)
	})
}

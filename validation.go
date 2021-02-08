package forms

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
)

var ErrMustBeString = fmt.Errorf("Field must be a string")
var ErrMustBeBool = fmt.Errorf("Field must be a boolean")
var ErrMustBeInt = fmt.Errorf("Field must be an integer")

var numberRx = regexp.MustCompile(`^-?\d*$`)
var emailRx = regexp.MustCompile(`^\S+@\S+$`)
var lettersWithSpacesRx = regexp.MustCompile(`^[- 'a-zA-ZÀ-ÖØ-öø-ÿ]+$`)
var lettersWithNumbersRx = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var usernameRx = regexp.MustCompile(`^[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*$`)
var lettersSpacesAndNumbersRx = regexp.MustCompile(`^[- 'a-zA-ZÀ-ÖØ-öø-ÿ0-9]+$`)
var urlRx = regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
var aTrueRx = regexp.MustCompile(`^true$`)

// ValidateOrChain Chains two or more validations in an or arrangement
func ValidateOrChain(validators ...Validator) Validator {
	return func(s interface{}) error {
		var errors []string
		for _, vd := range validators {
			err := vd(s)
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				// Just need one in chain to pass for this to be fine
				return nil
			}
		}

		// Happens if there are no validators
		if len(errors) == 0 {
			return nil
		}

		return fmt.Errorf(strings.Join(errors, " or "))
	}
}

// ValidateRegex Confirms that value matches the provided regex
func ValidateRegex(rx *regexp.Regexp, message string) Validator {
	return func(s interface{}) error {
		value, ok := s.(string)
		if !ok {
			return fmt.Errorf("Field must be string")
		}
		if !rx.MatchString(value) {
			return fmt.Errorf(message)
		}

		return nil
	}
}

// ValidateEmail Checks that email address is valid
func ValidateEmail() Validator {
	return ValidateRegex(emailRx, "Email address is invalid")
}

// ValidateFail Used for testing, always fails and adds the message
func ValidateFail(message string) Validator {
	return func(v interface{}) error {
		return fmt.Errorf(message)
	}
}

// ValidateNoError used when there's an error.  If there's an error, then it fails
func ValidateNoError(message string) Validator {
	return func(v interface{}) error {
		err, ok := v.(error)
		if !ok {
			return fmt.Errorf("Value must be an error")
		}
		if err != nil {

			return fmt.Errorf("message")
		}

		return nil
	}
}

// ValidateFalse Checks that a value is true
func ValidateFalse(msg string) Validator {
	return func(v interface{}) error {
		b, ok := v.(bool)
		if !ok {
			return fmt.Errorf("Value must be a bool")
		}

		if b {
			return fmt.Errorf(msg)
		}

		return nil
	}
}

// ValidateLength Requires the string to be a length between m and n inclusive
func ValidateLength(m int, n int) Validator {
	return func(v interface{}) error {
		field, ok := v.(string)
		if !ok {
			return ErrMustBeString
		}

		var msg string
		if m == n {
			msg = fmt.Sprintf("Must be exactly %d characters long", m)
		} else {
			msg = fmt.Sprintf("Must be between %d and %d characters long", m, n)
		}

		if len(field) < m || len(field) > n {
			return fmt.Errorf(msg)
		}

		return nil
	}
}

// ValidateMinimumLength Requires the minimum length for a string to be n characters
func ValidateMinimumLength(n int) Validator {
	return func(s interface{}) error {
		v, ok := s.(string)
		if !ok {
			return ErrMustBeString
		}
		if len(v) < n {
			return fmt.Errorf("Must be at least %d characters long", n)
		}
		return nil
	}
}

// ValidateUUID Verifies that provided string is a UUID
func ValidateUUID() Validator {
	return func(s interface{}) error {
		v, ok := s.(string)
		if !ok {
			return ErrMustBeString
		}
		_, err := uuid.FromString(v)

		if err != nil {
			return fmt.Errorf("Expected UUID for field")
		}

		return nil
	}
}

// ValidateNumbers Must be a string containing only numbers
func ValidateNumbers() Validator {
	return ValidateRegex(numberRx, "Must be numbers only")
}

// ValidatePositive Checks that an integer is positive
func ValidatePositive() Validator {
	return func(s interface{}) error {
		v, ok := s.(int)
		if !ok {
			return ErrMustBeInt
		}
		if v < 0 {
			return fmt.Errorf("Must be positive")
		}

		return nil
	}
}

// ValidateUsername Username validation
func ValidateUsername() Validator {
	return ValidateRegex(usernameRx, "Usernames can only contain letters, numbers, underscores and hyphens, and must not begin or end with an underscore or hyphen")
}

// ValidateLettersWithSpaces Must only contain letters and spaces
func ValidateLettersWithSpaces() Validator {
	return ValidateRegex(lettersWithSpacesRx, "Must only contain letters and spaces")
}

// ValidateLettersWithNumbers Must be a string containing only numbers
func ValidateLettersWithNumbers() Validator {
	return ValidateRegex(lettersWithNumbersRx, "Must only contain letters and numbers")
}

// ValidateLettersSpacesAndNumbers Must be a string containing only numbers
func ValidateLettersSpacesAndNumbers() Validator {
	return ValidateRegex(lettersSpacesAndNumbersRx, "Must only contain letters, spaces and numbers")
}

// ValidateURL Checks that URL is valid
func ValidateURL() Validator {
	return ValidateRegex(urlRx, "Website URL is invalid")
}

// ValidateTrue Checks that a value is true
func ValidateTrue(msg string) Validator {
	return func(s interface{}) error {
		v, ok := s.(bool)
		if !ok {
			return ErrMustBeBool
		}

		if !v {
			return fmt.Errorf(msg)
		}

		return nil
	}
}

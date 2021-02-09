package forms

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type TypeName string

var (
	TypeBool   TypeName = "bool"
	TypeString TypeName = "string"
	TypeUUID   TypeName = "uuid"
)

type Validator func(interface{}) error

type FieldDefinition struct {
	Name        string
	Group       string // Group is used when we want to break fields out into separate groups.  E.g., if creating or updating involves modifying multiple tables
	FieldType   TypeName
	Validations []Validator
}

type Definition struct {
	Fields map[string]FieldDefinition
}

// GetGroup Separates out all fields that belong to the specified group.  Should
// only be used on a processed map
func GetGroup(def Definition, group string, m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for name, field := range def.Fields {
		if field.Group == group {
			v, ok := m[name]
			if ok {
				out[name] = v
			}
		}
	}

	return out
}

type ApplyOptions struct {
	// If true, returns an error if the field is not in the permittedFields list,
	// otherwise it just removes that field silently
	ErrorOnPermission bool
}

func inList(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func ApplyDefinition(
	def Definition,
	permittedFields []string,
	m map[string]interface{},
	options ApplyOptions,
) (map[string]interface{}, map[string][]error, error) {
	if len(permittedFields) == 0 {
		return nil, nil, fmt.Errorf("Must provide a list of permitted fields with one or more fields")
	}
	applied := make(map[string]interface{})
	// Loop over this form's definitions, enforcing fields match the type that
	// we need.  This is because sometimes a html form will submit values, like
	// a checkbox, in a format other than a bool.  We convert that value into
	// the type we expect
	errors := make(map[string][]error)
	for k, v := range m {

		if !inList(permittedFields, k) {
			if options.ErrorOnPermission {
				errors[k] = append(errors[k], fmt.Errorf("Using field %s not permitted", k))
			}

			// Skip this field, as it's not permitted
			continue
		}

		var err error
		// Check if this is in our definitions
		fieldDef, ok := def.Fields[k]
		if !ok {
			continue
		}

		switch fieldDef.FieldType {
		case TypeBool:
			applied[k], err = ensureBool(v)
		case TypeString:
			applied[k], err = ensureString(v)
		case TypeUUID:
			applied[k], err = ensureUUID(v)
		default:
			err = fmt.Errorf("Unknown field type %s when converting field %s", fieldDef.FieldType, k)
		}

		if err != nil {
			// Continue, since if it's the wrong type, no point doing validations
			errors[k] = append(errors[k], err)
			continue
		}

		// Check validations:
		for _, validator := range fieldDef.Validations {
			err = validator(applied[k])
			if err != nil {
				errors[k] = append(errors[k], err)
			} else {
			}
		}
	}

	return applied, errors, nil
}

func ensureBool(s interface{}) (interface{}, error) {
	switch v := s.(type) {
	case bool:
		return s, nil
	case string:
		var b bool
		str := s.(string)
		switch str {
		case "true", "on":
			b = true
			return b, nil
		case "false", "off":
			b = false
			return b, nil
		default:
			return nil, fmt.Errorf("field must be true or false")
		}
	default:
		return nil, fmt.Errorf("Cannot convert type %T to bool", v)
	}
}

func ensureString(s interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", s)
	return str, nil
}

func ensureUUID(s interface{}) (interface{}, error) {
	switch v := s.(type) {
	case uuid.UUID:
		return s, nil
	case string:
		id, err := uuid.FromString(s.(string))
		return id, err
	default:
		return nil, fmt.Errorf("Cannot convert type %T to uuid.UUID", v)
	}
}

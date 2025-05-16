package forms

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testDefinition = Definition{
		Fields: map[string]FieldDefinition{
			"field1": {
				Group:       "test-group", // regular valid bool
				FieldType:   TypeBool,
				Validations: []Validator{},
			},
			"field2": {
				Group:       "test-group", // regular valid string
				FieldType:   TypeString,
				Validations: []Validator{},
			},
			"field3": {
				Group:       "test-group", // leading tabspace but trimmed
				FieldType:   TypeString,
				Validations: []Validator{},
			},
			"field4": {
				Group:         "test-group", // leading tabspace but not trimmed
				FieldType:     TypeString,
				Validations:   []Validator{},
				NotTrimSpaces: true,
			},
			"field5": {
				Group:       "test-group", // regular valid UUID
				FieldType:   TypeUUID,
				Validations: []Validator{},
			},
			"field6": {
				Group:       "test-group",
				FieldType:   TypeUUID,
				Validations: []Validator{}, // regular string but is not permitted field
			},
		},
	}

	testPermittedFields = []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
	}

	testInput = map[string]interface{}{
		"field1": "true",
		"field2": "this is a string",
		"field3": "	  this is also string but is tabbed ",
		"field4": "	this is also string but is tabbed and should not trim	 ",
		"field5": "123e4567-e89b-12d3-a456-426614174000",
		"field6": "this is a string but it won't get permitted",
	}

	testApplyOptionsNoError = ApplyOptions{
		ErrorOnPermission: MissingErrorNone,
	}

	testApplyOptionsFailOnError = ApplyOptions{
		ErrorOnPermission: MissingErrorFail,
	}
)

type dummyStruct struct{}

func Test_ApplyDefinition(t *testing.T) {
	t.Run("Success - No Error", func(t *testing.T) {
		transformedFields, errors, err := ApplyDefinition(testDefinition, testPermittedFields, testInput, testApplyOptionsNoError)
		assert.Equal(t, map[string]interface{}{
			"field1": true,
			"field2": "this is a string",
			"field3": "this is also string but is tabbed",
			"field4": "	this is also string but is tabbed and should not trim	 ",
			"field5": uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
		}, transformedFields)
		assert.Empty(t, errors)
		assert.NoError(t, err)
	})

	t.Run("Success - Error List on Non Permitted Field", func(t *testing.T) {
		transformedFields, errors, err := ApplyDefinition(testDefinition, testPermittedFields, testInput, testApplyOptionsFailOnError)
		assert.Equal(t, map[string]interface{}{
			"field1": true,
			"field2": "this is a string",
			"field3": "this is also string but is tabbed",
			"field4": "	this is also string but is tabbed and should not trim	 ",
			"field5": uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
		}, transformedFields)
		assert.NotEmpty(t, errors)
		assert.Equal(t, map[string][]error{
			"field6": {fmt.Errorf("Using field field6 not permitted")},
		}, errors)
		assert.NoError(t, err)
	})

	t.Run("Failure - No Permitted Fields", func(t *testing.T) {
		_, _, err := ApplyDefinition(testDefinition, []string{}, testInput, testApplyOptionsFailOnError)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Must provide a list of permitted fields with one or more fields")
	})
}

func Test_ensureBool(t *testing.T) {
	t.Run("Success - Is String", func(t *testing.T) {
		result, err := ensureBool("true")
		assert.Equal(t, true, result)
		assert.NoError(t, err)
	})

	t.Run("Success - Is String with Whitespace", func(t *testing.T) {
		result, err := ensureBool(" true  ")
		assert.Equal(t, true, result)
		assert.NoError(t, err)
	})

	t.Run("Success - Is Bool", func(t *testing.T) {
		result, err := ensureBool(true)
		assert.Equal(t, true, result)
		assert.NoError(t, err)
	})

	t.Run("Failure - Invalid String", func(t *testing.T) {
		_, err := ensureBool("dunno")
		assert.ErrorContains(t, err, "field must be true or false")
	})

	// This case should technically never be reachable
	t.Run("Failure - Invalid Interface", func(t *testing.T) {
		_, err := ensureBool(dummyStruct{})
		assert.ErrorContains(t, err, "Cannot convert type forms.dummyStruct to bool")
	})
}

func Test_ensureString(t *testing.T) {
	t.Run("Success - Is String", func(t *testing.T) {
		result, err := ensureString("lorem ipsum", false)
		assert.Equal(t, "lorem ipsum", result)
		assert.NoError(t, err)
	})

	t.Run("Success - Is String With Whitespace", func(t *testing.T) {
		result, err := ensureString("		lorem ipsum	 ", true)
		assert.Equal(t, "lorem ipsum", result)
		assert.NoError(t, err)
	})

	// This case should technically never be reachable
	t.Run("Failure - Invalid Interface", func(t *testing.T) {
		_, err := ensureString(dummyStruct{}, false)
		assert.ErrorContains(t, err, "Cannot convert type forms.dummyStruct to string")
	})
}

func Test_ensureUUID(t *testing.T) {
	dummyUUIDString := "123e4567-e89b-12d3-a456-426614174000"
	// uuid should contain 32 hex chars, total of 36 chars, in the form 8-4-4-4-12
	t.Run("Success - Is UUID", func(t *testing.T) {
		result, err := ensureUUID(uuid.Must(uuid.FromString(dummyUUIDString)))
		assert.Equal(t, dummyUUIDString, result.(uuid.UUID).String())
		assert.NoError(t, err)
	})

	t.Run("Success - Is String", func(t *testing.T) {
		result, err := ensureUUID(dummyUUIDString)
		assert.Equal(t, dummyUUIDString, result.(uuid.UUID).String())
		assert.NoError(t, err)
	})

	t.Run("Success - Is String With Whitespace", func(t *testing.T) {
		result, err := ensureUUID(fmt.Sprintf("	 	%s	 ", dummyUUIDString))
		assert.Equal(t, dummyUUIDString, result.(uuid.UUID).String())
		assert.NoError(t, err)
	})

	// This case should technically never be reachable
	t.Run("Failure - Invalid Interface", func(t *testing.T) {
		_, err := ensureUUID(dummyStruct{})
		assert.ErrorContains(t, err, "Cannot convert type forms.dummyStruct to uuid.UUID")
	})
}

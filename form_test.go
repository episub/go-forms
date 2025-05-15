package forms

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type dummyStruct struct{}

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
		result, err := ensureString("lorem ipsum")
		assert.Equal(t, "lorem ipsum", result)
		assert.NoError(t, err)
	})

	t.Run("Success - Is String With Whitespace", func(t *testing.T) {
		result, err := ensureString("		lorem ipsum	 ")
		assert.Equal(t, "lorem ipsum", result)
		assert.NoError(t, err)
	})

	// This case should technically never be reachable
	t.Run("Failure - Invalid Interface", func(t *testing.T) {
		_, err := ensureString(dummyStruct{})
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

package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategory_Validate(t *testing.T) {
	t.Run("Valid category", func(t *testing.T) {
		category := &Category{
			Name: "Electronics",
		}

		err := category.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid - empty name", func(t *testing.T) {
		category := &Category{
			Name: "",
		}

		err := category.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Category name is required")
	})
}

func TestCategory_BeforeCreate(t *testing.T) {
	t.Run("Should generate UUID if not set", func(t *testing.T) {
		category := &Category{
			Name: "Electronics",
		}

		err := category.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEqual(t, "", category.ID.String())
	})

	t.Run("Should not override existing UUID", func(t *testing.T) {
		id, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		
		category := &Category{
			ID:   id,
			Name: "Electronics",
		}

		err := category.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", category.ID.String())
	})
}

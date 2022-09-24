package mapper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unregisteredType struct{}

type testObject struct {
	GivenName  string
	FamilyName string
}

type testModel struct {
	Name string
}

func TestMapper(t *testing.T) {
	Register[*testObject, *testModel](func(item *testObject) *testModel {
		return &testModel{
			Name: fmt.Sprintf("%v %v", item.GivenName, item.FamilyName),
		}
	})

	t.Run("MapFrom", func(t *testing.T) {
		t.Run("when the mapping is registered", func(t *testing.T) {
			item := &testObject{
				GivenName:  "Tsuyoshi",
				FamilyName: "Garret",
			}

			model := MapFrom[*testObject, *testModel](item)

			require.NotNil(t, model)
			assert.Equal(t, "Tsuyoshi Garret", model.Name)
		})

		t.Run("When the mapping is not registered", func(t *testing.T) {
			item := &unregisteredType{}
			model := MapFrom[*unregisteredType, *testModel](item)

			assert.Nil(t, model)
		})
	})

	t.Run("MapEachFrom", func(t *testing.T) {
		t.Run("when the mapping is registered", func(t *testing.T) {
			datum := []*testObject{
				{GivenName: "Tsuyoshi", FamilyName: "Garret"},
				{GivenName: "Takashi", FamilyName: "Shirogane"},
			}

			results := MapEachFrom[*testObject, *testModel](datum)

			require.NotNil(t, results)
			require.Equal(t, 2, len(results))

			assert.Equal(t, "Tsuyoshi Garret", results[0].Name)
			assert.Equal(t, "Takashi Shirogane", results[1].Name)
		})

		t.Run("when the mapping is not registered", func(t *testing.T) {
			datum := []*unregisteredType{
				{},
			}

			results := MapEachFrom[*unregisteredType, *testModel](datum)

			require.NotNil(t, results)
			assert.Equal(t, 0, len(results))
		})
	})
}

package mapper

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unregisteredType struct{}

type greeter interface{ Greet() string }

type farewell interface{ Bye() string }

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

func TestMapperInterfaceTypes(t *testing.T) {
	Register[greeter, string](func(item greeter) string {
		return "hello"
	})

	t.Run("does not collide with another interface type", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.Equal(t, "", MapFrom[farewell, string](nil))
		})
	})
}

type raceInput struct{}

type raceOutput struct{}

func TestMapperConcurrentRegister(t *testing.T) {
	t.Run("does not race", func(t *testing.T) {
		var wg sync.WaitGroup

		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Register[*raceInput, *raceOutput](func(item *raceInput) *raceOutput {
					return &raceOutput{}
				})
			}()
		}

		wg.Wait()
	})
}

func TestMapEachFromZeroResults(t *testing.T) {
	Register[int, string](func(item int) string {
		if item == 0 {
			return ""
		}
		return strconv.Itoa(item)
	})

	t.Run("keeps a mapped result that is the zero value", func(t *testing.T) {
		results := MapEachFrom[int, string]([]int{0, 1, 2})

		assert.Equal(t, []string{"", "1", "2"}, results)
	})

	t.Run("drops every item when no mapping is registered", func(t *testing.T) {
		results := MapEachFrom[int, bool]([]int{0, 1, 2})

		assert.Empty(t, results)
	})
}

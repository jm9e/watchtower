package filters

import (
	"testing"

	"github.com/jm9e/watchtower/pkg/container/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWatchtowerContainersFilter(t *testing.T) {
	container := new(mocks.FilterableContainer)

	container.On("IsWatchtower").Return(true)

	assert.True(t, WatchtowerContainersFilter(container))

	container.AssertExpectations(t)
}

func TestNoFilter(t *testing.T) {
	container := new(mocks.FilterableContainer)

	assert.True(t, NoFilter(container))

	container.AssertExpectations(t)
}

func TestFilterByNames(t *testing.T) {
	var names []string

	filter := FilterByNames(names, nil)
	assert.Nil(t, filter)

	names = append(names, "test")

	filter = FilterByNames(names, NoFilter)
	assert.NotNil(t, filter)

	container := new(mocks.FilterableContainer)
	container.On("Name").Return("test")
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("NoTest")
	assert.False(t, filter(container))
	container.AssertExpectations(t)
}

func TestFilterByEnableLabel(t *testing.T) {
	filter := FilterByEnableLabel(NoFilter)
	assert.NotNil(t, filter)

	container := new(mocks.FilterableContainer)
	container.On("Enabled").Return(true, true)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, true)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, false)
	assert.False(t, filter(container))
	container.AssertExpectations(t)
}

func TestFilterByDisabledLabel(t *testing.T) {
	filter := FilterByDisabledLabel(NoFilter)
	assert.NotNil(t, filter)

	container := new(mocks.FilterableContainer)
	container.On("Enabled").Return(true, true)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, true)
	assert.False(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, false)
	assert.True(t, filter(container))
	container.AssertExpectations(t)
}

func TestBuildFilter(t *testing.T) {
	var names []string
	names = append(names, "test")

	filter := BuildFilter(names, false)

	container := new(mocks.FilterableContainer)
	container.On("Name").Return("Invalid")
	container.On("Enabled").Return(false, false)
	assert.False(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("test")
	container.On("Enabled").Return(false, false)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("Invalid")
	container.On("Enabled").Return(true, true)
	assert.False(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("test")
	container.On("Enabled").Return(true, true)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, true)
	assert.False(t, filter(container))
	container.AssertExpectations(t)
}

func TestBuildFilterEnableLabel(t *testing.T) {
	var names []string
	names = append(names, "test")

	filter := BuildFilter(names, true)

	container := new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, false)
	assert.False(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("Invalid")
	container.On("Enabled").Twice().Return(true, true)
	assert.False(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Name").Return("test")
	container.On("Enabled").Twice().Return(true, true)
	assert.True(t, filter(container))
	container.AssertExpectations(t)

	container = new(mocks.FilterableContainer)
	container.On("Enabled").Return(false, true)
	assert.False(t, filter(container))
	container.AssertExpectations(t)
}

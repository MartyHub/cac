package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	c, err := NewCache("testNewCache")
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, 0, c.Len())
	assert.True(t, c.exists())
	assert.Equal(t, "testNewCache.json", c.fileName())
}

func TestNewCache_noConfig(t *testing.T) {
	c, err := NewCache("")
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, 0, c.Len())
	assert.False(t, c.exists())
}

func TestCache_clean(t *testing.T) {
	c, err := NewCache("testClean")
	require.NoError(t, err)

	c.merge([]account{
		{
			Object:     "o1",
			Value:      "value1",
			StatusCode: 200,
			Timestamp:  now.Add(-1 * time.Hour),
		},
		{
			Object:     "o2",
			Value:      "value2",
			StatusCode: 200,
			Timestamp:  now.Add(-2 * time.Hour),
		},
		{
			Object:     "o3",
			Value:      "value3",
			StatusCode: 200,
			Timestamp:  now.Add(-3 * time.Hour),
		},
	})
	assert.Equal(t, 3, c.Len())

	c.clean(newFixedClock(), 2*time.Hour)

	assert.Equal(t, 1, c.Len())
	assert.Contains(t, c.Accounts, "o1")
}

func TestCache_SortedObjects(t *testing.T) {
	c := Cache{
		Accounts: map[string]account{
			"o1": {
				Object: "o1",
			},
			"o2": {
				Object: "o2",
			},
		},
	}

	assert.Equal(t, []string{"o1", "o2"}, c.SortedObjects(""))
	assert.Equal(t, []string{"o1", "o2"}, c.SortedObjects("o"))
	assert.Equal(t, []string{"o1"}, c.SortedObjects("o1"))
	assert.Equal(t, []string{}, c.SortedObjects("a"))
}

func TestCache_workflow(t *testing.T) {
	c, err := NewCache("testWorkflow")
	require.NoError(t, err)
	assert.Equal(t, 0, c.Len())

	c.merge([]account{
		{
			Object:     "o1",
			Value:      "value1",
			StatusCode: 200,
		},
	})
	require.NoError(t, c.save())

	c, err = NewCache("testWorkflow")
	require.NoError(t, err)
	assert.Equal(t, 1, c.Len())

	require.NoError(t, c.Remove())

	c, err = NewCache("testWorkflow")
	require.NoError(t, err)
	assert.Equal(t, 0, c.Len())
}

func TestCache_merge(t *testing.T) {
	c, err := NewCache("testMerge")
	require.NoError(t, err)

	c.Accounts["o1"] = account{
		Object: "o1",
		Value:  "value1",
	}
	c.Accounts["o2"] = account{
		Object: "o2",
		Value:  "value2",
	}

	c.merge([]account{
		{
			Object:     "o2",
			Value:      "newValue2",
			StatusCode: 200,
		},
		{
			Object:     "o3",
			Value:      "value3",
			StatusCode: 200,
		},
	})

	assert.Equal(t, 3, c.Len())
	assert.Contains(t, c.Accounts, "o1")
	assert.Contains(t, c.Accounts, "o2")
	assert.Contains(t, c.Accounts, "o3")
}

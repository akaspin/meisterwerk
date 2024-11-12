package model

import (
	"testing"

	"github.com/akaspin/meisterwerk/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote_Update(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()

	cId := uint64(storage.TestPackNextID())

	t.Run("ok", func(t *testing.T) {
		q1 := &Quote{
			CustomerID:  cId,
			Description: "quote1",
			Items: []*Item{
				{
					ItemID:  "product-1",
					Segment: "product",
					Price:   100,
					Tax:     0.1,
				},
				{
					ItemID:  "product-1",
					Segment: "product",
					Price:   100,
					Tax:     0.1,
				},
			},
		}
		id, err := q1.Create(conn.GORM)
		require.NoError(t, err)

		q2 := &Quote{
			ID:          id,
			CustomerID:  cId,
			Description: "quote1123",
			Status:      "accepted",
			Items: []*Item{
				{
					ItemID:  "product-1",
					Segment: "product",
					Price:   200,
					Tax:     0.1,
				},
			},
		}
		err = q2.Update(conn.GORM)
		require.NoError(t, err)

		q3 := &Quote{ID: id}
		err = q3.Find(conn.GORM)
		require.NoError(t, err)
		assert.Len(t, q3.Items, 1)
		assert.EqualValues(t, 200, q3.Items[0].Price)
		t.Log(q3.Status)
	})
	t.Run("not found", func(t *testing.T) {
		q := &Quote{
			ID:          cId,
			CustomerID:  cId,
			Description: "quote2",
		}
		err := q.Update(conn.GORM)
		assert.EqualError(t, err, "record not found")
	})
}

func TestQuote_Delete(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()
	cId := uint64(storage.TestPackNextID())

	q1 := &Quote{
		CustomerID:  cId,
		Description: "quote1",
		Items: []*Item{
			{
				ItemID:  "product-1",
				Segment: "product",
				Price:   100,
				Tax:     0.1,
			},
			{
				ItemID:  "product-1",
				Segment: "product",
				Price:   100,
				Tax:     0.1,
			},
		},
	}
	id, err := q1.Create(conn.GORM)
	require.NoError(t, err)

	q2 := &Quote{ID: id}
	err = q2.Delete(conn.GORM)
	require.NoError(t, err)

	err = q1.Find(conn.GORM) // cascade guarantees features removal
	assert.EqualError(t, err, "record not found")

	err = q2.Delete(conn.GORM)
	assert.EqualError(t, err, "record not found")
}

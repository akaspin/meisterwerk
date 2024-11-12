package model

import (
	"fmt"
	"testing"

	"github.com/akaspin/meisterwerk/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedQuotes(t *testing.T, db *gorm.DB, cID int32, name string, n int) []uint64 {
	t.Helper()
	var qs Quotes
	for i := 0; i < n; i++ {
		q := &Quote{
			CustomerID:  uint64(cID),
			Description: fmt.Sprintf("quote-%s-%d", name, i),
			Status:      "pending",
		}
		for j := 0; j < i; j++ {
			q.Items = append(q.Items, &Item{
				ItemID:  fmt.Sprintf("item-%s-%d-%d", name, i, j),
				Segment: "product",
				Price:   100,
				Tax:     0.1,
			})
		}
		qs = append(qs, q)
	}
	ids, err := qs.Create(db)
	require.NoError(t, err)
	return ids
}

func TestQuotes_Create(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()
	cID := uint64(storage.TestPackNextID())

	t.Run("create", func(t *testing.T) {
		quotes := Quotes{
			{
				CustomerID:  cID,
				Description: "quote1",
				Items: []*Item{
					{
						ItemID:  "product-1",
						Segment: "feature-1",
						Price:   100,
						Tax:     0.1,
					},
				},
			},
			{
				CustomerID:  cID,
				Description: "quote2",
			},
		}
		ids, err := quotes.Create(conn.GORM)
		require.NoError(t, err)
		assert.Len(t, ids, 2)
	})
	t.Run("duplicate quotes", func(t *testing.T) {
		quotes := Quotes{
			{
				CustomerID:  cID,
				Description: "dup-quote1",
				Items: []*Item{
					{
						ItemID:  "dup-quote-product-1",
						Segment: "product",
						Price:   100,
						Tax:     0.1,
					},
				},
			},
			{
				CustomerID:  cID,
				Description: "dup-quote1", // boom
			},
		}
		_, err := quotes.Create(conn.GORM)
		require.Error(t, err)
	})
}

func TestQuotes_Update(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()

	cId := uint64(storage.TestPackNextID())

	t.Run("update", func(t *testing.T) {
		quotes1 := Quotes{
			{
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
			},
			{
				CustomerID:  cId,
				Description: "quote2",
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
			},
		}
		ids, err := quotes1.Create(conn.GORM)
		require.NoError(t, err)
		assert.Len(t, ids, 2)

		quotes2 := Quotes{
			{
				ID:          quotes1[0].ID,
				CustomerID:  cId,
				Description: "quote1",
				Items: []*Item{
					{
						ItemID:  "product-1",
						Segment: "product",
						Price:   200,
						Tax:     0.1,
					},
				},
			},
			{
				ID:          quotes1[1].ID,
				CustomerID:  cId,
				Description: "quote2",
				Items: []*Item{
					{
						ItemID:  "product-1",
						Segment: "product",
						Price:   200,
						Tax:     0.1,
					},
				},
			},
		}
		err = quotes2.Update(conn.GORM)
		require.NoError(t, err)

		var quotes3 Quotes
		err = conn.GORM.Preload(clause.Associations).Find(&quotes3, "quotes.id IN ?", ids).Error
		require.NoError(t, err)
		assert.Len(t, quotes3, 2)
		assert.Len(t, quotes3[0].Items, 1)
		assert.EqualValues(t, 200, quotes3[0].Items[0].Price)
	})
}

func TestQuotes_Find(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()

	cId := uint64(storage.TestPackNextID())

	t.Run("find", func(t *testing.T) {
		quotes := Quotes{
			{
				CustomerID:  cId,
				Description: "quote1",
				Items: []*Item{
					{
						ItemID:  "product-1",
						Segment: "product",
						Price:   200,
						Tax:     0.1,
					},
				},
			},
			{
				CustomerID:  cId,
				Description: "quote2",
				Items: []*Item{
					{
						ItemID:  "product-1",
						Segment: "product",
						Price:   200,
						Tax:     0.1,
					},
				},
			},
		}
		ids, err := quotes.Create(conn.GORM)
		require.NoError(t, err)

		var quotes2 Quotes
		err = (&quotes2).Find(conn.GORM, ids)
		require.NoError(t, err)
		assert.Len(t, quotes2, 2)
		assert.Len(t, quotes2[0].Items, 1)
		assert.Contains(t, ids, quotes2[0].ID)
	})
}

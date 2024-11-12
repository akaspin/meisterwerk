package model

import (
	"slices"
	"testing"

	"github.com/akaspin/meisterwerk/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCriteria_Find(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()

	cID := storage.TestPackNextID()
	ids := seedQuotes(t, conn.GORM, cID, "criteria", 10)

	t.Run("by ids", func(t *testing.T) {
		res, err := (&Criteria{IDs: ids}).Find(conn.GORM)
		require.NoError(t, err)
		var newIds []uint64
		for _, q := range res {
			newIds = append(newIds, q.ID)
		}

		res2 := slices.DeleteFunc(res, func(i *Quote) bool {
			return slices.Contains(ids, i.ID)
		})
		assert.Empty(t, res2)
		assert.True(t, slices.IsSorted(newIds))
	})
	t.Run("desc", func(t *testing.T) {
		res, err := (&Criteria{IDs: ids, Desc: true}).Find(conn.GORM)
		require.NoError(t, err)
		assert.Len(t, res, 10)
		var newIds []uint64
		for _, q := range res {
			newIds = append(newIds, q.ID)
		}
		slices.Reverse(newIds)
		assert.True(t, slices.IsSorted(newIds))
	})
}

func TestCriteria_Delete(t *testing.T) {
	conn, connErr := storage.Connect(storage.TestPackDBConfig("meisterwerk"))
	require.NoError(t, connErr)
	defer conn.Close()

	cId := storage.TestPackNextID()
	ids := seedQuotes(t, conn.GORM, cId, "quotes-delete", 3)
	t.Run("ok", func(t *testing.T) {
		cri := &Criteria{IDs: ids}
		err := cri.Delete(conn.GORM)
		require.NoError(t, err)

		res, err := cri.Find(conn.GORM)
		require.NoError(t, err)
		assert.Len(t, res, 0)

		t.Run("not found", func(t *testing.T) {
			err1 := cri.Delete(conn.GORM)
			assert.EqualError(t, err1, "record not found")
		})
	})
}

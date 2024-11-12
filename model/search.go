package model

import (
	"errors"

	"github.com/akaspin/meisterwerk/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidCriteria = errors.New("invalid criteria")

// Criteria defines search criteria for quotes
type Criteria struct {
	IDs         []uint64
	CustomerIDs []uint64
	Limit       int
	Skip        int
	Desc        bool
}

func (c *Criteria) Find(db *gorm.DB) (Quotes, error) {
	var cls []clause.Expression
	if len(c.IDs) > 0 {
		cls = append(cls, clause.Where{
			Exprs: []clause.Expression{
				clause.IN{Column: "id", Values: toAnySlice(c.IDs)},
			}})
	}
	if len(c.CustomerIDs) > 0 {
		cls = append(cls, clause.Where{
			Exprs: []clause.Expression{
				clause.IN{Column: "qustomer_id", Values: toAnySlice(c.IDs)},
			}})
	}
	if c.Limit > 0 || c.Skip > 0 {
		cl := clause.Limit{
			Offset: c.Skip,
		}
		if c.Limit > 0 {
			cl.Limit = &c.Limit
		}
		cls = append(cls, cl)
	}

	var res Quotes
	q1 := db.Preload(clause.Associations).Clauses(cls...).Order(
		clause.OrderByColumn{
			Column: clause.Column{Name: "id"},
			Desc:   c.Desc,
		}).Find(&res)
	return res, q1.Error
}

// Delete deletes quotes only by their IDs
func (c *Criteria) Delete(db *gorm.DB) error {
	if len(c.IDs) == 0 {
		return ErrInvalidCriteria
	}
	cls := clause.Where{
		Exprs: []clause.Expression{
			clause.IN{Column: "id", Values: toAnySlice(c.IDs)},
		}}
	err := storage.WithTx(db, func(tx *gorm.DB) error {
		q := tx.Clauses(cls).Delete(&Quotes{})
		if q.RowsAffected != int64(len(c.IDs)) {
			return gorm.ErrRecordNotFound
		}
		return q.Error
	})
	return err
}

func toAnySlice[S ~[]E, E any](s S) []any {
	var res []any
	for _, v := range s {
		res = append(res, v)
	}
	return res
}

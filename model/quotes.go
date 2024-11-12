package model

import (
	"github.com/akaspin/meisterwerk/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Quotes can be used for bulk operations
//
// Note: GORM is already repository abstraction. One more dedicated intermediate type
// like "QuotesRepository" in separate package is redundant in scope of this project.
//
// Of course. In projects with more complex business logic the dedicated repository will make sense.
type Quotes []*Quote

func (q *Quotes) Create(db *gorm.DB) ([]uint64, error) {
	err := db.Create(q).Error
	if err != nil {
		return nil, err
	}
	var res []uint64
	for _, quote := range *q {
		res = append(res, quote.ID)
	}
	return res, nil
}

func (q *Quotes) Update(db *gorm.DB) error {
	return storage.WithTx(db, func(tx *gorm.DB) error {
		for _, quote := range *q {
			err := quote.updateUnsafe(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (q *Quotes) Find(db *gorm.DB, ids []uint64) error {
	return db.Preload(clause.Associations).Where("id in ?", ids).Find(q).Error
}

package model

import (
	"time"

	"github.com/akaspin/meisterwerk/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Quote represents quote in storage
//
// Note: GORM is already repository abstraction. One more dedicated intermediate type
// like "QuotesRepository" in separate package is redundant in scope of this project.
//
// Of course. In projects with more complex business logic the dedicated repository will make sense.
type Quote struct {
	ID          uint64    `gorm:"primaryKey"`
	CustomerID  uint64    `gorm:"uniqueIndex:unique_customer_id_description;not null"`
	Description string    `gorm:"uniqueIndex:unique_customer_id_description;not null"`
	Items       []*Item   `gorm:"foreignKey:QuoteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status      string    `gorm:"not null;default=pending"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null"`
}

func (q *Quote) Create(db *gorm.DB) (uint64, error) {
	err := db.Create(q).Error
	if err != nil {
		return 0, err
	}
	return q.ID, nil
}

func (q *Quote) Update(db *gorm.DB) error {
	return storage.WithTx(db, q.updateUnsafe)
}

func (q *Quote) updateUnsafe(db *gorm.DB) error { // optimisations removed
	m := db.Omit(clause.Associations).Updates(q)
	if m.Error != nil {
		return m.Error
	}
	if m.RowsAffected != 1 { // not found
		return gorm.ErrRecordNotFound
	}
	err := db.Unscoped().Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Model(q).Association("Items").Unscoped().Replace(q.Items)
	return err
}

func (q *Quote) Find(db *gorm.DB) error {
	return db.Preload(clause.Associations).First(q).Error
}

func (q *Quote) Delete(db *gorm.DB) error {
	x := db.Delete(q)
	if x.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return x.Error
}

// The Item in Quote
//
// Note that use float for money and percentages is not best approach.
// But in scope of this task this approach reduces complexity.
//
// The Segment can be implemented with ENUM in Postgres. But this requires extra
// manual migration and omited for now.
type Item struct {
	QuoteID   uint64    `gorm:"primaryKey"`
	ItemID    string    `gorm:"primaryKey"`
	Segment   string    `gorm:"primaryKey"`
	Price     float32   `gorm:"not null"`
	Tax       float32   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}

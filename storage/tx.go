package storage

import "gorm.io/gorm"

func WithTx(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

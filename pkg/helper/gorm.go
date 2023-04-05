package helper

import "gorm.io/gorm"

func TransactionHandle(txn *gorm.DB, err *error) {
	if *err != nil {
		txn.Rollback()
	} else {
		txn.Commit()
	}
}

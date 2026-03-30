package repository

import "gorm.io/gorm"

func applyOrder(db *gorm.DB, sortBy, order, defaultOrder string) *gorm.DB {
	ord := order
	if ord != "asc" && ord != "desc" {
		ord = "asc"
	}

	o := defaultOrder
	if sortBy != "" {
		o = sortBy + " " + ord
	}
	return db.Order(o)
}


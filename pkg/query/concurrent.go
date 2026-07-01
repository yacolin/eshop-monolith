package query

import (
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// ConcurrentCountList 并发执行 COUNT 和 LIMIT 查询。
// db 的 WHERE/Joins 等条件需要在调用前预先设置好。
// 泛型参数 T 为列表元素类型，应与 db.Model() 一致。
func ConcurrentCountList[T any](db *gorm.DB, page, size int) ([]T, int64, error) {
	g, ctx := errgroup.WithContext(db.Statement.Context)

	var total int64
	g.Go(func() error {
		return db.Session(&gorm.Session{Context: ctx}).Count(&total).Error
	})

	var list []T
	g.Go(func() error {
		offset := (page - 1) * size
		return db.Session(&gorm.Session{Context: ctx}).Offset(offset).Limit(size).Find(&list).Error
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

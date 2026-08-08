//go:build ignore

// gen_schema.go 一次性工具:在临时数据库上跑 AutoMigrate 生成 CREATE TABLE DDL
// 用法: go run scripts/gen_schema.go > docs/schema.sql
// 说明: 在临时库 eshop_schema_gen 上执行, 完成后自动 DROP, 不影响业务库
// 注意: 不用 -tags ignore(显式文件参数时 go run 本就忽略 build tag;加 tag 会破坏标准库构建)
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"eshop-monolith/internal/base"
	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/marketing"
	"eshop-monolith/internal/product"
	"eshop-monolith/internal/review"
	"eshop-monolith/internal/staff"
	"eshop-monolith/internal/trade"
	"eshop-monolith/internal/user"
)

const (
	dsn         = "root:123456@tcp(localhost:3306)/eshop_db?charset=utf8mb4&parseTime=True&loc=Local"
	genDatabase = "eshop_schema_gen"
)

func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 创建临时库
	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", genDatabase)).Error; err != nil {
		log.Fatalf("清理临时库失败: %v", err)
	}
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", genDatabase)).Error; err != nil {
		log.Fatalf("创建临时库失败: %v", err)
	}
	genDB, err := gorm.Open(mysql.Open(fmt.Sprintf("root:123456@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", genDatabase)),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("连接临时库失败: %v", err)
	}

	// 全部模型
	models := []any{
		// base
		&base.Notification{}, &base.NotificationRead{}, &base.NotificationTemplate{},
		// inventory
		&inventory.Inventory{}, &inventory.InventoryLog{},
		// marketing
		&marketing.Promotion{}, &marketing.PromotionRule{}, &marketing.PromotionProduct{},
		&marketing.UsageLog{}, &marketing.UserPromotion{},
		// product
		&product.Attribute{}, &product.Brand{}, &product.CategoryBrand{}, &product.Category{},
		&product.Description{}, &product.ProductAttribute{}, &product.SKU{}, &product.SPU{},
		// review
		&review.Review{}, &review.ReviewMedia{}, &review.ReviewReply{}, &review.ReviewRating{}, &review.ReviewAuditLog{},
		// trade
		&trade.Cart{}, &trade.CartItem{}, &trade.OrderItem{}, &trade.OrderLog{}, &trade.Order{},
		&trade.PaymentLog{}, &trade.Payment{}, &trade.Refund{},
		// user
		&user.Address{}, &user.UserInfo{}, &user.LoginHistory{}, &user.User{},
		// staff
		&staff.Staff{}, &staff.SysRole{}, &staff.SysPermission{}, &staff.SysStaffRole{},
		&staff.SysRolePermission{}, &staff.StaffLoginHistory{},
	}

	if err := genDB.AutoMigrate(models...); err != nil {
		log.Fatalf("AutoMigrate 失败: %v", err)
	}

	// 导出每个表的 SHOW CREATE TABLE
	var tables []string
	if err := genDB.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		log.Fatalf("列出表失败: %v", err)
	}
	for _, t := range tables {
		var rows []map[string]interface{}
		if err := genDB.Raw("SHOW CREATE TABLE `" + t + "`").Scan(&rows).Error; err != nil {
			log.Fatalf("导出 %s 失败: %v", t, err)
		}
		fmt.Printf("-- ------------------------------------------------------------\n-- Table: %s\n-- ------------------------------------------------------------\n", t)
		fmt.Println(rows[0]["Create Table"].(string) + ";")
		fmt.Println()
	}

	// 清理临时库
	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", genDatabase)).Error; err != nil {
		log.Printf("警告: 清理临时库失败: %v", err)
	}
}

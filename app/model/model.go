// package model
//
// import (
// 	"time"
//
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// )
//
// //database table model
// type LoadedFunds struct {
//   Id      uint  `gorm:"primaryKey"`
//   Customer_id   uint
//   Load_amount    float64
//   Time time.Time
// }
//
//
// // DBMigrate will create and migrate the tables, and then make the some relationships if necessary
// func DBMigrate(db *gorm.DB) *gorm.DB {
// 	db.AutoMigrate(&Project{}, &Task{})
// 	db.Model(&Task{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
// 	return db
// }

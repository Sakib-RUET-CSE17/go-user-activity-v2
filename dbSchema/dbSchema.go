package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID             uint
	FirstName      string
	LastName       string
	Country        string
	ProfilePicture string
}

type Activity struct {
	ID     uint
	Points uint
}

type ActivityLog struct {
	ID         uint
	UserId     uint
	ActivityID uint
	LoggedAt   time.Time
}

func main() {
	// Connect to the MySQL database
	connectionString := "sakib:changeMe@tcp(localhost:49153)/userActivity?charset=utf8mb4&parseTime=True&loc=Local"
	db, dbErr := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// logger.Info এর কারণে কুয়েরি গুলো স্ট্যান্ডার্ড
		// আউটপুটে প্রিন্ট হতে থাকবে
	})
	if dbErr != nil {
		panic(dbErr)
	}
	// Auto Migrate the table
	err := db.AutoMigrate(&User{}, &Activity{}, &ActivityLog{})
	if err != nil {
		panic("failed to auto migrate table")
	}
}

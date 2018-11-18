package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type user struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type order struct {
	ID          int    `json:"id" gorm:"primary_key"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	User        *user  `json:"user,omitempty" gorm:"ForeignKey:UserID"`
}

func newDB() *gorm.DB {
	cs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", "root", "secret", "db", 3306, "example")
	connection, err := gorm.Open("mysql", cs)
	if err != nil {
		msg := "failed to connect to database: " + err.Error()
		panic(msg)
	}
	connection.LogMode(true)
	return connection
}

func createOrder(userName string, db *gorm.DB) {
	u := &user{
		Name: userName,
	}

	errs := db.Create(u).GetErrors()
	if len(errs) > 0 {
		msg := fmt.Sprintf("failed to save test user %s", userName)
		panic(msg)
	}

	o1 := &order{
		UserID:      u.ID,
		Description: fmt.Sprintf("Test Order for %s", userName),
	}

	errs = db.Create(o1).GetErrors()
	if len(errs) > 0 {
		msg := fmt.Sprintf("failed to save test order %s", userName)
		panic(msg)
	}
}

func main() {
	db := newDB()
	db.DropTable(&order{}).DropTable(&user{}).AutoMigrate(&user{}, &order{})

	createOrder("Test User 1", db)
	createOrder("Test User 2", db)

	// Get a user by their name - this works fine.
	var users []user
	errs := db.Where("name = ?", "Test User 2").Find(&users).GetErrors()
	if len(errs) > 0 {
		panic("failed to find users by name")
	} else {
		msg := fmt.Sprintf("found %v users", len(users))
		println(msg)
	}

	// Get all orders - this works fine.
	var allOrders []order
	errs = db.Find(&allOrders).GetErrors()
	if len(errs) > 0 {
		panic("failed to get all orders")
	} else {
		msg := fmt.Sprintf("found %v orders", len(allOrders))
		println(msg)
	}

	// Get orders by description - this also works fine.
	var whereOrders []order
	errs = db.Where("description = ?", "Test Order for Test User 2").Find(&allOrders).GetErrors()
	if len(errs) > 0 {
		panic("failed to get orders by id")
	} else {
		msg := fmt.Sprintf("found %v orders", len(whereOrders))
		println(msg)
	}

	// Trying to get an order by the user's name causes an error because the
	// Preloads are processed after the query is executed.
	var orders []order
	errs = db.Preload("User").Where("users.name = ?", "Test User 2").Find(&orders).GetErrors()
	if len(errs) > 0 {
		msg := fmt.Errorf("failed to get order by user name: %v", errs)
		println(msg)
	} else {
		msg := fmt.Sprintf("found %v orders", len(orders))
		println(msg)
	}

	// Trying to use Join works!
	var joinsOrders []order
	errs = db.Joins("LEFT JOIN users ON orders.user_id = users.id").Where("users.name = ?", "Test User 2").Find(&joinsOrders).GetErrors()
	if len(errs) > 0 {
		msg := fmt.Errorf("failed to get order by user name: %v", errs)
		println(msg)
	} else {
		msg := fmt.Sprintf("found %v orders", len(joinsOrders))
		println(msg)
	}

	println("done")
}

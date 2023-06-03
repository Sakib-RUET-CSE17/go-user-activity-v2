package main

import (
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// type Activity struct {
// 	ID             uint
// 	UserActivities []UserActivity
// 	Name           string
// 	Point          uint
// }

// type UserActivity struct {
// 	ID         uint
// 	UserId     uint
// 	ActivityID uint
// 	CreatedAt  time.Time
// }

var db *gorm.DB

func main() {
	connectionString := "sakib:changeMe@tcp(localhost:49154)/userActivity?charset=utf8mb4&parseTime=True&loc=Local"
	var dbErr error
	db, dbErr = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if dbErr != nil {
		panic(dbErr)
	}

	e := echo.New()
	fmt.Println(db)
	e.GET("/hello", hello)
	e.GET("/userActivities", getUserActivities)
	e.POST("/users", createUser)
	e.PATCH("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	err := e.Start(":1324")
	if err != nil {
		panic(err)
	}
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type UserActivity struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"name"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profilePicture"`
	TotalPoint     uint   `json:"totalPoint"`
	Ranking        uint   `json:"rank"`
	// UserActivities []UserActivity `json:"userActivities"`
}

func getUserActivities(c echo.Context) error {
	rankingTypeStr := c.QueryParam("rankingType")
	weekStr := c.QueryParam("week")
	week, _ := strconv.Atoi(weekStr)
	monthStr := c.QueryParam("month")
	month, _ := strconv.Atoi(monthStr)
	yearStr := c.QueryParam("year")
	year, _ := strconv.Atoi(yearStr)

	var userActivities []UserActivity

	if rankingTypeStr == "weekly" {
		query := `select
		users.id,
		first_name,
		country,
		profile_picture,
		week(logged_at) as week,
		year(logged_at) as year,
		sum(points) as total_point,
		dense_rank() over (
		  order by
			sum(points) desc
		) as ranking
	  from
		users
		join activity_logs on users.id = activity_logs.user_id
		join activities on activity_logs.activity_id = activities.id
	  group by
		users.id,
		week,
		year
	  having
		week = ? && year = ?`
		resp := db.Raw(query, week, year).Find(&userActivities)
		if resp.Error != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}
	} else if rankingTypeStr == "monthly" {
		query := `select
		users.id,
		first_name,
		country,
		profile_picture,
		month(logged_at) as month,
		year(logged_at) as year,
		sum(points) as total_point,
		dense_rank() over (
		  order by
			sum(points) desc
		) as ranking
	  from
		users
		join activity_logs on users.id = activity_logs.user_id
		join activities on activity_logs.activity_id = activities.id
	  group by
		users.id,
		month,
		year
	  having
		month = ? && year = ?`
		resp := db.Raw(query, month, year).Find(&userActivities)
		if resp.Error != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}
	} else {
		query := `select
		users.id,
		first_name,
		country,
		profile_picture,
		sum(points) as total_point,
		dense_rank() over (
		  order by
			sum(points) desc
		) as ranking
	  from
		users
		join activity_logs on users.id = activity_logs.user_id
		join activities on activity_logs.activity_id = activities.id
	  group by
		users.id`
		resp := db.Raw(query).Find(&userActivities)
		if resp.Error != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}
	}

	return c.JSON(http.StatusOK, userActivities)
}

type User struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profilePicture"`
}

func createUser(c echo.Context) error {
	user := &User{}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func deleteUser(c echo.Context) error {
	var user User
	id := c.Param("id")

	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := db.Where("id = ?", ID).Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "User deleted successfully")
}

func updateUser(c echo.Context) error {
	user := &User{}
	id := c.Param("id")

	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := db.Where("id = ?", ID).Updates(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, "User updated successfully")
}

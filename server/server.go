package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Task struct {
	Id string `json:"id"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Todoaroo!")
	})
	e.GET("/tasks", getTasks)
	e.GET("/tasks/:id", getTask)
	e.GET("/tasks/:id/", getTask)

	e.Logger.Fatal(e.Start(":1323"))
}

func getTasks(c echo.Context) error {
	t := &Task{
		Id: "1",
	}
	return c.JSON(http.StatusOK, t)
}

func getTask(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

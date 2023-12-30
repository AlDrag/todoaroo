package main

import (
	"context"
	"net/http"
	database "todoaroo/database/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=admin password=password123 dbname=todoaroo sslmode=disable host=localhost port=5432")
	if err != nil {
		// TODO Handle
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	e.GET("/tasks", getTasks(queries))
	e.GET("/tasks/:id", getTask(queries))
	e.POST("/tasks", createTask(queries))

	e.Logger.Fatal(e.Start(":1323"))
}

func getTask(queries *database.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var uuid pgtype.UUID
		err := uuid.Scan(c.Param("id"))
		if err != nil {
			// TODO Handle
		}

		task, err := queries.GetTask(c.Request().Context(), uuid)
		if err != nil {
			// TODO Handle
		}
		return c.JSON(http.StatusOK, task)
	}
}

func getTasks(queries *database.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		tasks, err := queries.ListTasks(c.Request().Context())
		if err != nil {
			// TODO Handle
		}
		return c.JSON(http.StatusOK, tasks)
	}
}

type CreateTaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createTask(queries *database.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		taskParams := new(CreateTaskBody)
		if err := c.Bind(taskParams); err != nil {
			return err
		}
		task, err := queries.CreateTask(c.Request().Context(), database.CreateTaskParams{
			Title:       taskParams.Title,
			Description: pgtype.Text{String: taskParams.Description, Valid: true},
		})
		if err != nil {
			// TODO Handle
		}
		return c.JSON(http.StatusCreated, task)
	}
}

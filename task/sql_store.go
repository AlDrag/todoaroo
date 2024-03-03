package task

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TaskSqlStore struct {
	db *sql.DB
}

func ConnectNew() (*TaskSqlStore, error) {
	// Open a database file, creating it if it doesn't exist
	db, err := sql.Open("sqlite3", "database.db")

	if err != nil {
		return nil, err
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &TaskSqlStore{
		db: db,
	}, nil
}

func (store *TaskSqlStore) Init() error {
	// Create Tasks table
	_, err := store.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (store *TaskSqlStore) Close() error {
	return store.db.Close()
}

func (store *TaskSqlStore) Create(title string, description string) (*Task, error) {
	query := "INSERT INTO tasks (title, description) VALUES (?, ?)"

	result, err := store.db.Exec(query, title, description)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Task{
		id:          int(lastID),
		title:       title,
		description: description,
	}, nil
}

func (store *TaskSqlStore) List() (list []Task, err error) {
	rows, err := store.db.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var id int
		var title string
		var description string
		err := rows.Scan(&id, &title, &description)
		if err != nil {
			log.Fatal(err)
		}
		var task = Task{
			id:          id,
			title:       title,
			description: description,
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (store *TaskSqlStore) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"

	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

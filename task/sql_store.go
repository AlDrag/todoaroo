package task

import (
	"database/sql"
	"log"
)

type TaskSqlStore struct {
	*sql.DB
}

func NewStore(db *sql.DB) *TaskSqlStore {
	return &TaskSqlStore{
		DB: db,
	}
}

func (db *TaskSqlStore) Create(title string, description string) (*Task, error) {
	query := "INSERT INTO tasks (title, description) VALUES (?, ?)"

	result, err := db.Exec(query, title, description)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Task{
		id:          int(lastID),
		title:       title,
		description: description,
	}, nil
}

func (db *TaskSqlStore) List() (list []Task, err error) {
	rows, err := db.Query("SELECT id, title, description FROM tasks")
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
		log.Fatal(err)
	}

	return tasks, nil
}

func (db *TaskSqlStore) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"

	_, err := db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

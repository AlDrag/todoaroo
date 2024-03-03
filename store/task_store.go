package store

import "todoaroo/task"

type TaskStore interface {
	Init() error
	Close() error
	Create(title string, description string) (*task.Task, error)
	List() ([]task.Task, error)
	Delete(id int) error
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"todoaroo/task"
	"todoaroo/task_input"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type taskStore interface {
	Create(title string, description string) (task.Task, error)
	List() ([]task.Task, error)
	Delete(id int) error
}

type model struct {
	list      list.Model
	taskStore *task.TaskSqlStore

	TaskInputModel task_input.TaskInput
	ShowTextInput  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "a":
			if !m.ShowTextInput {
				m.TaskInputModel = task_input.InitialModel()
				m.ShowTextInput = true
			}
		case "x":
			if !m.ShowTextInput {
				selectedItem := m.list.Items()[m.list.Index()]
				selectedTask := selectedItem.(task.Task)
				m.taskStore.Delete(selectedTask.ID())
				m.list.RemoveItem(m.list.Index())
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.ShowTextInput {
		taskInputModel, cmd := m.TaskInputModel.Update(msg)
		m.TaskInputModel = taskInputModel.(task_input.TaskInput)
		if m.TaskInputModel.Submitted {
			title, description := m.TaskInputModel.GetNewTask()
			newTask, _ := m.taskStore.Create(title, description)
			m.list.InsertItem(len(m.list.Items()), newTask)
			m.TaskInputModel.Submitted = false
			m.ShowTextInput = false
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.ShowTextInput {
		return m.TaskInputModel.View()
	}
	return docStyle.Render(m.list.View())
}

func main() {
	// Open a database file, creating it if it doesn't exist
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to SQLite database")

	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created tasks table")

	taskStore := task.NewStore(db)

	tasks := taskStore.List()

	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}

	tasksList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	tasksList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add item"),
			),
			key.NewBinding(
				key.WithKeys("x"),
				key.WithHelp("a", "delete item"),
			),
		}
	}
	m := model{list: tasksList, taskStore: taskStore}
	m.list.Title = "Todoaroo list"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

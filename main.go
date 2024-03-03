package main

import (
	"fmt"
	"log"
	"os"

	"todoaroo/task"
	"todoaroo/task_input"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type taskStore interface {
	Init() error
	Close() error
	Create(title string, description string) (*task.Task, error)
	List() ([]task.Task, error)
	Delete(id int) error
}

type model struct {
	list      list.Model
	taskStore taskStore

	taskInputModel task_input.TaskInput
	showTextInput  bool
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
			if !m.showTextInput {
				m.taskInputModel = task_input.InitialModel()
				m.showTextInput = true
				return m, nil
			}
		case "x":
			if !m.showTextInput {
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

	if m.showTextInput {
		taskInputModel, cmd := m.taskInputModel.Update(msg)
		m.taskInputModel = taskInputModel.(task_input.TaskInput)
		if m.taskInputModel.Submitted {
			title, description := m.taskInputModel.GetNewTask()
			newTask, _ := m.taskStore.Create(title, description)
			m.list.InsertItem(len(m.list.Items()), *newTask)
			m.taskInputModel.Submitted = false
			m.showTextInput = false
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showTextInput {
		return m.taskInputModel.View()
	}
	return docStyle.Render(m.list.View())
}

func main() {
	// Connect to store
	taskStore, err := task.ConnectNew()
	if err != nil {
		log.Fatal(err)
	}

	// Initialise store structure
	err = taskStore.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer taskStore.Close()

	// Request initial task list
	tasks, err := taskStore.List()
	if err != nil {
		log.Fatal(err)
	}

	// Transform tasks slice to list.Item interface slice
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}

	// Construct tasks list with custom key bindings
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

	// Initialise Bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run application
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

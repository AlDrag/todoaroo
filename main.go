package main

import (
	"fmt"
	"log"
	"os"

	"todoaroo/store"
	"todoaroo/task"
	"todoaroo/task_input"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tasksListMode int = iota
	createMode
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	mode      int
	taskStore store.TaskStore

	list           list.Model
	taskInputModel task_input.TaskInput
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case tasksListMode:
		return m.taskListUpdate(msg)
	case createMode:
		return m.createTaskUpdate(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.mode == createMode {
		return m.taskInputModel.View()
	}
	return docStyle.Render(m.list.View())
}

func (m model) taskListUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "a":
			m.taskInputModel = task_input.InitialModel()
			m.mode = createMode
			return m, nil
		case "x":
			if len(m.list.Items()) > 0 {
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) createTaskUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	taskInputModel, cmd := m.taskInputModel.Update(msg)
	m.taskInputModel = taskInputModel.(task_input.TaskInput)
	if m.taskInputModel.Submitted {
		title, description := m.taskInputModel.GetNewTask()
		newTask, _ := m.taskStore.Create(title, description)
		m.list.InsertItem(len(m.list.Items()), *newTask)
		m.taskInputModel.Submitted = false
		m.mode = tasksListMode
	}
	return m, cmd
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
				key.WithHelp("x", "delete item"),
			),
		}
	}
	m := model{mode: tasksListMode, list: tasksList, taskStore: taskStore}
	m.list.Title = "Todoaroo list"

	// Initialise Bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run application
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

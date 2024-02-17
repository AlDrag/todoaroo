package task

type Task struct {
	id, title, description string
}

func (i Task) Title() string       { return i.title }
func (i Task) Description() string { return i.description }
func (i Task) FilterValue() string { return i.title }

func Create(id string, title string, description string) Task {
	return Task{
		id:          id,
		title:       title,
		description: description,
	}
}

package task

type Task struct {
	id                 int
	title, description string
}

func (i Task) ID() int             { return i.id }
func (i Task) Title() string       { return i.title }
func (i Task) Description() string { return i.description }
func (i Task) FilterValue() string { return i.title }

package app

// App is one application on a device: its identifier, a human name, and a note carrying its version or process id.
type App struct {
	id   string
	name string
	note string
}

type Input struct {
	Id   string
	Name string
	Note string
}

func New(input Input) App {
	return App{id: input.Id, name: input.Name, note: input.Note}
}

func (app App) Id() string {
	return app.id
}

func (app App) Name() string {
	return app.name
}

func (app App) Note() string {
	return app.note
}

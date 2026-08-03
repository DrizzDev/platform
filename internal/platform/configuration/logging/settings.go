package logging

type Settings struct {
	level Level
}

func (settings Settings) Level() Level {
	return settings.level
}

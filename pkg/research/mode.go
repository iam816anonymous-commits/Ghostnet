package research

type Mode struct {
	Enabled bool
}

func (m *Mode) EvaluateStrategy(name string) string {
	if !m.Enabled {
		return "Research mode disabled."
	}
	return "Evaluating experimental strategy: " + name
}

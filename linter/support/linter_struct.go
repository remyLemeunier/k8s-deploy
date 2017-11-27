package support

import "fmt"

type HelmDeployLinter struct {
	Messages []Message
	ChartDir string
}

type Message struct {
	Path string
	Err  error
}

func (l *HelmDeployLinter) RunLinterRule(path string, err error) bool {
	if err != nil {
		l.Messages = append(l.Messages, Message{path, err})

	}
	return err == nil
}

func (m Message) Error() string {
	return fmt.Sprintf("%s: %s", m.Path, m.Err.Error())
}

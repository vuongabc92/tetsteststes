package lang

import (
	"github.com/vuongabc92/octocv/config"
	goMsg "golang.org/x/text/message"
)

type message struct {
}

func NewMessage() message {
	return message{}
}

// Set string message
func (m message) Load() error {
	text := NewText()
	for _, l := range config.Languages {
		texts := text.GetTexts(l)
		for k, t := range texts {
			goMsg.SetString(l, k, t)
		}

	}

	return nil
}

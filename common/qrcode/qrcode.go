package qrcode

import (
	"bytes"

	"github.com/mdp/qrterminal/v3"
)

func New(data string) string {
	var buf bytes.Buffer
	qrterminal.Generate(data, qrterminal.M, &buf)
	return buf.String()
}

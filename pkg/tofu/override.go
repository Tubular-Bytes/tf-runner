package tofu

import (
	"bytes"
	"text/template"
)

var overrideTemplate = `terraform {
  backend "http" {
    address = "{{ .stateURL }}"
    lock_address = "{{ .lockURL }}"
    unlock_address = "{{ .unlockURL }}"
  }
}`

func Render(stateURL, lockURL, unlockURL string) ([]byte, error) {
	tpl, err := template.New("override").Parse(overrideTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	tpl.Execute(buf, map[string]string{
		"stateURL":  stateURL,
		"lockURL":   lockURL,
		"unlockURL": unlockURL,
	})

	return buf.Bytes(), nil
}

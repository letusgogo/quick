package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func String[T any](m *T) string {
	b, err := json.Marshal(*m)
	if err != nil {
		return fmt.Sprintf("%+v", *m)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *m)
	}
	return out.String()
}

func StringSlice[T any](m []*T) string {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("%+v", m)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", m)
	}
	return out.String()
}

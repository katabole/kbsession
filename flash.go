package kbsession

import (
	"encoding/gob"
	"net/http"
)

func init() {
	// gorilla/sessions uses gob to encode/decode session values, and requires us to register this type we're using to
	// store flash data.
	gob.Register(map[string][]string{})
}

const flashKey = "_flash_"

// AddFlash adds a flash message to the session for display at the next page render.
// The key groups together messages of a similar category and will be interpeted by the template.
// It's common to use keys like "success", "info", "warning", and "error" which map to CSS classes.
// Value is the string to be displayed.
func AddFlash(r *http.Request, key, value string) {
	s := Get(r)
	if flashMap, ok := s.Values[flashKey]; ok && flashMap != nil {
		flashMap := flashMap.(map[string][]string)
		flashMap[key] = append(flashMap[key], value)
	} else {
		s.Values[flashKey] = map[string][]string{key: []string{value}}
	}
}

// Flash grabs the flash messages from the session and removes them so they'll only be rendered once.
func Flash(r *http.Request) map[string][]string {
	s := Get(r)
	if flashMap, ok := s.Values[flashKey]; ok && flashMap != nil {
		delete(s.Values, flashKey)
		return flashMap.(map[string][]string)
	} else {
		return map[string][]string{}
	}
}

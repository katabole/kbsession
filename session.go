package kbsession

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionContextKey int

const sessionKey sessionContextKey = 0

type Handler struct {
	sessionStore sessions.Store
	next         http.Handler
}

func (sh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := sh.sessionStore.Get(r, "RootSession")
	if err != nil {
		// This should be extremely rare, it only happens if there's a session that can't be decoded.
		// If there's no existing session yet it's not an error, Get just returns a new one.
		slog.Error("Failed to load session", "err", err)
		http.Error(w, "Failed to load session, check logs for details", http.StatusInternalServerError)
		return
	}
	sh.next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), sessionKey, session)))
}

// NewHandler decodes the session from the request and adds it to the request context, guaranteeing it'll be
// there and making gorilla/sessions easier to work with.
func NewHandler(sessionStore sessions.Store, next http.Handler) http.Handler {
	return &Handler{sessionStore: sessionStore, next: next}
}

// Get returns the session from the request context.
func Get(r *http.Request) *sessions.Session {
	return r.Context().Value(sessionKey).(*sessions.Session)
}

// SaveSession saves the final session in the request if it's been accessed (i.e. new or modified).
func Save(w http.ResponseWriter, r *http.Request) {
	s := Get(r)

	// Avoid unnecessarily saving a session if the request didn't come with one and nothing was added to it.
	if s.IsNew && len(s.Values) == 0 {
		return
	}

	if err := s.Save(r, w); err != nil {
		slog.Error("Failed to save session", "err", err)
	}
}

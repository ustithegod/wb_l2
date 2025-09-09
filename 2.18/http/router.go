package http

import (
	"calendar/event"
	"net/http"
)

func NewRouter(repo event.IEventRepository) *http.ServeMux {
	mux := http.NewServeMux()
	h := NewEventHandler(repo)

	// mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.Get)
	mux.HandleFunc("/events_for_week", h.Get)
	mux.HandleFunc("/events_for_montj", h.Get)

	return mux
}

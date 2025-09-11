package http

import (
	"calendar/event"
	"calendar/http/middleware"
	"net/http"
)

func NewRouter(repo event.IEventRepository) *http.ServeMux {
	mux := http.NewServeMux()
	h := NewEventHandler(repo)

	mux.HandleFunc("/create_event", middleware.Logger(h.CreateEvent))
	mux.HandleFunc("/update_event", middleware.Logger(h.UpdateEvent))
	mux.HandleFunc("/delete_event", middleware.Logger(h.DeleteEvent))
	mux.HandleFunc("/events_for_day", middleware.Logger(h.Get))
	mux.HandleFunc("/events_for_week", middleware.Logger(h.Get))
	mux.HandleFunc("/events_for_month", middleware.Logger(h.Get))

	return mux
}

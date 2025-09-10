package http

import (
	"calendar/event"
	"calendar/http/logger"
	"net/http"
)

func NewRouter(repo event.IEventRepository) *http.ServeMux {
	mux := http.NewServeMux()
	h := NewEventHandler(repo)

	mux.HandleFunc("/create_event", logger.Logger(h.CreateEvent))
	mux.HandleFunc("/update_event", logger.Logger(h.UpdateEvent))
	mux.HandleFunc("/delete_event", logger.Logger(h.DeleteEvent))
	mux.HandleFunc("/events_for_day", logger.Logger(h.Get))
	mux.HandleFunc("/events_for_week", logger.Logger(h.Get))
	mux.HandleFunc("/events_for_month", logger.Logger(h.Get))

	return mux
}

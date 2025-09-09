package http

import (
	"calendar/event"
	"net/http"
	"strconv"
	"time"
)

type EventHandler struct {
	repo event.IEventRepository
}

func NewEventHandler(repo event.IEventRepository) *EventHandler {
	return &EventHandler{
		repo: repo,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, r, http.StatusBadRequest, "only POST method is allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error while parsing form")
		return
	}

	uidInput := r.FormValue("user_id")
	user_id, err := strconv.Atoi(uidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error while parsing 'user_id'")
		return
	}

	eidInput := r.FormValue("event_id")
	event_id, err := strconv.Atoi(eidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error while parsing 'event_id'")
		return
	}

	dateInput := r.FormValue("date")
	date, err := time.Parse(time.RFC3339, dateInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error while parsing 'date'. use RFC3339 format")
		return
	}

	text := r.FormValue("text")
	if text == "" {
		errorResponse(w, r, http.StatusBadRequest, "no text provided")
		return
	}

	e := event.Event{
		ID:   uint64(event_id),
		Date: date,
		Text: text,
	}

	err = h.repo.Create(uint64(user_id), e)
	if err != nil {
		errorResponse(w, r, event.GetStatusCode(err), "error while creating event")
	}
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, r, http.StatusBadRequest, "only GET method is allowed")
		return
	}

	uidInput := r.URL.Query().Get("user_id")
	user_id, err := strconv.Atoi(uidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error parsing 'user_id'")
		return
	}

	dateInput := r.URL.Query().Get("date")
	date, err := time.Parse(time.RFC3339, dateInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, "error parsing 'date'. use RFC3339 format")
		return
	}

	events := make([]event.Event, 0)
	switch r.URL.Path {
	case "/events_for_day":
		events, err = h.repo.GetForDay(uint64(user_id), date)
	case "/events_for_week":
		events, err = h.repo.GetForWeek(uint64(user_id), date)
	case "/events_for_month":
		events, err = h.repo.GetForMonth(uint64(user_id), date)
	}

	if err != nil {
		errorResponse(w, r, event.GetStatusCode(err), "error while fetching events")
	}

	writeJson(w, http.StatusOK, envelope{"result": events}, nil)
}

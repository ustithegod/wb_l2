package http

import (
	"calendar/event"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

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
		errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "only POST method is allowed")
		return
	}

	user_id, e, err := parseEvent(r)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "error while parsing form")
		return
	}

	err = h.repo.Create(uint64(user_id), e)
	if err != nil {
		errorResponse(w, r, event.GetStatusCode(err), err, "error while creating event")
	}
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "only POST method is allowed")
		return
	}

	user_id, e, err := parseEvent(r)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, err.Error())
		return
	}

	err = h.repo.Update(uint64(user_id), e)
	if err != nil {
		errorResponse(w, r, event.GetStatusCode(err), err, "error while updating event")
	}
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "only POST method is allowed")
		return
	}

	uidInput := r.URL.Query().Get("user_id")
	user_id, err := strconv.Atoi(uidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "error parsing 'user_id'")
		return
	}

	eidInput := r.URL.Query().Get("event_id")
	event_id, err := strconv.Atoi(eidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "error parsing 'event_id'")
		return
	}

	err = h.repo.Delete(uint64(user_id), uint64(event_id))
	if err != nil {
		errorResponse(w, r, event.GetStatusCode(err), err, "error while deleting event")
	}
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "only GET method is allowed")
		return
	}

	uidInput := r.URL.Query().Get("user_id")
	user_id, err := strconv.Atoi(uidInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "error parsing 'user_id'")
		return
	}

	dateInput := r.URL.Query().Get("date")
	date, err := time.Parse(dateFormat, dateInput)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "error parsing 'date'. use 'YYYY-MM-DD' format")
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
		errorResponse(w, r, event.GetStatusCode(err), err, "error while fetching events")
		return
	}

	writeJson(w, http.StatusOK, envelope{"result": events}, nil)
}

func parseEvent(r *http.Request) (int, event.Event, error) {
	err := r.ParseForm()
	if err != nil {
		return -1, event.Event{}, fmt.Errorf("error while parsing form")
	}

	uidInput := r.FormValue("user_id")
	user_id, err := strconv.Atoi(uidInput)
	if err != nil {
		return -1, event.Event{}, fmt.Errorf("error while parsing 'user_id'")
	}

	eidInput := r.FormValue("event_id")
	event_id, err := strconv.Atoi(eidInput)
	if err != nil {
		return -1, event.Event{}, fmt.Errorf("error while parsing 'event_id'")
	}

	dateInput := r.FormValue("date")
	date, err := time.Parse(dateFormat, dateInput)
	if err != nil {
		return -1, event.Event{}, fmt.Errorf("error while parsing 'date'. use 'YYYY-MM-DD' format")
	}

	text := r.FormValue("text")
	if text == "" {
		return -1, event.Event{}, fmt.Errorf("no text provided")
	}

	e := event.Event{
		ID:   uint64(event_id),
		Date: date,
		Text: text,
	}

	return user_id, e, nil
}

package event

import (
	"time"
)

type IEventRepository interface {
	Create(user_id uint64, e Event) error
	Update(user_id uint64, e Event) error
	Delete(user_id uint64, event_id uint64) error
	GetForDay(user_id uint64, day time.Time) ([]Event, error)
	GetForWeek(user_id uint64, week time.Time) ([]Event, error)
	GetForMonth(user_id uint64, month time.Time) ([]Event, error)
}

type Event struct {
	ID   uint64    `json:"event_id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

type User struct {
	ID     uint64
	Events map[uint64]Event
}

func (u *User) AddEvent(e Event) {
	u.Events[e.ID] = e
}

func (u *User) GetAllEvents() []Event {
	events := make([]Event, 0, len(u.Events))
	for _, event := range u.Events {
		events = append(events, event)
	}

	return events
}

func (u *User) GetEventByID(event_id uint64) (Event, error) {
	event, ok := u.Events[event_id]
	if ok {
		return event, nil
	}

	return Event{}, ErrServiceUnavailable
}

func (u *User) DeleteEvent(event_id uint64) error {
	_, err := u.GetEventByID(event_id)
	if err != nil {
		return err
	}

	delete(u.Events, event_id)
	return nil
}

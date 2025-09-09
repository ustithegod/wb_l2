package event

import "time"

type EventRepository struct {
	Store map[uint64]User
}

func NewRepository() IEventRepository {
	return &EventRepository{
		Store: make(map[uint64]User),
	}
}

func (r *EventRepository) Create(user_id uint64, e Event) error {
	if _, ok := r.Store[user_id]; !ok {
		r.Store[user_id] = User{
			ID:     user_id,
			Events: make(map[uint64]Event),
		}
	}
	user := r.Store[user_id]
	_, err := user.GetEventByID(e.ID)
	if err != nil {
		user.AddEvent(e)
		return nil
	}
	return ErrServiceUnavailable
}

func (r *EventRepository) Update(user_id uint64, e Event) error {
	user, ok := r.Store[user_id]
	if !ok {
		return ErrServiceUnavailable
	}
	_, err := user.GetEventByID(e.ID)
	if err != nil {
		return err
	}
	user.AddEvent(e)
	return nil
}

func (r *EventRepository) Delete(user_id uint64, event_id uint64) error {
	user, ok := r.Store[user_id]
	if !ok {
		return ErrServiceUnavailable
	}

	return user.DeleteEvent(event_id)
}

func (r *EventRepository) GetForDay(user_id uint64, day time.Time) ([]Event, error) {
	user, ok := r.Store[user_id]
	if !ok {
		return nil, ErrServiceUnavailable
	}

	var result []Event
	events := user.GetAllEvents()
	for _, event := range events {
		if event.Date.Year() == day.Year() && event.Date.Month() == day.Month() && event.Date.Day() == day.Day() {
			result = append(result, event)
		}
	}

	return result, nil
}

func (r *EventRepository) GetForWeek(user_id uint64, week time.Time) ([]Event, error) {
	user, ok := r.Store[user_id]
	if !ok {
		return nil, ErrServiceUnavailable
	}

	var result []Event
	events := user.GetAllEvents()
	for _, event := range events {
		eY, eW := event.Date.ISOWeek()
		cY, cW := week.ISOWeek()

		if eY == cY && eW == cW {
			result = append(result, event)
		}
	}
	return result, nil
}

func (r *EventRepository) GetForMonth(user_id uint64, month time.Time) ([]Event, error) {
	user, ok := r.Store[user_id]
	if !ok {
		return nil, ErrServiceUnavailable
	}

	var result []Event
	events := user.GetAllEvents()
	for _, event := range events {
		eY := event.Date.Year()
		eM := event.Date.Month()
		cY := month.Year()
		cM := month.Month()

		if eY == cY && eM == cM {
			result = append(result, event)
		}
	}
	return result, nil
}

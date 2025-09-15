package event

import (
	"errors"
	"testing"
	"time"
)

func TestEventRepository_Create(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	event := Event{ID: 1, Date: time.Now()}

	err := repo.Create(userID, event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	repoImpl := repo.(*EventRepository)
	user, exists := repoImpl.Store[userID]
	if !exists {
		t.Error("Expected user to be created in store")
	}
	if e, err := user.GetEventByID(event.ID); err != nil || e.ID != event.ID {
		t.Errorf("Expected event with ID %d, got error %v or event %v", event.ID, err, e)
	}

	err = repo.Create(userID, event)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}
}

func TestEventRepository_Update(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	event := Event{ID: 1, Date: time.Now()}

	err := repo.Update(userID, event)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}

	err = repo.Create(userID, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	newEvent := Event{ID: 2, Date: time.Now()}
	err = repo.Update(userID, newEvent)
	if err == nil {
		t.Error("Expected error for non-existent event, got nil")
	}

	updatedEvent := Event{ID: 1, Date: time.Now().Add(time.Hour)}
	err = repo.Update(userID, updatedEvent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	repoImpl := repo.(*EventRepository)
	user, _ := repoImpl.Store[userID]
	if e, err := user.GetEventByID(updatedEvent.ID); err != nil || e.Date != updatedEvent.Date {
		t.Errorf("Expected updated event with date %v, got error %v or event %v", updatedEvent.Date, err, e)
	}
}

func TestEventRepository_Delete(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	event := Event{ID: 1, Date: time.Now()}

	err := repo.Delete(userID, event.ID)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}

	err = repo.Create(userID, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	err = repo.Delete(userID, 2)
	if err == nil {
		t.Error("Expected error for non-existent event, got nil")
	}

	err = repo.Delete(userID, event.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	repoImpl := repo.(*EventRepository)
	user, _ := repoImpl.Store[userID]
	if _, err := user.GetEventByID(event.ID); err == nil {
		t.Error("Expected event to be deleted")
	}
}

func TestEventRepository_GetForDay(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	day := time.Date(2025, time.September, 11, 0, 0, 0, 0, time.UTC)
	event1 := Event{ID: 1, Date: day}
	event2 := Event{ID: 2, Date: day.Add(24 * time.Hour)}

	events, err := repo.GetForDay(userID, day)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}

	err = repo.Create(userID, event1)
	if err != nil {
		t.Fatalf("Failed to create event1: %v", err)
	}
	err = repo.Create(userID, event2)
	if err != nil {
		t.Fatalf("Failed to create event2: %v", err)
	}

	events, err = repo.GetForDay(userID, day)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(events) != 1 || events[0].ID != event1.ID {
		t.Errorf("Expected one event with ID %d, got %v", event1.ID, events)
	}
}

func TestEventRepository_GetForWeek(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	week := time.Date(2025, time.September, 8, 0, 0, 0, 0, time.UTC) // Start of week
	event1 := Event{ID: 1, Date: week}
	event2 := Event{ID: 2, Date: week.Add(7 * 24 * time.Hour)} // Next week

	events, err := repo.GetForWeek(userID, week)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}

	err = repo.Create(userID, event1)
	if err != nil {
		t.Fatalf("Failed to create event1: %v", err)
	}
	err = repo.Create(userID, event2)
	if err != nil {
		t.Fatalf("Failed to create event2: %v", err)
	}

	events, err = repo.GetForWeek(userID, week)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(events) != 1 || events[0].ID != event1.ID {
		t.Errorf("Expected one event with ID %d, got %v", event1.ID, events)
	}
}

func TestEventRepository_GetForMonth(t *testing.T) {
	repo := NewRepository()
	userID := uint64(1)
	month := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	event1 := Event{ID: 1, Date: month}
	event2 := Event{ID: 2, Date: month.AddDate(0, 1, 0)} // Next month

	events, err := repo.GetForMonth(userID, month)
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}

	err = repo.Create(userID, event1)
	if err != nil {
		t.Fatalf("Failed to create event1: %v", err)
	}
	err = repo.Create(userID, event2)
	if err != nil {
		t.Fatalf("Failed to create event2: %v", err)
	}

	events, err = repo.GetForMonth(userID, month)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(events) != 1 || events[0].ID != event1.ID {
		t.Errorf("Expected one event with ID %d, got %v", event1.ID, events)
	}
}

package accounting

import "pas/events"

// inMemoryEventStorage test in memory event storage
//
//
type inMemoryEventStorage struct {
	events []events.Event
}

// AddEvent add an event to storage  (see EventStorage::AddEvent)
//
//
func (s *inMemoryEventStorage) AddEvent(event events.Event) {
	if s.events == nil {
		s.events = []events.Event{}
	}

	s.events = append(s.events, event)
}

// GetEventStream get stream of events (see EventStorage::GetEventStream)
//
//
func (s *inMemoryEventStorage) GetEventStream() chan events.Event {
	ch := make(chan events.Event)

	go func() {
		defer close(ch)

		for _, event := range s.events {
			ch <- event
		}
	}()

	return ch
}

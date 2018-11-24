package accounting

// inMemoryEventStorage test in memory event storage
//
//
type inMemoryEventStorage struct {
	events []Event
}

// AddEvent add an event to storage  (see EventStorage::AddEvent)
//
//
func (s *inMemoryEventStorage) AddEvent(event Event) {
	if s.events == nil {
		s.events = []Event{}
	}

	s.events = append(s.events, event)
}

// GetEventStream get stream of events (see EventStorage::GetEventStream)
//
//
func (s *inMemoryEventStorage) GetEventStream() chan Event {
	ch := make(chan Event)

	go func() {
		defer close(ch)

		for _, event := range s.events {
			ch <- event
		}
	}()

	return ch
}

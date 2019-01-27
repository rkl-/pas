package events

// InMemoryEventStorage simple in memory event storage
//
//
type InMemoryEventStorage struct {
	events []Event
}

func (s *InMemoryEventStorage) AddEvent(event Event) {
	if s.events == nil {
		s.events = []Event{}
	}

	s.events = append(s.events, event)
}

func (s *InMemoryEventStorage) GetEventStream() chan Event {
	ch := make(chan Event)

	go func() {
		defer close(ch)

		for _, event := range s.events {
			ch <- event
		}
	}()

	return ch
}

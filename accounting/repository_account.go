package accounting

import (
	"github.com/satori/go.uuid"
	"pas/events"
)

// AccountRepository
//
//
type AccountRepository struct {
	eventStorage events.EventStorage
}

func (r *AccountRepository) loadById(id uuid.UUID) (*Account, error) {
	account := &Account{id: id}

	// we need this to check if the AccountCreatedEvent event is
	// the first one in the stream.
	gotExpectedFirstEvent := false

	for event := range r.getHistoryFor(id) {
		if !gotExpectedFirstEvent {
			if _, ok := event.(*AccountCreatedEvent); !ok {
				return nil, &AccountCreatedEventNotFoundError{}
			}
			gotExpectedFirstEvent = true
		}

		switch event.(type) {
		//
		// AccountCreatedEvent
		//
		case *AccountCreatedEvent:
			account.title = event.(*AccountCreatedEvent).accountTitle
			account.balance = Money{}.NewFromInt(0, event.(*AccountCreatedEvent).currencyId)
			break

		//
		// AccountValueAddedEvent
		//
		case *AccountValueAddedEvent:
			value := event.(*AccountValueAddedEvent).value
			reason := event.(*AccountValueAddedEvent).reason

			if err := account.addValue(value, reason); err != nil {
				return nil, err
			}
			break

		//
		// AccountValueSubtractedEvent
		//
		case *AccountValueSubtractedEvent:
			value := event.(*AccountValueSubtractedEvent).value
			reason := event.(*AccountValueSubtractedEvent).reason

			if err := account.subtractValue(value, reason); err != nil {
				return nil, err
			}
			break

		//
		// AccountValueTransferredEvent
		//
		case *AccountValueTransferredEvent:
			fromId := event.(*AccountValueTransferredEvent).fromId
			value := event.(*AccountValueTransferredEvent).value
			reason := event.(*AccountValueTransferredEvent).reason

			if fromId == id {
				if err := account.subtractValue(value, reason); err != nil {
					return nil, err
				}
			} else {
				if err := account.addValue(value, reason); err != nil {
					return nil, err
				}
			}

			break
		}
	}

	return account, nil

}

func (r *AccountRepository) getHistoryFor(accountId uuid.UUID) chan events.Event {
	ch := make(chan events.Event)

	go func() {
		defer close(ch)

		for event := range r.eventStorage.GetEventStream() {
			switch event.(type) {
			case SingleAccountEvent:
				if event.(SingleAccountEvent).GetAccountId() == accountId {
					ch <- event
				}
				break
			case *AccountValueTransferredEvent:
				if event.(*AccountValueTransferredEvent).fromId == accountId ||
					event.(*AccountValueTransferredEvent).toId == accountId {
					ch <- event
				}
				break
			}
		}
	}()

	return ch
}

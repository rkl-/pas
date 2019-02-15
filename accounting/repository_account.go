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

func (r *AccountRepository) save(account *Account) error {
	if account.recordedEvents == nil {
		return nil
	}

	for _, event := range account.recordedEvents {
		r.eventStorage.AddEvent(event)
	}

	// reset events
	account.recordedEvents = []events.Event{}

	return nil
}

func (r *AccountRepository) loadById(id uuid.UUID) (*Account, error) {
	var account *Account

	// we need this to check if the AccountCreatedEvent event is
	// the first one in the stream.
	gotExpectedFirstEvent := false

	for event := range r.getHistoryFor(id) {
		if account == nil {
			account = &Account{id: id}
		}

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
		// PlannedCashReceiptCreatedEvent
		//
		case *PlannedCashReceiptCreatedEvent:
			plannedReceipt := &PlannedCashFlow{
				id:        event.(*PlannedCashReceiptCreatedEvent).ReceiptId,
				accountId: event.(*PlannedCashReceiptCreatedEvent).AccountId,
				date:      event.(*PlannedCashReceiptCreatedEvent).Date,
				amount:    event.(*PlannedCashReceiptCreatedEvent).Amount,
				title:     event.(*PlannedCashReceiptCreatedEvent).Title,
			}
			if err := account.addPlannedCashReceipt(plannedReceipt); err != nil {
				return nil, err
			}

			break

		case *PlannedCashReceiptConfirmedEvent:
			value := event.(*PlannedCashReceiptConfirmedEvent).Amount
			reason := event.(*PlannedCashReceiptConfirmedEvent).Title

			if err := account.addValue(value, reason); err != nil {
				return nil, err
			}

			break

		case *PlannedCashWithdrawalConfirmedEvent:
			value := event.(*PlannedCashWithdrawalConfirmedEvent).Amount
			reason := event.(*PlannedCashWithdrawalConfirmedEvent).Title

			if err := account.subtractValue(value, reason); err != nil {
				return nil, err
			}

			break

		//
		// PlannedCashWithdrawalCreatedEvent
		//
		case *PlannedCashWithdrawalCreatedEvent:
			plannedWithdrawal := &PlannedCashFlow{
				id:        event.(*PlannedCashWithdrawalCreatedEvent).WithdrawalId,
				accountId: event.(*PlannedCashWithdrawalCreatedEvent).AccountId,
				date:      event.(*PlannedCashWithdrawalCreatedEvent).Date,
				amount:    event.(*PlannedCashWithdrawalCreatedEvent).Amount,
				title:     event.(*PlannedCashWithdrawalCreatedEvent).Title,
			}
			if err := account.addPlannedCashWithdrawal(plannedWithdrawal); err != nil {
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

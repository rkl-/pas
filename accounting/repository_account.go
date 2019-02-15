package accounting

import (
	"github.com/satori/go.uuid"
	"pas/accounting/errors"
	events2 "pas/accounting/events"
	"pas/accounting/structs"
	"pas/events"
	"pas/money"
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
			if _, ok := event.(*events2.AccountCreatedEvent); !ok {
				return nil, &errors.AccountCreatedEventNotFoundError{}
			}
			gotExpectedFirstEvent = true
		}

		switch event.(type) {
		//
		// AccountCreatedEvent
		//
		case *events2.AccountCreatedEvent:
			account.title = event.(*events2.AccountCreatedEvent).AccountTitle
			account.balance = money.Money{}.NewFromInt(0, event.(*events2.AccountCreatedEvent).AurrencyId)
			break

		//
		// AccountValueAddedEvent
		//
		case *events2.AccountValueAddedEvent:
			value := event.(*events2.AccountValueAddedEvent).Value
			reason := event.(*events2.AccountValueAddedEvent).Reason

			if err := account.addValue(value, reason); err != nil {
				return nil, err
			}
			break

		//
		// AccountValueSubtractedEvent
		//
		case *events2.AccountValueSubtractedEvent:
			value := event.(*events2.AccountValueSubtractedEvent).Value
			reason := event.(*events2.AccountValueSubtractedEvent).Reason

			if err := account.subtractValue(value, reason); err != nil {
				return nil, err
			}
			break

		//
		// PlannedCashReceiptCreatedEvent
		//
		case *events2.PlannedCashReceiptCreatedEvent:
			plannedReceipt := &structs.PlannedCashFlow{
				Id:        event.(*events2.PlannedCashReceiptCreatedEvent).ReceiptId,
				AccountId: event.(*events2.PlannedCashReceiptCreatedEvent).AccountId,
				Date:      event.(*events2.PlannedCashReceiptCreatedEvent).Date,
				Amount:    event.(*events2.PlannedCashReceiptCreatedEvent).Amount,
				Title:     event.(*events2.PlannedCashReceiptCreatedEvent).Title,
			}
			if err := account.addPlannedCashReceipt(plannedReceipt); err != nil {
				return nil, err
			}

			break

		case *events2.PlannedCashReceiptConfirmedEvent:
			value := event.(*events2.PlannedCashReceiptConfirmedEvent).Amount
			reason := event.(*events2.PlannedCashReceiptConfirmedEvent).Title

			if err := account.addValue(value, reason); err != nil {
				return nil, err
			}

			break

		case *events2.PlannedCashWithdrawalConfirmedEvent:
			value := event.(*events2.PlannedCashWithdrawalConfirmedEvent).Amount
			reason := event.(*events2.PlannedCashWithdrawalConfirmedEvent).Title

			if err := account.subtractValue(value, reason); err != nil {
				return nil, err
			}

			break

		//
		// PlannedCashWithdrawalCreatedEvent
		//
		case *events2.PlannedCashWithdrawalCreatedEvent:
			plannedWithdrawal := &structs.PlannedCashFlow{
				Id:        event.(*events2.PlannedCashWithdrawalCreatedEvent).WithdrawalId,
				AccountId: event.(*events2.PlannedCashWithdrawalCreatedEvent).AccountId,
				Date:      event.(*events2.PlannedCashWithdrawalCreatedEvent).Date,
				Amount:    event.(*events2.PlannedCashWithdrawalCreatedEvent).Amount,
				Title:     event.(*events2.PlannedCashWithdrawalCreatedEvent).Title,
			}
			if err := account.addPlannedCashWithdrawal(plannedWithdrawal); err != nil {
				return nil, err
			}

			break

		//
		// AccountValueTransferredEvent
		//
		case *events2.AccountValueTransferredEvent:
			fromId := event.(*events2.AccountValueTransferredEvent).FromId
			value := event.(*events2.AccountValueTransferredEvent).Value
			reason := event.(*events2.AccountValueTransferredEvent).Reason

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
			case *events2.AccountValueTransferredEvent:
				if event.(*events2.AccountValueTransferredEvent).FromId == accountId ||
					event.(*events2.AccountValueTransferredEvent).ToId == accountId {
					ch <- event
				}
				break
			}
		}
	}()

	return ch
}

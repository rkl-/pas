package accounting

import (
	"github.com/satori/go.uuid"
	"math/big"
	"strings"
)

// Ledger ledger for accounting
//
//
type Ledger struct {
	eventDispatcher *EventDispatcher
	eventStorage    EventStorage
}

// GetInstance create new ledger instance
//
//
func (l Ledger) New(eventDispatcher *EventDispatcher, eventStorage EventStorage) *Ledger {
	le := &Ledger{
		eventDispatcher: eventDispatcher,
		eventStorage:    eventStorage,
	}

	return le
}

// Account a ledger accountId
//
//
type Account struct {
	id      uuid.UUID
	title   string
	balance Money
}

// CreateAccount create a new ledger account and dispatch AccountCreatedEvent
//
//
func (l *Ledger) CreateAccount(title, currencyId string) *Account {
	a := &Account{
		id:      uuid.NewV4(),
		title:   title,
		balance: Money{}.NewFromInt(0, currencyId),
	}

	l.eventDispatcher.Dispatch(&AccountCreatedEvent{
		accountId:    a.id,
		accountTitle: title,
		currencyId:   currencyId,
	})

	return a
}

// TransferValue transfer value fromId one accountId toId another
//
//
func (l *Ledger) TransferValue(fromAccount, toAccount *Account, value Money, reason string) error {
	ok, err := fromAccount.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &InsufficientFoundsError{}
	}

	currencyId := fromAccount.balance.GetCurrencyId()
	oldFromAmount := fromAccount.balance.GetAmount()
	oldToAmount := toAccount.balance.GetAmount()

	newFromAmount := (&big.Int{}).Sub(oldFromAmount, value.GetAmount())
	newToAmount := (&big.Int{}).Add(oldToAmount, value.GetAmount())

	newFromBalance, err := Money{}.NewFromString(newFromAmount.String(), currencyId)
	if err != nil {
		return err
	}
	newToBalance, err := Money{}.NewFromString(newToAmount.String(), currencyId)
	if err != nil {
		return err
	}

	fromAccount.balance = newFromBalance
	toAccount.balance = newToBalance

	l.eventDispatcher.Dispatch(&AccountValueTransferredEvent{
		fromId: fromAccount.id,
		toId:   toAccount.id,
		value:  value,
		reason: reason,
	})

	return nil
}

// AddValue add new value to an account and dispatch AccountValueAddedEvent
//
//
func (l *Ledger) AddValue(toAccount *Account, value Money, reason string) error {
	if err := l.addValue(toAccount, value, reason); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(&AccountValueAddedEvent{
		accountId: toAccount.id,
		value:     value,
		reason:    reason,
	})

	return nil
}

func (l *Ledger) addValue(toAccount *Account, value Money, reason string) error {
	if strings.Compare(toAccount.balance.currencyId, value.currencyId) != 0 {
		return &UnequalCurrenciesError{}
	}

	currencyId := toAccount.balance.currencyId
	newAmount := (&big.Int{}).Add(toAccount.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	toAccount.balance = newBalance

	return nil
}

// SubtractValue subtract value from an account and dispatch AccountValueSubtractedEvent
//
//
func (l *Ledger) SubtractValue(fromAccount *Account, value Money, reason string) error {
	if err := l.subtractValue(fromAccount, value, reason); err != nil {
		return err
	}

	l.eventDispatcher.Dispatch(&AccountValueSubtractedEvent{
		accountId: fromAccount.id,
		value:     value,
		reason:    reason,
	})

	return nil
}

func (l *Ledger) subtractValue(fromAccount *Account, value Money, reason string) error {
	ok, err := fromAccount.balance.IsLowerThan(value)
	if err != nil {
		return err
	}
	if ok {
		return &InsufficientFoundsError{}
	}

	currencyId := fromAccount.balance.currencyId
	newAmount := (&big.Int{}).Sub(fromAccount.balance.amount, value.amount)

	newBalance, err := Money{}.NewFromString(newAmount.String(), currencyId)
	if err != nil {
		return err
	}

	fromAccount.balance = newBalance

	return nil
}

// LoadAccount an account by id
//
//
func (l *Ledger) LoadAccount(accountId uuid.UUID) (*Account, error) {
	account := &Account{id: accountId}

	// we need this to check if the AccountCreatedEvent event is
	// the first one in the stream.
	gotExpectedFirstEvent := false

	for event := range l.getHistoryFor(accountId) {
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

			if err := l.addValue(account, value, reason); err != nil {
				return nil, err
			}
			break

		//
		// AccountValueSubtractedEvent
		//
		case *AccountValueSubtractedEvent:
			value := event.(*AccountValueSubtractedEvent).value
			reason := event.(*AccountValueSubtractedEvent).reason

			if err := l.subtractValue(account, value, reason); err != nil {
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

			if fromId == accountId {
				if err := l.subtractValue(account, value, reason); err != nil {
					return nil, err
				}
			} else {
				if err := l.addValue(account, value, reason); err != nil {
					return nil, err
				}
			}

			break
		}
	}

	return account, nil
}

func (l *Ledger) getHistoryFor(accountId uuid.UUID) chan Event {
	ch := make(chan Event)

	go func() {
		defer close(ch)

		for event := range l.eventStorage.GetEventStream() {
			switch event.(type) {
			case *AccountCreatedEvent:
				if event.(*AccountCreatedEvent).accountId == accountId {
					ch <- event
				}
				break
			case *AccountValueAddedEvent:
				if event.(*AccountValueAddedEvent).accountId == accountId {
					ch <- event
				}
				break
			case *AccountValueSubtractedEvent:
				if event.(*AccountValueSubtractedEvent).accountId == accountId {
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

package accounting

import "github.com/satori/go.uuid"

// DefaultLedger common ledger
//
//
type Ledger interface {
	// CreateAccount create a new account in ledger.
	//
	//
	CreateAccount(title, currencyId string) (*Account, error)

	// TransferValue transfer value from one account to another.
	//
	//
	TransferValue(fromAccount, toAccount *Account, value Money, reason string) error

	// AddValue add new value to an account
	//
	//
	AddValue(account *Account, value Money, reason string) error

	// SubtractValue subtract value from an account
	//
	//
	SubtractValue(account *Account, value Money, reason string) error

	// HasAccount efficient way to check if an account exists or not.
	//
	//
	HasAccount(accountId uuid.UUID) bool

	// LoadAccount an account by id
	//
	//
	LoadAccount(accountId uuid.UUID) (*Account, error)

	// AddPlannedCashReceipt add a planned cash receipt to an account
	//
	//
	AddPlannedCashReceipt(account *Account, receipt *PlannedCashReceipt) error
}

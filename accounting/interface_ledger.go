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
	TransferValue(fromAccountId, toAccountId uuid.UUID, value Money, reason string) error

	// AddValue add new value to an account
	//
	//
	AddValue(accountId uuid.UUID, value Money, reason string) error

	// SubtractValue subtract value from an account
	//
	//
	SubtractValue(accountId uuid.UUID, value Money, reason string) error

	// LoadAccount an account by id
	//
	//
	LoadAccount(accountId uuid.UUID) (*Account, error)

	// AddPlannedCashReceipt add a planned cash receipt to an account
	//
	//
	AddPlannedCashReceipt(accountId uuid.UUID, receipt *PlannedCashFlow) error

	// ConfirmPlannedCashReceipt confirm a planned cash receipt
	//
	//
	ConfirmPlannedCashReceipt(accountId uuid.UUID, receiptId uuid.UUID) error

	// AddPlannedCashWithdrawal add a planned cash withdrawal to an account
	//
	//
	AddPlannedCashWithdrawal(accountId uuid.UUID, withdrawal *PlannedCashFlow) error

	// ConfirmPlannedCashWithdrawal confirm a planned cash withdrawal
	//
	//
	ConfirmPlannedCashWithdrawal(accountId uuid.UUID, withdrawalId uuid.UUID) error
}

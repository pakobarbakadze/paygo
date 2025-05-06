package service_test

import (
	"errors"
	"paygo/internal/domain/model"
	"paygo/internal/domain/service"
	"paygo/internal/infra/database"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockDBManager struct {
	mock.Mock
}

func (m *MockDBManager) WithTransaction(fn func(tx database.Transaction) error) error {
	args := m.Called(fn)
	// Execute the function with nil as transaction to simulate the transaction
	if fn != nil {
		_ = fn(nil)
	}
	return args.Error(0)
}

func (m *MockDBManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBManager) Migrate() error {
	args := m.Called()
	return args.Error(0)
}

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindByID(tx database.Transaction, id uuid.UUID, forUpdate bool) (*model.Account, error) {
	args := m.Called(tx, id, forUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(tx database.Transaction, account *model.Account) (*model.Account, error) {
	args := m.Called(tx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx database.Transaction, transaction *model.Transaction) error {
	args := m.Called(tx, transaction)
	// Set ID for transaction because normally the DB would do this
	if args.Error(0) == nil {
		transaction.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateLedgerEntry(tx database.Transaction, entry *model.LedgerEntry) error {
	args := m.Called(tx, entry)
	return args.Error(0)
}

func setupTest() (*service.TransferService, *MockDBManager, *MockAccountRepository, *MockTransactionRepository) {
	mockDB := new(MockDBManager)
	mockAccountRepo := new(MockAccountRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	transferService := service.NewTransferService(mockDB, mockAccountRepo, mockTransactionRepo)

	return transferService, mockDB, mockAccountRepo, mockTransactionRepo
}

func TestTransferServiceTransferMoneySuccess(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "active",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "USD",
	}

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockTransactionRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Transaction")).Return(nil)
	mockTransactionRepo.On("CreateLedgerEntry", mock.Anything, mock.AnythingOfType("*model.LedgerEntry")).Return(nil).Twice()
	mockAccountRepo.On("Update", mock.Anything, mock.MatchedBy(func(account *model.Account) bool {
		return account.ID == fromAccountID
	})).Return(fromAccount, nil)
	mockAccountRepo.On("Update", mock.Anything, mock.MatchedBy(func(account *model.Account) bool {
		return account.ID == toAccountID
	})).Return(toAccount, nil)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(nil)

	// Execute
	transaction, updatedFromAccount, updatedToAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, "transfer", transaction.TransactionType)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, "USD", transaction.CurrencyCode)
	assert.Equal(t, "completed", transaction.Status)
	assert.Equal(t, description, transaction.Description)

	assert.Equal(t, 400.0, updatedFromAccount.Balance)
	assert.Equal(t, 400.0, updatedFromAccount.AvailableBalance)

	assert.Equal(t, 300.0, updatedToAccount.Balance)
	assert.Equal(t, 300.0, updatedToAccount.AvailableBalance)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneySourceAccountNotFound(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	expectedErr := errors.New("account not found")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(nil, expectedErr)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneyDestinationAccountNotFound(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	expectedErr := errors.New("account not found")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(nil, expectedErr)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneySourceAccountNotActive(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "suspended",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "active",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "USD",
	}

	expectedErr := errors.New("source account is not active")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneyDestinationAccountNotActive(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "suspended",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "USD",
	}

	expectedErr := errors.New("destination account is not active")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneyCurrencyMismatch(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "active",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "EUR",
	}

	expectedErr := errors.New("currency mismatch between accounts")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneyInsufficientFunds(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 600.0 // More than available balance
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "active",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "USD",
	}

	expectedErr := errors.New("insufficient funds")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestTransferServiceTransferMoneyTransactionCreationFailed(t *testing.T) {
	// Setup
	transferService, mockDB, mockAccountRepo, mockTransactionRepo := setupTest()

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := 100.0
	description := "Test transfer"

	fromAccount := &model.Account{
		ID:               fromAccountID,
		Status:           "active",
		Balance:          500.0,
		AvailableBalance: 500.0,
		CurrencyCode:     "USD",
	}

	toAccount := &model.Account{
		ID:               toAccountID,
		Status:           "active",
		Balance:          200.0,
		AvailableBalance: 200.0,
		CurrencyCode:     "USD",
	}

	expectedErr := errors.New("failed to create transaction")

	// Expectations
	mockAccountRepo.On("FindByID", mock.Anything, fromAccountID, true).Return(fromAccount, nil)
	mockAccountRepo.On("FindByID", mock.Anything, toAccountID, true).Return(toAccount, nil)
	mockTransactionRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Transaction")).Return(expectedErr)
	mockDB.On("WithTransaction", mock.AnythingOfType("func(database.Transaction) error")).Return(expectedErr)

	// Execute
	transaction, fromAccount, toAccount, err := transferService.TransferMoney(fromAccountID, toAccountID, amount, description)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, transaction)
	assert.Nil(t, fromAccount)
	assert.Nil(t, toAccount)

	// Verify expectations
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

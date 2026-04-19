package impl

import (
	"Hades/internal/errs"
	mockLogger "Hades/internal/logger/mocks"
	"Hades/internal/models"
	mockStorage "Hades/internal/repository/mocks"
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockStorage)

	require.NotNil(t, service)
	require.Equal(t, mockLogger, service.logger)
	require.Equal(t, mockStorage, service.storage)

}

func TestValidateItem(t *testing.T) {

	validItem := models.Item{
		Type:        "income",
		Amount:      decimal.NewFromFloat(100.50),
		Date:        time.Now().UTC(),
		Category:    "Salary",
		Description: "Monthly salary",
	}

	t.Run("valid item", func(t *testing.T) {
		err := validateItem(validItem)
		require.NoError(t, err)
	})

	t.Run("missing type", func(t *testing.T) {
		item := validItem
		item.Type = "   "
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrMissingType)
	})

	t.Run("invalid type", func(t *testing.T) {
		item := validItem
		item.Type = "invalid"
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrInvalidType)
	})

	t.Run("negative amount", func(t *testing.T) {
		item := validItem
		item.Amount = decimal.NewFromFloat(-50)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrNegativeAmount)
	})

	t.Run("zero amount", func(t *testing.T) {
		item := validItem
		item.Amount = decimal.Zero
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrZeroAmount)
	})

	t.Run("amount too large", func(t *testing.T) {
		item := validItem
		item.Amount = decimal.NewFromInt(1_000_000_001)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrAmountTooLarge)
	})

	t.Run("date too old", func(t *testing.T) {
		item := validItem
		item.Date = time.Now().UTC().AddDate(-2, 0, 0)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrDateTooOld)
	})

	t.Run("date too far in future", func(t *testing.T) {
		item := validItem
		item.Date = time.Now().UTC().AddDate(2, 0, 0)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrDateTooFar)
	})

	t.Run("missing category", func(t *testing.T) {
		item := validItem
		item.Category = " "
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrMissingCategory)
	})

	t.Run("category too short", func(t *testing.T) {
		item := validItem
		item.Category = "ab"
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrCategoryTooShort)
	})

	t.Run("category too long", func(t *testing.T) {
		item := validItem
		item.Category = strings.Repeat("a", 101)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrCategoryTooLong)
	})

	t.Run("description too long", func(t *testing.T) {
		item := validItem
		item.Description = strings.Repeat("a", 1001)
		err := validateItem(item)
		require.ErrorIs(t, err, errs.ErrDescriptionTooLong)
	})

}

func TestValidateOptions(t *testing.T) {

	validOptions := models.Options{
		Type:    "expense",
		SortBy:  "amount",
		Sort:    "ASC",
		GroupBy: "category",
	}

	t.Run("valid options", func(t *testing.T) {
		err := validateOptions(&validOptions)
		require.NoError(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		opts := validOptions
		opts.Type = "something"
		err := validateOptions(&opts)
		require.ErrorIs(t, err, errs.ErrInvalidType)
	})

	t.Run("invalid sort order", func(t *testing.T) {
		opts := validOptions
		opts.Sort = "INVALID"
		err := validateOptions(&opts)
		require.ErrorIs(t, err, errs.ErrInvalidSortOrder)
	})

	t.Run("invalid sort by", func(t *testing.T) {
		opts := validOptions
		opts.SortBy = "unknown"
		err := validateOptions(&opts)
		require.ErrorIs(t, err, errs.ErrInvalidSortBy)
	})

	t.Run("invalid group by", func(t *testing.T) {
		opts := validOptions
		opts.GroupBy = "month"
		err := validateOptions(&opts)
		require.ErrorIs(t, err, errs.ErrInvalidGroupBy)
	})

	t.Run("default values are set", func(t *testing.T) {
		opts := models.Options{}
		err := validateOptions(&opts)
		require.NoError(t, err)
		require.Equal(t, "DESC", opts.Sort)
		require.Equal(t, "date", opts.SortBy)
	})

}

func TestService_CreateItem(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	service := NewService(mockLogger, mockStorage)

	validItem := models.Item{
		Type:     "income",
		Amount:   decimal.NewFromFloat(500),
		Date:     time.Now().UTC(),
		Category: "Salary",
	}

	t.Run("validateItem fails", func(t *testing.T) {
		invalid := validItem
		invalid.Type = ""
		item, err := service.CreateItem(ctx, invalid)
		require.Equal(t, models.Item{}, item)
		require.ErrorIs(t, err, errs.ErrMissingType)
	})

	t.Run("storage.CreateItem succeeds", func(t *testing.T) {
		expectedID := 42
		mockStorage.EXPECT().CreateItem(ctx, gomock.Any()).Return(expectedID, nil)
		item, err := service.CreateItem(ctx, validItem)
		require.NoError(t, err)
		require.Equal(t, expectedID, item.ID)
		require.NotZero(t, item.CreatedAt)
	})

	t.Run("storage.CreateItem returns error", func(t *testing.T) {
		dbErr := errors.New("unique violation")
		mockStorage.EXPECT().CreateItem(ctx, gomock.Any()).Return(0, dbErr)
		mockLogger.EXPECT().LogError("service — failed to save item to storage", dbErr, "layer", "service.impl")
		item, err := service.CreateItem(ctx, validItem)
		require.Equal(t, models.Item{}, item)
		require.EqualError(t, err, dbErr.Error())
	})

}

func TestService_GetItems(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)
	service := NewService(mockLogger, mockStorage)

	validOpts := models.Options{SortBy: "date"}

	t.Run("validateOptions fails", func(t *testing.T) {
		invalid := validOpts
		invalid.Type = "wrong"
		items, err := service.GetItems(ctx, invalid)
		require.Nil(t, items)
		require.ErrorIs(t, err, errs.ErrInvalidType)
	})

	t.Run("storage.GetItems succeeds", func(t *testing.T) {
		expected := []models.Item{{ID: 1}, {ID: 2}}
		mockStorage.EXPECT().GetItems(ctx, gomock.Any()).Return(expected, nil)
		items, err := service.GetItems(ctx, validOpts)
		require.NoError(t, err)
		require.Equal(t, expected, items)
	})

	t.Run("storage.GetItems fails", func(t *testing.T) {
		dbErr := errors.New("db error")
		mockStorage.EXPECT().GetItems(ctx, gomock.Any()).Return(nil, dbErr)
		mockLogger.EXPECT().LogError("service — failed to get items from storage", dbErr, "layer", "service.impl")
		items, err := service.GetItems(ctx, validOpts)
		require.Nil(t, items)
		require.EqualError(t, err, dbErr.Error())
	})

}

func TestService_UpdateItem(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)
	service := NewService(mockLogger, mockStorage)

	validItem := models.Item{
		Type:     "expense",
		Amount:   decimal.NewFromFloat(100),
		Date:     time.Now().UTC(),
		Category: "Food",
	}

	t.Run("validateItem fails", func(t *testing.T) {
		invalid := validItem
		invalid.Amount = decimal.Zero
		item, err := service.UpdateItem(ctx, 10, invalid)
		require.Equal(t, models.Item{}, item)
		require.ErrorIs(t, err, errs.ErrZeroAmount)
	})

	t.Run("storage returns sql.ErrNoRows → ErrItemNotFound", func(t *testing.T) {
		mockStorage.EXPECT().UpdateItem(ctx, 10, gomock.Any()).Return(models.Item{}, sql.ErrNoRows)
		item, err := service.UpdateItem(ctx, 10, validItem)
		require.Equal(t, models.Item{}, item)
		require.ErrorIs(t, err, errs.ErrItemNotFound)
	})

	t.Run("storage update succeeds", func(t *testing.T) {
		returned := validItem
		returned.ID = 10
		mockStorage.EXPECT().UpdateItem(ctx, 10, gomock.Any()).Return(returned, nil)
		item, err := service.UpdateItem(ctx, 10, validItem)
		require.NoError(t, err)
		require.Equal(t, returned, item)
	})

	t.Run("storage generic error", func(t *testing.T) {
		dbErr := errors.New("constraint violation")
		mockStorage.EXPECT().UpdateItem(ctx, 10, gomock.Any()).Return(models.Item{}, dbErr)
		mockLogger.EXPECT().LogError("service — failed to update item in storage", dbErr, "itemID", 10, "layer", "service.impl")
		item, err := service.UpdateItem(ctx, 10, validItem)
		require.Equal(t, models.Item{}, item)
		require.EqualError(t, err, dbErr.Error())
	})

}

func TestService_DeleteItem(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)
	service := NewService(mockLogger, mockStorage)

	itemID := 451

	t.Run("storage.DeleteItem succeeds", func(t *testing.T) {
		mockStorage.EXPECT().DeleteItem(ctx, itemID).Return(nil)
		err := service.DeleteItem(ctx, itemID)
		require.NoError(t, err)
	})

	t.Run("storage returns ErrItemNotFound", func(t *testing.T) {
		mockStorage.EXPECT().DeleteItem(ctx, itemID).Return(errs.ErrItemNotFound)
		err := service.DeleteItem(ctx, itemID)
		require.ErrorIs(t, err, errs.ErrItemNotFound)
	})

	t.Run("storage generic error", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockStorage.EXPECT().DeleteItem(ctx, itemID).Return(dbErr)
		mockLogger.EXPECT().LogError("service — failed to delete item from storage", dbErr, "itemID", itemID, "layer", "service.impl")
		err := service.DeleteItem(ctx, itemID)
		require.EqualError(t, err, dbErr.Error())
	})

}

func TestService_GetAnalytics(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)
	service := NewService(mockLogger, mockStorage)

	validOpts := models.Options{Type: "income"}

	t.Run("validateOptions fails", func(t *testing.T) {
		invalid := validOpts
		invalid.SortBy = "invalid_field"
		res, err := service.GetAnalytics(ctx, invalid)
		require.Equal(t, models.Analytics{}, res)
		require.ErrorIs(t, err, errs.ErrInvalidSortBy)
	})

	t.Run("storage.GetAnalytics succeeds", func(t *testing.T) {
		expected := models.Analytics{}
		mockStorage.EXPECT().GetAnalytics(ctx, gomock.Any()).Return(expected, nil)
		res, err := service.GetAnalytics(ctx, validOpts)
		require.NoError(t, err)
		require.Equal(t, expected, res)
	})

	t.Run("storage.GetAnalytics fails", func(t *testing.T) {
		dbErr := errors.New("analytics error")
		mockStorage.EXPECT().GetAnalytics(ctx, gomock.Any()).Return(nil, dbErr)
		mockLogger.EXPECT().LogError("service — failed to get analytics from storage", dbErr, "layer", "service.impl")
		res, err := service.GetAnalytics(ctx, validOpts)
		require.Equal(t, models.Analytics{}, res)
		require.EqualError(t, err, dbErr.Error())
	})

}

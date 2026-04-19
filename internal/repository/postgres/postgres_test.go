package postgres_test

import (
	"Hades/internal/config"
	"Hades/internal/errs"
	"Hades/internal/logger"
	"Hades/internal/models"
	"Hades/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"
	wbf "github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
)

var (
	testStorage repository.Storage
	testDB      *dbpg.DB
)

const nonExistentID = 451

var fixedDate = time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

func TestMain(m *testing.M) {

	cfg := wbf.New()

	if err := cfg.LoadEnvFiles("../../../.env"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cfg.LoadConfigFiles("../../../config.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var appCfg config.Config
	if err := cfg.Unmarshal(&appCfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testCfg := config.Storage{
		Host:               appCfg.Storage.Host,
		Port:               appCfg.Storage.Port,
		Username:           os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWORD"),
		DBName:             appCfg.Storage.DBName,
		SSLMode:            appCfg.Storage.SSLMode,
		MaxOpenConns:       appCfg.Storage.MaxOpenConns,
		MaxIdleConns:       appCfg.Storage.MaxIdleConns,
		ConnMaxLifetime:    appCfg.Storage.ConnMaxLifetime,
		QueryRetryStrategy: appCfg.Storage.QueryRetryStrategy,
	}

	logger, _ := logger.NewLogger(config.Logger{Debug: true})

	var err error
	testDB, err = repository.ConnectDB(testCfg)
	if err != nil {
		logger.LogFatal("postgres_test — failed to connect to test DB", err, "layer", "repository.postgres_test")
	}

	if err := migrate(testDB.Master); err != nil {
		logger.LogError("postgres_test — failed to run migrations", err, "layer", "repository.postgres_test")
		os.Exit(1)
	}

	testStorage = repository.NewStorage(logger, testCfg, testDB)

	exitCode := m.Run()
	testStorage.Close()
	os.Exit(exitCode)

}

func migrate(db *sql.DB) error {
	_ = goose.SetDialect("postgres")
	if err := goose.Up(db, "../../../migrations"); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}
	return nil
}

func setupTest(t *testing.T) {

	ctx := context.Background()
	_, err := testDB.Master.ExecContext(ctx, `

	TRUNCATE TABLE items 
	RESTART IDENTITY CASCADE`)

	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

}

func TestCreateItem(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	item := models.Item{
		Type:        "income",
		Amount:      decimal.RequireFromString("1234.56"),
		Date:        fixedDate,
		Category:    "salary",
		Description: "Test income item",
		CreatedAt:   time.Now().UTC(),
	}

	id, err := testStorage.CreateItem(ctx, item)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive ID, got %d", id)
	}

}

func TestGetItems(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	itemsToInsert := []models.Item{
		{
			Type:        "income",
			Amount:      decimal.RequireFromString("1000.00"),
			Date:        fixedDate,
			Category:    "salary",
			Description: "Salary April",
			CreatedAt:   time.Now().UTC(),
		},
		{
			Type:        "expense",
			Amount:      decimal.RequireFromString("250.75"),
			Date:        fixedDate.AddDate(0, 0, -5),
			Category:    "food",
			Description: "Groceries",
			CreatedAt:   time.Now().UTC(),
		},
		{
			Type:        "income",
			Amount:      decimal.RequireFromString("750.00"),
			Date:        fixedDate.AddDate(0, 0, 3),
			Category:    "freelance",
			Description: "Side project",
			CreatedAt:   time.Now().UTC(),
		},
	}

	for _, it := range itemsToInsert {
		_, err := testStorage.CreateItem(ctx, it)
		if err != nil {
			t.Fatalf("failed to insert test item: %v", err)
		}
	}

	t.Run("all_items_sorted_by_date_desc", func(t *testing.T) {
		options := models.Options{SortBy: "date", Sort: "DESC"}
		got, err := testStorage.GetItems(ctx, options)
		if err != nil {
			t.Fatalf("GetItems failed: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("expected 3 items, got %d", len(got))
		}
		if !got[0].Date.Equal(fixedDate.AddDate(0, 0, 3)) {
			t.Errorf("expected newest date first, got %v", got[0].Date)
		}
	})

	t.Run("filter_by_type_and_date", func(t *testing.T) {
		options := models.Options{
			Type:   "income",
			From:   fixedDate.AddDate(0, 0, -1),
			To:     fixedDate.AddDate(0, 0, 10),
			SortBy: "date",
			Sort:   "ASC",
		}
		got, err := testStorage.GetItems(ctx, options)
		if err != nil {
			t.Fatalf("GetItems failed: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 income items in range, got %d", len(got))
		}
	})

}

func TestGetItems_Empty(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	options := models.Options{SortBy: "date", Sort: "DESC"}

	got, err := testStorage.GetItems(ctx, options)
	if err != nil {
		t.Fatalf("GetItems empty failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 items, got %d", len(got))
	}

}

func TestGetItems_CategoryFilterOnly(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("100.00"),
		Date:      fixedDate,
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})

	options := models.Options{
		Category: "food",
		SortBy:   "id",
		Sort:     "ASC",
	}

	got, err := testStorage.GetItems(ctx, options)
	if err != nil {
		t.Fatalf("GetItems category filter failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item with category 'food', got %d", len(got))
	}

}

func TestUpdateItem(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	original := models.Item{
		Type:        "income",
		Amount:      decimal.RequireFromString("1000.00"),
		Date:        fixedDate,
		Category:    "salary",
		Description: "Old description",
		CreatedAt:   time.Now().UTC(),
	}
	id, err := testStorage.CreateItem(ctx, original)
	if err != nil {
		t.Fatal(err)
	}

	updatedInput := models.Item{
		Type:        "expense",
		Amount:      decimal.RequireFromString("1200.50"),
		Date:        fixedDate.Add(24 * time.Hour),
		Category:    "rent",
		Description: "Updated description",
	}

	updated, err := testStorage.UpdateItem(ctx, id, updatedInput)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}

	if updated.ID != id ||
		updated.Type != "expense" ||
		!updated.Amount.Equal(decimal.RequireFromString("1200.50")) ||
		updated.Description != "Updated description" {
		t.Fatalf("updated item does not match expected values: %+v", updated)
	}

}

func TestUpdateItem_NotFound(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	updatedInput := models.Item{
		Type:   "income",
		Amount: decimal.RequireFromString("999.00"),
		Date:   fixedDate,
	}

	_, err := testStorage.UpdateItem(ctx, nonExistentID, updatedInput)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

}

func TestDeleteItem(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	item := models.Item{
		Type:        "expense",
		Amount:      decimal.RequireFromString("500.00"),
		Date:        fixedDate,
		Category:    "test",
		Description: "To be deleted",
		CreatedAt:   time.Now().UTC(),
	}
	id, err := testStorage.CreateItem(ctx, item)
	if err != nil {
		t.Fatal(err)
	}

	err = testStorage.DeleteItem(ctx, id)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

}

func TestDeleteItem_NotFound(t *testing.T) {

	setupTest(t)
	ctx := context.Background()
	err := testStorage.DeleteItem(ctx, nonExistentID)
	if !errors.Is(err, errs.ErrItemNotFound) {
		t.Fatalf("expected errs.ErrItemNotFound, got %v", err)
	}

}

func TestGetAnalytics_Ungrouped(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("1000.00"),
		Date:      fixedDate,
		Category:  "salary",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("500.00"),
		Date:      fixedDate,
		Category:  "freelance",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("300.00"),
		Date:      fixedDate,
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})

	res, err := testStorage.GetAnalytics(ctx, models.Options{})
	if err != nil {
		t.Fatalf("GetAnalytics (ungrouped) failed: %v", err)
	}

	analytics, ok := res.(models.Analytics)
	if !ok {
		t.Fatalf("expected models.Analytics, got %T", res)
	}

	if analytics.Count != 3 {
		t.Errorf("expected count 3, got %d", analytics.Count)
	}

	if !analytics.TotalIncome.Equal(decimal.NewFromInt(1500)) {
		t.Errorf("expected total_income 1500, got %v", analytics.TotalIncome)
	}
	if !analytics.TotalExpense.Equal(decimal.NewFromInt(300)) {
		t.Errorf("expected total_expense 300, got %v", analytics.TotalExpense)
	}
	if !analytics.Balance.Equal(decimal.NewFromInt(1200)) {
		t.Errorf("expected balance 1200, got %v", analytics.Balance)
	}

}

func TestGetAnalytics_Ungrouped_WithFilters(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("1000.00"),
		Date:      fixedDate,
		Category:  "salary",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("300.00"),
		Date:      fixedDate.AddDate(0, 0, -10),
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})

	options := models.Options{
		From: fixedDate.AddDate(0, 0, -5),
		To:   fixedDate,
	}

	res, err := testStorage.GetAnalytics(ctx, options)
	if err != nil {
		t.Fatalf("GetAnalytics ungrouped with filters failed: %v", err)
	}

	analytics, ok := res.(models.Analytics)
	if !ok {
		t.Fatalf("expected models.Analytics, got %T", res)
	}
	if analytics.Count != 1 {
		t.Errorf("expected count 1 (only income in range), got %d", analytics.Count)
	}

}

func TestGetAnalytics_GroupedByCategory(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("1000.00"),
		Date:      fixedDate,
		Category:  "salary",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("500.00"),
		Date:      fixedDate,
		Category:  "freelance",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("300.00"),
		Date:      fixedDate,
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})

	res, err := testStorage.GetAnalytics(ctx, models.Options{GroupBy: "category"})
	if err != nil {
		t.Fatalf("GetAnalytics (grouped) failed: %v", err)
	}

	grouped, ok := res.([]models.GroupedAnalytics)
	if !ok {
		t.Fatalf("expected []models.GroupedAnalytics, got %T", res)
	}

	if len(grouped) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(grouped))
	}

	found := map[string]bool{}
	for _, g := range grouped {
		found[g.GroupKey] = true
	}
	if !found["salary"] || !found["freelance"] || !found["food"] {
		t.Errorf("missing expected categories: %v", found)
	}

}

func TestGetAnalytics_GroupedByDay(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("1000.00"),
		Date:      fixedDate,
		Category:  "salary",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("300.00"),
		Date:      fixedDate,
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("500.00"),
		Date:      fixedDate.AddDate(0, 0, 1),
		Category:  "freelance",
		CreatedAt: time.Now().UTC(),
	})

	res, err := testStorage.GetAnalytics(ctx, models.Options{GroupBy: "day"})
	if err != nil {
		t.Fatalf("GetAnalytics grouped by day failed: %v", err)
	}

	grouped, ok := res.([]models.GroupedAnalytics)
	if !ok {
		t.Fatalf("expected []models.GroupedAnalytics, got %T", res)
	}
	if len(grouped) != 2 {
		t.Fatalf("expected 2 day groups, got %d", len(grouped))
	}

}

func TestGetAnalytics_GroupedByWeek(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "income",
		Amount:    decimal.RequireFromString("1000.00"),
		Date:      fixedDate,
		Category:  "salary",
		CreatedAt: time.Now().UTC(),
	})
	_, _ = testStorage.CreateItem(ctx, models.Item{
		Type:      "expense",
		Amount:    decimal.RequireFromString("300.00"),
		Date:      fixedDate.AddDate(0, 0, 2),
		Category:  "food",
		CreatedAt: time.Now().UTC(),
	})

	res, err := testStorage.GetAnalytics(ctx, models.Options{GroupBy: "week"})
	if err != nil {
		t.Fatalf("GetAnalytics grouped by week failed: %v", err)
	}

	grouped, ok := res.([]models.GroupedAnalytics)
	if !ok {
		t.Fatalf("expected []models.GroupedAnalytics, got %T", res)
	}
	if len(grouped) != 1 {
		t.Fatalf("expected 1 week group, got %d", len(grouped))
	}

}

func TestGetAnalytics_InvalidGroupBy(t *testing.T) {

	setupTest(t)
	ctx := context.Background()
	_, err := testStorage.GetAnalytics(ctx, models.Options{GroupBy: "invalid"})
	if !errors.Is(err, errs.ErrInvalidGroupBy) {
		t.Fatalf("expected ErrInvalidGroupBy, got %v", err)
	}

}

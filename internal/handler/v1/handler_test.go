package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	mockService "Hades/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/ginext"
	"go.uber.org/mock/gomock"
)

func setupRouter(handler *Handler) *ginext.Engine {
	router := ginext.New("")
	v1 := router.Group("/api/v1")
	{
		v1.POST("/items", handler.CreateItem)
		v1.GET("/items", handler.GetItems)
		v1.PUT("/items/:id", handler.UpdateItem)
		v1.DELETE("/items/:id", handler.DeleteItem)
		v1.GET("/analytics", handler.GetAnalytics)
	}
	return router
}

func TestHandler_CreateItem(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBufferString(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("empty date", func(t *testing.T) {
		body := CreateItemDTO{
			Type:     "income",
			Amount:   decimal.NewFromInt(100),
			Date:     "",
			Category: "test",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid date", func(t *testing.T) {
		body := CreateItemDTO{
			Type:     "income",
			Amount:   decimal.NewFromInt(100),
			Date:     "invalid-date",
			Category: "test",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns error", func(t *testing.T) {
		body := CreateItemDTO{
			Type:     "income",
			Amount:   decimal.NewFromInt(100),
			Date:     "2024-01-01",
			Category: "test",
		}
		b, _ := json.Marshal(body)
		date, _ := parseTime(body.Date)
		mockService.EXPECT().CreateItem(gomock.Any(),
			models.Item{
				Type:        body.Type,
				Amount:      body.Amount,
				Date:        date,
				Category:    body.Category,
				Description: body.Description,
			}).Return(models.Item{}, errs.ErrItemNotFound)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		body := CreateItemDTO{
			Type:        "income",
			Amount:      decimal.NewFromInt(100),
			Date:        "2024-01-01",
			Category:    "test",
			Description: "desc",
		}
		b, _ := json.Marshal(body)
		date, _ := parseTime(body.Date)
		created := models.Item{
			ID:          123,
			Type:        body.Type,
			Amount:      body.Amount,
			Date:        date,
			Category:    body.Category,
			Description: body.Description,
		}
		mockService.EXPECT().CreateItem(gomock.Any(), models.Item{
			Type:        body.Type,
			Amount:      body.Amount,
			Date:        date,
			Category:    body.Category,
			Description: body.Description,
		}).Return(created, nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.Contains(t, w.Body.String(), `"id":123`)
	})

}

func TestHandler_DeleteItem(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("item not found", func(t *testing.T) {
		mockService.EXPECT().DeleteItem(gomock.Any(), 999).Return(errs.ErrItemNotFound)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success delete", func(t *testing.T) {
		mockService.EXPECT().DeleteItem(gomock.Any(), 123).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), models.StatusDeleted)
	})

}

func TestHandler_UpdateItem(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/123", bytes.NewBufferString(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid date", func(t *testing.T) {
		body := UpdateItemDTO{
			Type:     "income",
			Amount:   decimal.NewFromInt(100),
			Date:     "invalid-date",
			Category: "test",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/123", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns error", func(t *testing.T) {
		body := UpdateItemDTO{
			Type:     "income",
			Amount:   decimal.NewFromInt(100),
			Date:     "2024-01-01",
			Category: "test",
		}
		b, _ := json.Marshal(body)
		date, _ := parseTime(body.Date)
		mockService.EXPECT().UpdateItem(gomock.Any(), 123, models.Item{
			Type:        body.Type,
			Amount:      body.Amount,
			Date:        date,
			Category:    body.Category,
			Description: body.Description,
		}).Return(models.Item{}, errs.ErrItemNotFound)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/123", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		body := UpdateItemDTO{
			Type:        "income",
			Amount:      decimal.NewFromInt(100),
			Date:        "2024-01-01",
			Category:    "test",
			Description: "updated desc",
		}
		b, _ := json.Marshal(body)
		date, _ := parseTime(body.Date)
		updated := models.Item{
			ID:          123,
			Type:        body.Type,
			Amount:      body.Amount,
			Date:        date,
			Category:    body.Category,
			Description: body.Description,
		}
		mockService.EXPECT().UpdateItem(gomock.Any(), 123, models.Item{
			Type:        body.Type,
			Amount:      body.Amount,
			Date:        date,
			Category:    body.Category,
			Description: body.Description,
		}).Return(updated, nil)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/items/123", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"id":123`)
	})

}

func TestHandler_GetItems(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	items := []models.Item{
		{ID: 1, Type: "income", Amount: decimal.NewFromInt(100),
			Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Category: "test1", Description: "desc1"},
		{ID: 2, Type: "expense", Amount: decimal.NewFromInt(-50),
			Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Category: "test2", Description: "desc2"},
	}

	t.Run("invalid query - invalid from date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items?from=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid query - invalid to date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items?to=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetItems(gomock.Any(), qp).Return(nil, errors.New("db down :("))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetItems(gomock.Any(), qp).Return(items, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, `"id":1`)
		require.Contains(t, body, `"type":"income"`)
		require.Contains(t, body, `"id":2`)
		require.Contains(t, body, `"type":"expense"`)
	})

	t.Run("export csv - items", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "csv",
		}
		mockService.EXPECT().GetItems(gomock.Any(), qp).Return(items, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items?export=csv", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "ID,Type,Amount,Date,Category,Description")
		require.Contains(t, body, "1,income,100")
		require.Contains(t, body, "2,expense,-50")
	})

	t.Run("valid from and to", func(t *testing.T) {
		fromStr := "2024-01-01"
		toStr := "2025-01-01"
		from, _ := parseTime(fromStr)
		to, _ := parseTime(toStr)
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         from,
			To:           to,
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetItems(gomock.Any(), qp).Return([]models.Item{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items?from=2024-01-01&to=2025-01-01", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("valid category", func(t *testing.T) {
		qp := models.Options{
			Category:     "food",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetItems(gomock.Any(), qp).Return([]models.Item{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/items?category=food", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

}

func TestHandler_GetAnalytics(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mockService.NewMockService(controller)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	analytics := models.Analytics{
		Count:        10,
		TotalIncome:  decimal.NewFromInt(1000),
		TotalExpense: decimal.NewFromInt(500),
		Balance:      decimal.NewFromInt(500),
		AvgAmount:    decimal.NewFromInt(100),
		Median:       decimal.NewFromInt(80),
		Percentile90: decimal.NewFromInt(200),
	}

	grouped := []models.GroupedAnalytics{
		{
			GroupKey:     "2024-01",
			Count:        5,
			TotalIncome:  decimal.NewFromInt(500),
			TotalExpense: decimal.NewFromInt(100),
			Balance:      decimal.NewFromInt(400),
			AvgAmount:    decimal.NewFromInt(80),
		},
	}

	t.Run("invalid query - invalid from date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics?from=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid query - invalid to date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics?to=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetAnalytics(gomock.Any(), qp).Return(nil, errors.New("db down :("))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "",
		}
		mockService.EXPECT().GetAnalytics(gomock.Any(), qp).Return(analytics, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, `"count":10`)
		require.Contains(t, body, `"total_income":"1000"`)
	})

	t.Run("export csv - analytics", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "csv",
		}
		mockService.EXPECT().GetAnalytics(gomock.Any(), qp).Return(analytics, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics?export=csv", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "Metric,Value")
		require.Contains(t, w.Body.String(), "Count,10")
	})

	t.Run("export csv - grouped analytics", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "month",
			ExportFormat: "csv",
		}
		mockService.EXPECT().GetAnalytics(gomock.Any(), qp).Return(grouped, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics?group_by=month&export=csv", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "GroupKey,Count,Total Income,Total Expense,Balance,Average")
		require.Contains(t, w.Body.String(), "2024-01,5,500")
	})

	t.Run("export csv - unsupported type", func(t *testing.T) {
		qp := models.Options{
			Category:     "",
			Type:         "",
			Sort:         "",
			SortBy:       "",
			From:         time.Time{},
			To:           time.Time{},
			GroupBy:      "",
			ExportFormat: "csv",
		}
		mockService.EXPECT().GetAnalytics(gomock.Any(), qp).Return("unsupported-type", nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics?export=csv", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

package grm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestRequest struct {
	Id    string `json:"id" uri:"id" form:"id"`
	Start string `json:"start" form:"start"`
	Stop  string `json:"stop" form:"stop"`
	Name  string `json:"name"`
}

func TestValidatePath_WithQueryAndJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/test/123?start=2020-01-01&stop=2025-01-01", strings.NewReader(`{"name": "test-name"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Params = gin.Params{{Key: "id", Value: "123"}}

	var received TestRequest
	handler := ValidatePath(func(c *gin.Context, req TestRequest) {
		received = req
		c.Status(http.StatusOK)
	})

	handler(c)

	if received.Id != "123" {
		t.Errorf("expected Id=123, got %s", received.Id)
	}
	if received.Start != "2020-01-01" {
		t.Errorf("expected Start=2020-01-01, got %s", received.Start)
	}
	if received.Stop != "2025-01-01" {
		t.Errorf("expected Stop=2025-01-01, got %s", received.Stop)
	}
	if received.Name != "test-name" {
		t.Errorf("expected Name=test-name, got %s", received.Name)
	}
}

func TestValidatePath_QueryOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/test/456?start=2021-06-15&stop=2024-12-31", nil)

	c.Params = gin.Params{{Key: "id", Value: "456"}}

	var received TestRequest
	handler := ValidatePath(func(c *gin.Context, req TestRequest) {
		received = req
		c.Status(http.StatusOK)
	})

	handler(c)

	if received.Id != "456" {
		t.Errorf("expected Id=456, got %s", received.Id)
	}
	if received.Start != "2021-06-15" {
		t.Errorf("expected Start=2021-06-15, got %s", received.Start)
	}
	if received.Stop != "2024-12-31" {
		t.Errorf("expected Stop=2024-12-31, got %s", received.Stop)
	}
}

func TestValidateQuery_WithPathParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/test/789?start=2020-01-01&stop=2025-01-01", nil)

	c.Params = gin.Params{{Key: "id", Value: "789"}}

	var received TestRequest
	handler := ValidateQuery(func(c *gin.Context, req TestRequest) {
		received = req
		c.Status(http.StatusOK)
	})

	handler(c)

	if received.Id != "789" {
		t.Errorf("expected Id=789, got %s", received.Id)
	}
	if received.Start != "2020-01-01" {
		t.Errorf("expected Start=2020-01-01, got %s", received.Start)
	}
	if received.Stop != "2025-01-01" {
		t.Errorf("expected Stop=2025-01-01, got %s", received.Stop)
	}
}

func TestValidateRequest_WithJSONAndPathParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/test/999", strings.NewReader(`{"name": "json-name", "start": "2022-01-01", "stop": "2023-12-31"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Params = gin.Params{{Key: "id", Value: "999"}}

	var received TestRequest
	handler := ValidateRequest(func(c *gin.Context, req TestRequest) {
		received = req
		c.Status(http.StatusOK)
	})

	handler(c)

	if received.Id != "999" {
		t.Errorf("expected Id=999, got %s", received.Id)
	}
	if received.Name != "json-name" {
		t.Errorf("expected Name=json-name, got %s", received.Name)
	}
}

func TestValidateRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/test/123", strings.NewReader(`{"name":`))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Params = gin.Params{{Key: "id", Value: "123"}}

	handler := ValidateRequest(func(c *gin.Context, req TestRequest) {
		c.Status(http.StatusOK)
	})

	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

package utils_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
)

func TestIsEmptyOrNull(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"Empty string", "", true},
		{"Whitespace string", "   ", true},
		{"Non-empty string", "test", false},
		{"String with spaces", " test ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.IsEmptyOrNull(tt.value); got != tt.want {
				t.Errorf("IsEmptyOrNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPort(t *testing.T) {
	tests := []struct {
		name         string
		envPort      string
		defaultPort  string
		expectedPort string
	}{
		{"Port from environment variable", "8080", "3000", "8080"},
		{"Empty environment variable", "", "3000", "3000"},
		{"Whitespace environment variable", "   ", "3000", "3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PORT", tt.envPort)
			defer os.Unsetenv("PORT")

			if got := utils.GetPort(tt.defaultPort); got != tt.expectedPort {
				t.Errorf("GetPort() = %v, want %v", got, tt.expectedPort)
			}
		})
	}
}

func TestGetCorrelationId(t *testing.T) {
	tests := []struct {
		name           string
		headerValue    string
		expectedID     string
		expectingError bool
	}{
		{"Valid correlation ID", "12345", "12345", false},
		{"Missing correlation ID", "", "", true},
		{"Whitespace correlation ID", "   ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			ctx.Request, _ = http.NewRequest("GET", "/", nil)
			ctx.Request.Header.Set("x-mgc-correlationId", tt.headerValue)

			got, err := utils.GetCorrelationId(ctx)
			if (err != nil) != tt.expectingError {
				t.Errorf("GetCorrelationId() error = %v, expectingError %v", err, tt.expectingError)
				return
			}
			if got != tt.expectedID {
				t.Errorf("GetCorrelationId() = %v, want %v", got, tt.expectedID)
			}
		})
	}
}

func TestGenerateGUID(t *testing.T) {
	guid1 := utils.GenerateGUID()
	guid2 := utils.GenerateGUID()

	if guid1 == "" {
		t.Errorf("GenerateGUID() returned an empty string")
	}

	if guid2 == "" {
		t.Errorf("GenerateGUID() returned an empty string")
	}

	if guid1 == guid2 {
		t.Errorf("GenerateGUID() returned the same GUID twice: %v", guid1)
	}
}

func TestToJSONReader(t *testing.T) {
	tests := []struct {
		name     string
		payload  interface{}
		expected string
	}{
		{"Valid struct", struct{ Name string }{"test"}, `{"Name":"test"}`},
		{"Empty struct", struct{}{}, `{}`},
		{"Nil payload", nil, `null`},
		{"Map payload", map[string]string{"key": "value"}, `{"key":"value"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := utils.ToJSONReader(tt.payload)
			result, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("ToJSONReader() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("ToJSONReader() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name     string
		payload  interface{}
		expected string
	}{
		{"Valid struct", struct{ Name string }{"test"}, `{"Name":"test"}`},
		{"Empty struct", struct{}{}, `{}`},
		{"Nil payload", nil, `null`},
		{"Map payload", map[string]string{"key": "value"}, `{"key":"value"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ToJSON(tt.payload)
			if result != tt.expected {
				t.Errorf("ToJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

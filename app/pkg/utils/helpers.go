package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IsEmptyOrNull checks if a given string is either empty or contains only whitespace characters.
// It returns true if the string is empty or consists solely of whitespace, and false otherwise.
//
// Parameters:
//
//	value - the string to be checked.
//
// Returns:
//
//  bool - true if the string is empty or contains only whitespace, false otherwise.

func IsEmptyOrNull(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

// GetEnvPortOrDefault retrieves the port number from the environment variable "PORT".
// If the environment variable is not set or is empty, it returns the provided default value.
//
// Parameters:
//
//	valueDefault - The default port value to return if the "PORT" environment variable is not set or is empty.
//
// Returns:
//
//	The port number as a string.
func GetEnvPortOrDefault(valueDefault string) string {
	port := os.Getenv("PORT")
	if IsEmptyOrNull(port) {
		port = valueDefault
	}
	return port
}

// GetCorrelationId retrieves the correlation ID from the request header.
// It expects the header "x-mgc-correlationId" to be present in the request.
// If the header is missing or empty, it returns an error.
//
// Parameters:
//
//	ctx (*gin.Context): The context of the incoming HTTP request.
//
// Returns:
//
//	(string, error): The correlation ID if present, otherwise an error indicating the header is missing.
func GetCorrelationId(ctx *gin.Context) (string, error) {
	correlationId := ctx.GetHeader("x-mgc-correlationId")
	if IsEmptyOrNull(correlationId) {
		return "", fmt.Errorf("missing required header: x-mgc-correlationId")
	}

	return correlationId, nil
}

// GenerateGUID generates a new globally unique identifier (GUID) as a string.
// It uses the uuid package to create a new UUID and returns it in string format.
func GenerateGUID() string {
	return uuid.New().String()
}

// ToJSONReader converts a given payload to a JSON-encoded io.Reader.
// It takes an interface{} as input, marshals it into JSON, and returns
// an io.Reader containing the JSON data.
//
// Note: This function ignores any errors that occur during JSON marshaling.
func ToJSONReader(payload interface{}) io.Reader {
	jsonData, _ := json.Marshal(payload)

	return strings.NewReader(string(jsonData))
}

// ToJSON converts a given payload to its JSON string representation.
// It takes an interface{} as input and returns a JSON string.
// Note: This function ignores any errors that occur during the marshalling process.
func ToJSON(payload interface{}) string {
	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}

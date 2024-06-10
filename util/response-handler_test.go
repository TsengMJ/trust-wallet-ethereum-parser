package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum-parser/util"
)

func TestGetSuccessResponse(t *testing.T) {
	// Test with integer data
	intResponse := util.GetSuccessResponse(42)
	assert.NotNil(t, intResponse, "Response should not be nil")
	assert.Equal(t, 42, intResponse.Data, "Data should be equal")

	// Test with string data
	strResponse := util.GetSuccessResponse("success")
	assert.NotNil(t, strResponse, "Response should not be nil")
	assert.Equal(t, "success", strResponse.Data, "Data should be equal")
	assert.Empty(t, strResponse.Error, "Error should be empty")
}

func TestGetFailResponse(t *testing.T) {
	// Test with error message
	errResponse := util.GetFailResponse("error")
	assert.NotNil(t, errResponse, "Response should not be nil")
	assert.Empty(t, errResponse.Data, "Data should be empty")
	assert.Equal(t, "error", errResponse.Error, "Error should be equal")
}

package util

import "github.com/gin-gonic/gin"

// The function returns a Gin H map with an error message as its value.
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

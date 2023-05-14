package api

import "github.com/gin-gonic/gin"

// The function returns a Gin H map with an error message as its value.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

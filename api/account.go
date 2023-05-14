package api

import (
	db "go-backend/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// The `createAccountRequest` type is a struct that represents a request to create an account with
// required fields for owner and currency, where currency must be one of CAD, USD, or EUR.
// @property {string} Owner - Owner is a property of the createAccountRequest struct. It is a string
// that represents the name of the owner of the account. The `json:"owner"` tag is used to specify the
// name of the property when the struct is serialized to JSON. The `binding:"required"` tag is used to
// @property {string} Currency - Currency is a property of the createAccountRequest struct. It is a
// string type and is tagged with `json:"currency" binding:"required,oneof=CAD USD EUR"`. This means
// that when a request is made to create an account, the currency field must be included in the request
// body
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=CAD USD EUR"`
}

// This is a function that creates a new account for a user. It receives a request with the owner's
// name and the currency of the account, and then it creates a new account with a balance of 0 using
// the `CreateAccountParams` struct from the database package. If there is an error during the creation
// of the account, it returns a 500 Internal Server Error response. Otherwise, it returns a 200 OK
// response with the newly created account.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

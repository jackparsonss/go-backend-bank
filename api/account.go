package api

import (
	"database/sql"
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

// The above code defines a struct type for a GET request to retrieve an account by its ID.
// @property {int64} ID - ID is a field of type int64 that is used to represent the unique identifier
// of an account request. It is tagged with `uri:"id"` to indicate that it should be extracted from the
// URI path of an HTTP request. It is also tagged with `binding:"required"` to indicate that it
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// This is a function that retrieves an account by its ID. It first extracts the ID from the URI path
// of the HTTP request. If there is an error during this process, it returns a 400 Bad Request response.
// Otherwise, it uses the `GetAccount` method from the database package to retrieve the account with the given ID.
// If there is an error during this process, it returns a 500 Internal Server Error response. Otherwise, it returns
// a 200 OK response with the retrieved account.
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// The deleteAccountRequest type is a struct that contains an ID field with URI binding and a minimum
// value of 1.
// @property {int64} ID - ID is a field of type int64 that represents the unique identifier of an
// account that needs to be deleted. It is tagged with `uri:"id"` to indicate that it is a parameter in
// the URI of an HTTP request and with `binding:"required,min=1"` to specify that it is
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// The `deleteAccount` function is a method of the `Server` struct that handles HTTP requests to delete
// an account. It first extracts the account ID from the URI path of the HTTP request using the
// `ShouldBindUri` method. If there is an error during this process, it returns a 400 Bad Request
// response. Otherwise, it uses the `DeleteAccount` method from the database package to delete the
// account with the given ID. If there is an error during this process, it returns a 500 Internal
// Server Error response. Otherwise, it returns a 200 OK response with a JSON message indicating that
// the account was successfully deleted.
func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully delete user"})
}

// The above type defines a request to update an account's ID and balance in a Go program.
// @property {int64} ID - ID is an integer field that represents the unique identifier of an account.
// @property {int64} Balance - The `Balance` property is an integer that represents the current balance
// of an account. It is used in the `updateAccountRequest` struct to update the balance of an account
// identified by its `ID`.
type updateAccountRequest struct {
	ID      int64 `uri:"id" binding:"required,min=1"`
	Balance int64 `json:"balance"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Balance < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid balance sent"})
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

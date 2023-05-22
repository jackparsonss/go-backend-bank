package api

import (
	"errors"
	db "go-backend/db/sqlc"
	"go-backend/token"
	"go-backend/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// The `addAccountRoutes` function is a method of the `Server` struct that adds routes for
// account-related HTTP requests to the `apiRouter` instance of the `gin.RouterGroup` type. It creates
// a new `accountRouter` instance of the `gin.RouterGroup` type with the base path of "/accounts" and
// then adds HTTP request handlers for creating, listing, retrieving, updating, and deleting accounts
// using the `createAccount`, `listAccounts`, `getAccount`, `updateAccount`, and `deleteAccount`
// methods of the `Server` struct, respectively.
func (server *Server) addAccountRoutes(apiRouter *gin.RouterGroup) {
	accountRouter := apiRouter.Group("/accounts")
	accountRouter.POST("", server.createAccount)
	accountRouter.GET("", server.listAccounts)
	accountRouter.GET("/:id", server.getAccount)
	accountRouter.PUT("/:id", server.updateAccount)
	accountRouter.DELETE("/:id", server.deleteAccount)
}

// The `createAccountRequest` type is a struct that represents a request to create an account with
// required fields for owner and currency, where currency must be one of CAD, USD, or EUR.
// @property {string} Owner - Owner is a property of the createAccountRequest struct. It is a string
// that represents the name of the owner of the account. The `json:"owner"` tag is used to specify the
// name of the property when the struct is serialized to JSON. The `binding:"required"` tag is used to
// @property {string} Currency - Currency is a property of the createAccountRequest struct. It is a
// string type and is tagged with `json:"currency" binding:"required,currency"`. This means
// that when a request is made to create an account, the currency field must be included in the request
// body
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// This is a function that creates a new account for a user. It receives a request with the owner's
// name and the currency of the account, and then it creates a new account with a balance of 0 using
// the `CreateAccountParams` struct from the database package. If there is an error during the creation
// of the account, it returns a 500 Internal Server Error response. Otherwise, it returns a 200 OK
// response with the newly created account.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if !util.CheckError(ctx, err) {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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
	Balance int64 `json:"balance" binding:"min=0"`
}

// The `updateAccount` function is a method of the `Server` struct that handles HTTP requests to update
// an account's balance. It first extracts the account ID from the URI path of the HTTP request using
// the `ShouldBindUri` method. If there is an error during this process, it returns a 400 Bad Request
// response. Otherwise, it uses the `ShouldBindJSON` method to extract the new balance from the request
// body. If there is an error during this process, it returns a 400 Bad Request response. If the new
// balance is less than 0, it returns a 400 Bad Request response with an error message. Otherwise, it
// uses the `UpdateAccount` method from the database package to update the account with the given ID
// and new balance. If there is an error during this process, it returns a 500 Internal Server Error
// response. Otherwise, it returns a 200 OK response with the updated account.
func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

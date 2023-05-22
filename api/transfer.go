package api

import (
	"errors"
	"fmt"
	db "go-backend/db/sqlc"
	"go-backend/token"
	"go-backend/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) addTransferRoutes(apiRouter *gin.RouterGroup) {
	accountRouter := apiRouter.Group("/transfers")
	accountRouter.POST("", server.createTransfer)
}

// This is a Go struct type for creating a transfer request with required fields for from and to
// account IDs, amount, and currency.
// @property {int64} FromAccountID - This property represents the ID of the account from which the
// transfer is being initiated. It is of type int64 and is required, meaning it must be provided in the
// request body. The binding tag is used to specify validation rules for this property. In this case,
// it must be a positive integer (
// @property {int64} ToAccountID - ToAccountID is an integer property that represents the ID of the
// account to which the transfer request is being made. It is a required field and must have a minimum
// value of 1.
// @property {int64} Amount - The amount property represents the amount of money that is being
// transferred from one account to another. It is of type int64, which means it can hold integer values
// up to 64 bits in size. The value of this property must be greater than zero, as specified by the
// binding tag "gt=
// @property {string} Currency - Currency is a string property that represents the currency of the
// transfer amount. It is a required field and can only have one of the three values: CAD, USD, or EUR.
type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// This is a function that handles the creation of a transfer request. It first binds the request body
// to a `createTransferRequest` struct and validates it using the `binding` tag rules. It then checks
// if the `FromAccountID` and `ToAccountID` are valid accounts with matching currencies using the
// `validAccount` function. If both accounts are valid, it creates a `db.TransferTxParams` struct with
// the necessary parameters and calls the `TransferTx` function from the `store` to execute the
// transfer transaction. If there are any errors during this process, it returns an error response with
// the appropriate status code. If the transfer is successful, it returns a success response with the
// transfer details.
func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// The `validAccount` function is a helper function that checks if an account with the given
// `accountID` and `currency` exists in the database. It takes in a `gin.Context` object, an
// `accountID` of type `int64`, and a `currency` of type `string`. It returns a boolean value
// indicating whether the account is valid or not.
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if !util.CheckError(ctx, err) {
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return account, false
	}

	return account, true
}

package pkg

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	appErr error
	status int
}

var (
	Internal   = ApiError{errors.New("internal fails"), 500}
	BadRequest = ApiError{errors.New("bad request"), 400}
	NotFound   = ApiError{errors.New("404 not found"), 404}
)

func (a *ApiError) SendErr(errorMessage string, ctx *gin.Context) {
	ctx.JSON(a.status, gin.H{"msg": errorMessage})
}

func (a *ApiError) IsErr(sevErr error ,errorMessage string, ctx *gin.Context) bool {
	if sevErr == nil {
		return false
	}	

	ctx.JSON(a.status, gin.H{"msg": errorMessage})
	return true
}

func (a *ApiError) IsErrWithCallback(sevErr error ,errorMessage string, ctx *gin.Context, callback func()) bool {
	if sevErr == nil {
		return false
	}	
	ctx.JSON(a.status, gin.H{"msg": errorMessage})
	callback()
	return true
}
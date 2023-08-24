package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

// prepareErrorResponse - function to prepare json error reply
func PrepareErrorResponse(c echo.Context, msgText, errorCode string, httpStatus int) error {
	errorResponse := models.ErrorResponse{
		Status:    consts.CErrorStatus,
		MsgText:   msgText,
		ErrorCode: errorCode,
	}

	b, err := json.Marshal(errorResponse)
	if err == nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Response().WriteHeader(httpStatus)
		io.WriteString(c.Response(), string(b[:]))
		return nil
	}
	return err
}

// prepareSuccessResponse - function to prepare json successful reply
func PrepareSuccessResponse(c echo.Context, preparedStruct interface{}) error {
	b, err := json.Marshal(preparedStruct)
	if err == nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Response().WriteHeader(http.StatusOK)
		io.WriteString(c.Response(), string(b[:]))
		return nil
	}
	return err
}

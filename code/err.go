package code

import (
	"github.com/linger1216/go-utils/common"
	"net/http"
)

var (
	ErrInvalidPara    = common.NewError(http.StatusBadRequest, "invalid para")
	ErrNotFound       = common.NewError(http.StatusNotFound, "not found")
	ErrInternalServer = common.NewError(http.StatusInternalServerError, "internal server error")
)

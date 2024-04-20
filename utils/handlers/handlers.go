package handlers

import (
	"emailserver-saimu/utils/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

func HandleCreate(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"data": data,
	})
}

func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case errors.AppError:
		c.JSON(e.Code, gin.H{
			"code":          e.Code,
			"error_message": e.Message,
		})
	case error:
		errMsg := "unexpected error"
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":          http.StatusInternalServerError,
			"error_message": errMsg,
		})
	}
}

func HandleRedirect(c *gin.Context, url string) {
	c.Redirect(http.StatusTemporaryRedirect, url)
}

package middleware

import (
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/util"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	tokenString := strings.Split(auth, " ")
	// status 401
	if len(tokenString) != 2 || tokenString[0] != "Bearer" {
		c.Error(dto.CustomError{
			ErrorStr:    constant.ErrorUnauthorized.Error(),
			InternalErr: fmt.Sprintf("Token length: %d, Token String 0: %s", len(tokenString), tokenString[0]),
			Status:      http.StatusUnauthorized,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := util.ParseToken(tokenString[1])
	if err != nil {
		if err == constant.ErrorHandleTokenExpired {
			c.Error(dto.CustomError{
				ErrorStr:    constant.ErrorUnauthorized.Error(),
				InternalErr: err.Error(),
				Status:      http.StatusUnauthorized,
			})

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Error(dto.CustomError{
			ErrorStr:    constant.ErrorUnauthorized.Error(),
			InternalErr: err.Error(),
			Status:      http.StatusUnauthorized,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	role := claims.Role
	id, err := claims.GetSubject()
	if err != nil {
		c.Error(dto.CustomError{
			ErrorStr:    constant.ErrorGetClaimSubject.Error(),
			InternalErr: err.Error(),
			Status:      http.StatusUnauthorized,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	subId, err := strconv.Atoi(id)
	if err != nil {
		c.Error(dto.CustomError{
			ErrorStr:    constant.ErrorInternalServer.Error(),
			InternalErr: err.Error(),
			Status:      http.StatusUnauthorized,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sub := entity.SubAuth{
		Role: role,
		Id:   subId,
	}

	c.Set("sub", sub)

	c.Next()
}

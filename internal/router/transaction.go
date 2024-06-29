package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkxro/squid/internal/common"
	"github.com/pkxro/squid/internal/model"
	"github.com/pkxro/squid/pkg"
)

// RegisterTransactionsRouter is a router register method that applies the routes to a router
func (m *Manager) RegisterTransactionsRouter(router *gin.RouterGroup) {
	router.POST("/signWithToken", m.SignWithTokenFee)
}

func (m *Manager) SignWithTokenFee(c *gin.Context) {
	var req model.SignWithTokenFeeRequest
	err := c.ShouldBindJSON(&req)

	out, err := m.Controller.Tx.SignWithTokenFee(c.Request.Context(), req)
	if err != nil {
		c.JSON(
			common.WrapAPIError(err.Error(),
				common.SquidBadRequestError,
				pkg.APIVersion,
			),
		)
		return
	}

	c.JSON(http.StatusOK, &out)
}

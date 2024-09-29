package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type controlPanelHandler struct {
	pageObject
	SiteInfo    siteInfo
	UserResults pagination.PagingInfo[domain.User]
	User        domain.User
	userService service.UserService
}

func NewControlPanelHandler() controlPanelHandler {
	return controlPanelHandler{
		SiteInfo:    siteInfo{},
		UserResults: pagination.PagingInfo[domain.User]{},
		userService: service.NewUserService(),
		User:        domain.User{},
	}
}

func (h controlPanelHandler) AppInsights(c *gin.Context) {
	h.Reset(c)

	persistence.DB.Model(&domain.User{}).Count(&h.SiteInfo.UserCount)
	persistence.DB.Model(&domain.Customer{}).Count(&h.SiteInfo.CustomerCount)
	persistence.DB.Model(&domain.Bank{}).Count(&h.SiteInfo.BankCount)
	persistence.DB.Model(&domain.Transaction{}).Count(&h.SiteInfo.TransactionCount)

	c.HTML(http.StatusOK, "control_panel_app_insights", h)
}

func (h controlPanelHandler) GetUsers(c *gin.Context) {
	h.Reset(c)

	pageNumber, err := strconv.Atoi(c.Query("page"))
	search := c.Query("search")

	if err != nil || pageNumber <= 0 {
		pageNumber = 1
	}

	pagingInfo := pagination.PagingInfo[domain.User]{
		PageNumber:   pageNumber,
		ItemsPerPage: 15,
	}
	h.userService.Search(search, &pagingInfo)
	h.UserResults = pagingInfo

	c.HTML(http.StatusOK, "control_panel_users", h)
}

func (h controlPanelHandler) OpenUserModal(c *gin.Context) {
	h.ModalType = "control_panel_user_modal"
	h.userService.FindByID(c.Param("id"), &h.User)

	c.HTML(http.StatusOK, "modal", h)
}

func (h controlPanelHandler) SearchUsers(c *gin.Context) {
	search := c.Query("search")

	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil || pageNumber <= 0 {
		pageNumber = 1
	}

	pagingInfo := pagination.PagingInfo[domain.User]{
		PageNumber:   pageNumber,
		ItemsPerPage: 15,
	}

	h.userService.Search(search, &pagingInfo)
	h.UserResults = pagingInfo

	c.HTML(http.StatusOK, "control_panel_users_list", h.UserResults)
}

func (h controlPanelHandler) Polls(c *gin.Context) {
	h.Reset(c)
	c.HTML(http.StatusOK, "control_panel_polls", h)
}

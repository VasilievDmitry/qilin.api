package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"qilin-api/pkg/api/rbac_echo"
	"qilin-api/pkg/mapper"
	"qilin-api/pkg/model"
	"qilin-api/pkg/orm"
	"strconv"
	"strings"
	"time"
)

type OnboardingAdminRouter struct {
	service             *orm.AdminOnboardingService
	notificationService model.NotificationService
}

type ChangeStatusRequest struct {
	Message string `json:"message"`
	Status  string `json:"status" validate:"required"`
}

type NotificationRequest struct {
	Message string `json:"message"`
	Title   string `json:"title" validate:"required"`
}

type NotificationDTO struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
	IsRead    bool   `json:"isRead"`
}

type ShortNotificationDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
	IsRead    bool   `json:"isRead"`
	HaveMsg   bool   `json:"haveMsg"`
}

type ShortDocumentsInfoDTO struct {
	VendorID  string `json:"vendor_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Person    string `json:"person"`
	UpdatedAt string `json:"updatedAt"`
	Status    string `json:"status"`
}

func InitAdminOnboardingRouter(group *echo.Group, service *orm.AdminOnboardingService, notificationService model.NotificationService) (*OnboardingAdminRouter, error) {
	router := OnboardingAdminRouter{
		service:             service,
		notificationService: notificationService,
	}

	common := []string{"*", model.AdminDocumentsType, model.VendorDomain}

	r := rbac_echo.Group(group, "/vendors", &router, common)
	r.GET("/reviews", router.getReviews, nil)
	r.GET("/:vendorId/documents", router.getDocument, nil)
	r.PUT("/:vendorId/documents/status", router.changeStatus, nil)
	r.POST("/:vendorId/messages", router.sendNotification, nil)
	r.GET("/:vendorId/messages", router.getNotifications, nil)

	return &router, nil
}

func (api *OnboardingAdminRouter) GetOwner(ctx rbac_echo.AppContext) (string, error) {
	if strings.Contains(ctx.Path(), ":vendorId") {
		return GetOwnerForVendor(ctx)
	}

	return "*", nil
}

func (api *OnboardingAdminRouter) changeStatus(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("vendorId"))
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, errors.Wrap(err, "Bad id"))
	}

	request := new(ChangeStatusRequest)

	if err := ctx.Bind(request); err != nil {
		return orm.NewServiceError(http.StatusBadRequest, err)
	}

	if errs := ctx.Validate(request); errs != nil {
		return orm.NewServiceError(http.StatusUnprocessableEntity, errs)
	}

	status, err := model.ReviewStatusFromString(request.Status)
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, errors.Wrap(err, "Bad status"))
	}

	err = api.service.ChangeStatus(id, status)
	if err != nil {
		return err
	}

	if request.Message != "" {
		_, err := api.notificationService.SendNotification(&model.Notification{Title: request.Message, Message: request.Message, VendorID: id})
		if err != nil {
			zap.L().Error(err.Error())
		}
	}

	return ctx.JSON(http.StatusOK, "")
}

func (api *OnboardingAdminRouter) getDocument(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("vendorId"))
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, errors.Wrap(err, "Bad id"))
	}

	doc, err := api.service.GetForVendor(id)
	if err != nil {
		return err
	}

	dto := DocumentsInfoResponseDTO{}
	err = mapper.Map(doc, &dto)
	if err != nil {
		return orm.NewServiceError(http.StatusInternalServerError, errors.Wrap(err, "Mapping dto error"))
	}
	dto.Status = doc.ReviewStatus.ToString()

	return ctx.JSON(http.StatusOK, dto)
}

func (api *OnboardingAdminRouter) getReviews(ctx echo.Context) error {
	offset := 0
	limit := 20
	status := model.ReviewUndefined

	if offsetParam := ctx.QueryParam("offset"); offsetParam != "" {
		if num, err := strconv.Atoi(offsetParam); err == nil {
			offset = num
		} else {
			return orm.NewServiceError(http.StatusBadRequest, errors.Wrapf(err, "Bad limit"))
		}
	}

	if limitParam := ctx.QueryParam("limit"); limitParam != "" {
		if num, err := strconv.Atoi(limitParam); err == nil {
			limit = num
		} else {
			return orm.NewServiceError(http.StatusBadRequest, errors.Wrapf(err, "Bad limit"))
		}
	}

	name := ctx.QueryParam("name")
	status, err := model.ReviewStatusFromString(ctx.QueryParam("status"))
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, errors.Wrapf(err, "Bad status"))
	}

	sort := ctx.QueryParam("sort")
	requests, count, err := api.service.GetRequests(limit, offset, name, status, sort)

	if err != nil {
		return err
	}

	var dto = make([]ShortDocumentsInfoDTO, len(requests))

	for i, doc := range requests {
		if name, ok := doc.Company["Name"]; ok {
			dto[i].Name = name.(string)
		}

		if country, ok := doc.Company["Country"]; ok {
			dto[i].Country = country.(string)
		}

		if contact, ok := doc.Contact["Authorized"]; ok {
			contactMap := contact.(map[string]interface{})
			dto[i].Person = contactMap["FullName"].(string)
		}

		dto[i].VendorID = doc.VendorID.String()
		dto[i].Status = doc.ReviewStatus.ToString()
		dto[i].UpdatedAt = doc.UpdatedAt.Format(time.RFC3339)
	}

	if dto == nil {
		dto = make([]ShortDocumentsInfoDTO, 0)
	}

	ctx.Response().Header().Add("X-Items-Count", fmt.Sprintf("%d", count))

	return ctx.JSON(http.StatusOK, dto)
}

func (api *OnboardingAdminRouter) getNotifications(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("vendorId"))
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, err)
	}

	offset := 0
	limit := 20

	if offsetParam := ctx.QueryParam("offset"); offsetParam != "" {
		if num, err := strconv.Atoi(offsetParam); err == nil {
			offset = num
		} else {
			return orm.NewServiceError(http.StatusBadRequest, errors.Wrapf(err, "Bad offset"))
		}
	}

	if limitParam := ctx.QueryParam("limit"); limitParam != "" {
		if num, err := strconv.Atoi(limitParam); err == nil {
			limit = num
		} else {
			return orm.NewServiceError(http.StatusBadRequest, errors.Wrapf(err, "Bad limit"))
		}
	}

	query := ctx.QueryParam("query")
	sort := ctx.QueryParam("sort")

	notifications, count, err := api.notificationService.GetNotifications(id, limit, offset, query, sort)
	if err != nil {
		return err
	}

	var result []NotificationDTO
	err = mapper.Map(notifications, &result)
	if err != nil {
		return orm.NewServiceErrorf(http.StatusInternalServerError, "Can't map to dto %#v", notifications)
	}

	for i, n := range notifications {
		result[i].ID = n.ID.String()
		result[i].CreatedAt = n.CreatedAt.Format(time.RFC3339)
	}

	ctx.Response().Header().Add("X-Items-Count", fmt.Sprintf("%d", count))

	return ctx.JSON(http.StatusOK, result)
}

func (api *OnboardingAdminRouter) sendNotification(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("vendorId"))
	if err != nil {
		return orm.NewServiceError(http.StatusBadRequest, err)
	}

	request := new(NotificationRequest)
	if err := ctx.Bind(request); err != nil {
		return orm.NewServiceError(http.StatusBadRequest, err)
	}

	if errs := ctx.Validate(request); errs != nil {
		return orm.NewServiceError(http.StatusUnprocessableEntity, errs)
	}

	notification, err := api.notificationService.SendNotification(&model.Notification{Message: request.Message, Title: request.Title, VendorID: id})
	if err != nil {
		return err
	}

	result := NotificationDTO{}
	err = mapper.Map(notification, &result)
	if err != nil {
		return orm.NewServiceErrorf(http.StatusInternalServerError, "Can't map to DTO `%#v`", notification)
	}

	result.CreatedAt = notification.CreatedAt.Format(time.RFC3339)
	result.ID = notification.ID.String()

	return ctx.JSON(http.StatusOK, result)
}

package api

import (
	"qilin-api/pkg/orm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"qilin-api/pkg/model"
	"github.com/satori/go.uuid"
	"github.com/labstack/echo"
	 maper "gopkg.in/jeevatkm/go-model.v1"
)

type (
	//MediaRouter is router struct
	MediaRouter struct {
		mediaService model.MediaService
	}

	//Media is DTO object with full information about media for game
	Media struct {
		
		// localized cover image of game
		CoverImage *model.LocalizedString `json:"coverImage" validate:"required"`

		// localized cover video of game
		CoverVideo *model.LocalizedString `json:"coverVideo" validate:"required"`

		// localized cover video of game
		Trailers *model.LocalizedString `json:"trailers" validate:"required"`

		// localized cover video of game
		Store *Store `json:"store" validate:"required,dive"`

		Capsule *Capsule `json:"capsule" validate:"required,dive"`
	}

	//Capsule is DTO object with information about capsule media for game
	Capsule struct {
		Generic *model.LocalizedString `json:"generic" validate:"required"`

		Small *model.LocalizedString `json:"small" validate:"required"`
	}

	//Store is DTO object with information about store media for game
	Store struct {
		Special *model.LocalizedString `json:"special" validate:"required"`

		Friends *model.LocalizedString `json:"friends" validate:"required"`
	}
)

//InitMediaRouter is initializing router method
func InitMediaRouter(group *echo.Group, service model.MediaService) (*MediaRouter, error) {
	mediaRouter := MediaRouter{
		mediaService: service,
	}
	router := group.Group("/games/:id")
	router.GET("/media", mediaRouter.get)
	router.PUT("/media", mediaRouter.put)

	return &mediaRouter, nil
}

// @Summary Change media for game
// @Description Change media data about game
// @Success 200 {object} "OK"
// @Failure 401 {object} "Unauthorized"
// @Failure 403 {object} "Forbidden"
// @Failure 404 {object} "Not found"
// @Failure 422 {object} "Unprocessable object"
// @Failure 500 {object} "Internal server error"
// @Router /api/v1/games/:id/media [put]
func (api *MediaRouter) put(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	mediaDto := new(Media)

	if err := ctx.Bind(mediaDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	
	if errs := ctx.Validate(mediaDto); errs != nil {
		return orm.NewServiceError(http.StatusUnprocessableEntity, errs)
	}

	media := model.Media{}
	input, err := maper.Map(mediaDto)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

 	err = mapstructure.Decode(input, &media)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := api.mediaService.Update(id, &media); err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "")
}

// @Summary Get media for game
// @Description Get media data about game
// @Success 200 {object} Media "OK"
// @Failure 401 {object} "Unauthorized"
// @Failure 403 {object} "Forbidden"
// @Failure 404 {object} "Not found"
// @Failure 500 {object} "Internal server error"
// @Router /api/v1/games/:id/media [get]
func (api *MediaRouter) get(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	
	media, err := api.mediaService.Get(id)

	if err != nil {
		return err
	}

	result := Media {}
	input, err := maper.Map(media)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = mapstructure.Decode(input, &result)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	
	return ctx.JSON(http.StatusOK, result)
}
package delivery

import (
	"fmt"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/model"
	"github.com/KonstantinPronin/advertising-backend/pkg/constants"
	"github.com/KonstantinPronin/advertising-backend/pkg/middleware"
	pkg "github.com/KonstantinPronin/advertising-backend/pkg/model"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"strconv"
)

type AdvertisementHandler struct {
	logger    *zap.Logger
	server    *echo.Echo
	sanitizer *bluemonday.Policy
	usecase   advertisement.Usecase
}

func NewAdvertisementHandler(
	server *echo.Echo,
	usecase advertisement.Usecase,
	sanitizer *bluemonday.Policy,
	logger *zap.Logger) {
	handler := AdvertisementHandler{
		logger:    logger,
		server:    server,
		sanitizer: sanitizer,
		usecase:   usecase,
	}

	server.GET("/adv/:id", handler.GetAd, middleware.ParseErrors)
	server.GET("/adv", handler.GetPage, middleware.ParseErrors)
	server.POST("/adv", handler.CreateAd, middleware.ParseErrors)
}

func (ah *AdvertisementHandler) GetAd(ctx echo.Context) error {
	id := ctx.Param("id")
	desc, imgs := false, false

	if params, ok := ctx.QueryParams()[constants.FieldsParam]; ok {
		for _, val := range params {
			if val == constants.IncludeDescription {
				desc = true
			} else if val == constants.IncludeImages {
				imgs = true
			}
		}
	}

	ad, err := ah.usecase.GetAd(id, desc, imgs)
	if err != nil {
		return err
	}

	if _, err := easyjson.MarshalToWriter(ad, ctx.Response().Writer); err != nil {
		ah.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (ah *AdvertisementHandler) GetPage(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam(constants.PageNumber))
	if err != nil || page < 1 {
		page = 1
	}

	desc, err := strconv.ParseBool(ctx.QueryParam(constants.Desc))
	if err != nil {
		desc = false
	}

	var handlerFunc func(page uint32, desc bool) (model.AdList, uint32, error)
	if order := ctx.QueryParam(constants.Order); order == "date" {
		handlerFunc = ah.usecase.GetAdsOrderByDate
	} else {
		handlerFunc = ah.usecase.GetAdsOrderByPrice
	}

	list, pages, err := handlerFunc(uint32(page), desc)
	if err != nil {
		return err
	}

	ah.addPageHeader(page, int(pages), ctx)

	if _, err := easyjson.MarshalToWriter(list, ctx.Response().Writer); err != nil {
		ah.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (ah *AdvertisementHandler) CreateAd(ctx echo.Context) error {
	ad := new(model.Ad)

	if err := easyjson.UnmarshalFromReader(ctx.Request().Body, ad); err != nil {
		ah.logger.Error(fmt.Sprintf("Request unmarshal error: %s", err.Error()))
		return pkg.NewInvalidArgument("wrong request body format")
	}

	if ad.Name == "" || ad.Price < 0 ||
		len([]rune(ad.Name)) > constants.MaxNameLength ||
		len([]rune(ad.Description)) > constants.MaxDescLength ||
		len(ad.Images) > constants.MaxImagesNumber ||
		len(ad.Images) < 1 {
		return pkg.NewInvalidArgument("wrong request body format")
	}

	ad.Id = ""
	ad.Created = ""
	ad.Name = ah.sanitizer.Sanitize(ad.Name)
	ad.Description = ah.sanitizer.Sanitize(ad.Description)
	for i, val := range ad.Images {
		ad.Images[i] = ah.sanitizer.Sanitize(val)
	}

	id, err := ah.usecase.CreateAd(ad)
	if err != nil {
		return nil
	}

	ad = &model.Ad{Id: id}
	if _, err := easyjson.MarshalToWriter(ad, ctx.Response().Writer); err != nil {
		ah.logger.Error(fmt.Sprintf("Response marshal error: %s", err.Error()))
		return err
	}

	return nil
}

func (ah *AdvertisementHandler) addPageHeader(page, totalPages int, ctx echo.Context) {
	ctx.Response().Header().Set(constants.TotalPages, strconv.Itoa(totalPages))
	ctx.Response().Header().Set(constants.PerPage, strconv.Itoa(constants.PageSize))
	ctx.Response().Header().Set(constants.Page, strconv.Itoa(page))

	if page+1 <= totalPages {
		ctx.Response().Header().Set(constants.NextPage, strconv.Itoa(page+1))
	}

	if page-1 >= 1 {
		ctx.Response().Header().Set(constants.PrevPage, strconv.Itoa(page-1))
	}
}

package advertisement

import "github.com/KonstantinPronin/advertising-backend/internal/advertisement/model"

type Usecase interface {
	CreateAd(ad *model.Ad) (string, error)
	GetAdsOrderByPrice(page uint32, desc bool) (model.AdList, uint32, error)
	GetAdsOrderByDate(page uint32, desc bool) (model.AdList, uint32, error)
	GetAd(id string, description bool, images bool) (*model.Ad, error)
}

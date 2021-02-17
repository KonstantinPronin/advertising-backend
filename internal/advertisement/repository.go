package advertisement

import "github.com/KonstantinPronin/advertising-backend/internal/advertisement/model"

type Repository interface {
	CreateAd(ad *model.Ad) (string, error)
	GetAdsOrderByPrice(page uint32, desc bool) (model.AdList, uint32, error)
	GetAdsOrderByDate(page uint32, desc bool) (model.AdList, uint32, error)
	GetAd(id string) (*model.Ad, error)
}

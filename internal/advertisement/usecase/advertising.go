package usecase

import (
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/model"
	"go.uber.org/zap"
)

type Advertising struct {
	repository advertisement.Repository
	logger     *zap.Logger
}

func (a *Advertising) CreateAd(ad *model.Ad) (string, error) {
	id, err := a.repository.CreateAd(ad)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *Advertising) GetAdsOrderByPrice(page uint32, desc bool) (model.AdList, uint32, error) {
	return a.repository.GetAdsOrderByPrice(page, desc)
}

func (a *Advertising) GetAdsOrderByDate(page uint32, desc bool) (model.AdList, uint32, error) {
	return a.repository.GetAdsOrderByDate(page, desc)
}

func (a *Advertising) GetAd(id string, description bool, images bool) (*model.Ad, error) {
	ad, err := a.repository.GetAd(id)
	if err != nil {
		return nil, err
	}

	if !description {
		ad.Description = ""
	}

	if !images && len(ad.Images) > 0 {
		ad.Images = ad.Images[:1]
	}

	return ad, nil
}

func NewAdvertising(repository advertisement.Repository, logger *zap.Logger) advertisement.Usecase {
	return &Advertising{
		repository: repository,
		logger:     logger,
	}
}

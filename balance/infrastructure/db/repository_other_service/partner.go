package repository_other_service

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/entity"
)

type Partner struct {
}

type PartnerInterface interface {
	AllPartnerActive() []entity.Partner
	AllPartnerCodeActive() []string
}

func (p Partner) AllPartnerActive() []entity.Partner {
	//todo call api from other service get info
	// fake value
	var rs []entity.Partner

	rs = append(rs,
		entity.Partner{
			ID:          1,
			PartnerCode: "TEST_1",
			Email:       "TEST_1@gmail.com",
			Status:      "active",
			Name:        "TEST_1",
			Address:     "anonymous",
			CreatedAt:   1,
			UpdatedAt:   2,
		},
		entity.Partner{
			ID:          2,
			PartnerCode: "TEST_2",
			Email:       "TEST_2@gmail.com",
			Status:      "active",
			Name:        "TEST_2",
			Address:     "anonymous",
			CreatedAt:   1,
			UpdatedAt:   2,
		},
		entity.Partner{
			ID:          3,
			PartnerCode: "TEST_3",
			Email:       "TEST_3@gmail.com",
			Status:      "active",
			Name:        "TEST_3",
			Address:     "anonymous",
			CreatedAt:   3,
			UpdatedAt:   3,
		},
		entity.Partner{
			ID:          3,
			PartnerCode: "TEST_1",
			Email:       "TEST_1@gmail.com",
			Status:      "active",
			Name:        "TEST_1",
			Address:     "anonymous",
			CreatedAt:   3,
			UpdatedAt:   3,
		},
		entity.Partner{
			ID:          4,
			PartnerCode: "TEST_4",
			Email:       "TEST_4@gmail.com",
			Status:      "active",
			Name:        "TEST_1",
			Address:     "anonymous",
			CreatedAt:   4,
			UpdatedAt:   4,
		},
		entity.Partner{
			ID:          5,
			PartnerCode: "TEST_5",
			Email:       "TEST_5@gmail.com",
			Status:      "active",
			Name:        "TEST_5",
			Address:     "anonymous",
			CreatedAt:   5,
			UpdatedAt:   5,
		},
	)

	return rs
}

func (p Partner) AllPartnerCodeActive() []string {
	return []string{"TEST_1", "TEST_2", "TEST_3", "TEST_4", "TEST_5"}
}

func NewPartnerRepository() PartnerInterface {
	return &Partner{}
}

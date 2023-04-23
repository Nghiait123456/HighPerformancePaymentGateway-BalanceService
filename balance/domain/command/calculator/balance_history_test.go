package calculator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalancePlaceHolderHistory_GetAllPlaceHolder(t *testing.T) {
	b := balancePlaceHolderHistory{}
	assert.Equal(t, []string{"test"}, b.allPartnerCode())
}

func TestBalancePlaceHolderHistory_LoadAllPlaceHolderAmountFrLogs(t *testing.T) {
	b := balancePlaceHolderHistory{}
	want := allAmountPlaceHolderFrLogs{}
	want["test"] = amountPlaceHolderFrLogs{
		partnerCode:       "test",
		amountPlaceHolder: 50000,
	}

	assert.Equal(t, want, b.loadAllPlaceHolderAmountFrLogs())
}

func TestBalancePlaceHolderHistory_InitBalancePlaceHolderHistory(t *testing.T) {
	b := balancePlaceHolderHistory{
		allAmountPlaceHolder: allAmountPlaceHolderFrLogs{},
	}
	want := allAmountPlaceHolderFrLogs{}
	want["test"] = amountPlaceHolderFrLogs{
		partnerCode:       "test",
		amountPlaceHolder: 50000,
	}

	b.initBalancePlaceHolderHistory()
	assert.Equal(t, want, b.allAmountPlaceHolder)
}

func TestBalancePlaceHolderHistory_TestGetAllPlaceHolder(t *testing.T) {
	b := balancePlaceHolderHistory{
		allAmountPlaceHolder: allAmountPlaceHolderFrLogs{},
	}
	want := allAmountPlaceHolderFrLogs{}
	want["test"] = amountPlaceHolderFrLogs{
		partnerCode:       "test",
		amountPlaceHolder: 50000,
	}
	b.initBalancePlaceHolderHistory()

	assert.Equal(t, b.allAmountPlaceHolder, b.GetAllPlaceHolder())
}

func TestNewBalancePlaceHolderHistory(t *testing.T) {
	b := NewBalancePlaceHolderHistory()
	assert.NotNil(t, b)
}

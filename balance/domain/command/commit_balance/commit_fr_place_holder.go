package commit_balance

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
)

type (
	// commit balance place holder to balance, update status balance log
	CommitPlaceHolderToBalanceDB struct {
		PartnerCode            string
		IndexLogRequestLatest  uint64
		Balance                uint64
		TotalAmountPlaceHolder uint64
		CnBalance              sql.Connect
	}

	CommitPlaceHolderToBalanceDBInterface interface {
		Commit() error
	}
)

func (c *CommitPlaceHolderToBalanceDB) Commit() error {
	rpBalance := repository.NewBalanceRepository()
	rpBalance.SetConnect(c.CnBalance)

	commit := repository.CommitAmountPlaceHolder{
		PartnerCode:           c.PartnerCode,
		AmountPlaceHolder:     c.TotalAmountPlaceHolder,
		IndexLogRequestLatest: c.IndexLogRequestLatest,
	}
	err := rpBalance.CommitAmountPlaceHolderToBalance(commit)
	if err != nil {
		return err
	}

	return nil
}

func NewCommitAmountPlaceHolderToBalance() CommitPlaceHolderToBalanceDBInterface {
	return &CommitPlaceHolderToBalanceDB{}
}

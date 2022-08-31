package env

type (
	DBConfigSqlDefault struct {
		UserName string
		PassWord string
		Link     string
		NameDB   string
	}

	DBConfigSqlDefaultInterface interface {
		BaseConfigInterface
		GetDSN() (string, error)
	}
)

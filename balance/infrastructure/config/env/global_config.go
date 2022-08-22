package env

type (
	GlobalConfig struct {
		DefaultDBSqlUserName string `env:"DEFAULT_DB_SQL_USERNAME,notEmpty"`
		DefaultDBSqlPassWord string `env:"DEFAULT_DB_SQL_PASSWORD,notEmpty"`
		DefaultDBSqlLink     string `env:"DEFAULT_DB_SQL_LINK,notEmpty"`
		DefaultDBSqlNameDB   string `env:"DEFAULT_DB_SQL_NAME_DB,notEmpty"`
	}

	DBConfigSqlDefault struct {
		UserName string
		PassWord string
		Link     string
		NameDB   string
	}

	DBConfigSqlDefaultInterface interface {
		GetDSN() (string, error)
	}

	GlobalConfigInterface interface {
		LoadConfig() error
		MapConfigToStruct() error
	}
)

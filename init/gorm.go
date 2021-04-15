package init

import (
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitGorm(sqlDialect string, sqlDsn string, opts ...gorm.Option) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch strings.ToLower(sqlDialect) {
	case "mysql":
		dialector = mysql.Open(sqlDsn)
	case "sqlite":
		dialector = sqlite.Open(sqlDsn)
	default:
		err := errors.New("invalid sql dialect")
		return nil, errors.WithStack(err)
	}

	db, err := gorm.Open(dialector, opts...)
	return db, errors.WithStack(err)
}

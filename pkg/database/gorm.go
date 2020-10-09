package database

import (
	"strings"

	"github.com/VulpesFerrilata/library/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGorm(config *config.Config) (*gorm.DB, error) {
	driverName := strings.ToLower(config.SqlSettings.DriverName)
	var dialector gorm.Dialector
	switch driverName {
	case "mysql":
		dialector = mysql.Open(config.SqlSettings.DataSource)
	case "sqlite":
		dialector = sqlite.Open(config.SqlSettings.DataSource)
	}
	return gorm.Open(dialector, &gorm.Config{})
}

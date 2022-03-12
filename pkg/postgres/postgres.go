package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresAuth struct {
	Host       string
	User       string
	Password   string
	Dbname     string
	Port       string
	Sslmode    string
	TimeZone   string
	GormOption gorm.Config
}

// New -.
func New(auth *PostgresAuth) (*gorm.DB, error) {
	dsn := "host=" + auth.Host + " user=" + auth.User +
		" password=" + auth.Password + " dbname=" + auth.Dbname +
		" port=" + auth.Port + " sslmode=" + auth.Sslmode + " TimeZone=" + auth.TimeZone
	return gorm.Open(postgres.Open(dsn), &auth.GormOption)
}

// Package postgres implements postgres connection.
package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOption struct {
	hotst string
	user string
	password string
	dbname string
	port string
	sslmode string
	timeZone string
}
// New -.
func New(PostgresOption) (*gorm.DB, error) {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

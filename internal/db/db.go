package db

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DB struct {
	client *gorm.DB
}

func Connect(ctx context.Context, dsn string) (*DB, error) {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := client.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{client: client}, nil
}

func (d *DB) Client() *gorm.DB {
	return d.client
}

func (d *DB) Close(ctx context.Context) error {
	sqlDB, err := d.client.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

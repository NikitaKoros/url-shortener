package shortening

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/model"
	"gorm.io/gorm"
)

type gormDB struct {
	db *gorm.DB
}

func NewGormDB(db *gorm.DB) *gormDB {
	return &gormDB{db: db}
}

func (g *gormDB) Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error) {
	const op = "shortening.gorm.Put"

	shortening.CreatedAt = time.Now().UTC()
	
	var existing model.Shortening
	if err := g.db.WithContext(ctx).Where("identifier = ?", shortening.Identifier).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("%s: %w", op, model.ErrIdentifierExists)
	}

	if err := g.db.WithContext(ctx).Create(&shortening).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortening, nil
}

func (g *gormDB) Get(ctx context.Context, shorteningID string) (*model.Shortening, error) {
	const op = "shortening.gorm.Get"

	var shortening model.Shortening
	if err := g.db.WithContext(ctx).Where("identifier = ?", shorteningID).First(&shortening).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortening, nil
}

func (g *gormDB) IncrementVisits(ctx context.Context, shorteningID string) error {
	const op = "shortening.gorm.IncrementVisits"

	if err := g.db.WithContext(ctx).
	Model(&model.Shortening{}).
	Where("identifier = ?", shorteningID).
	Updates(map[string]interface{}{
		"visits": gorm.Expr("visits + ?", 1),
		"updated_at": time.Now().UTC(),
	}).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}





package container

import "gorm.io/gorm"

type Container interface {
	GetDatabase() *gorm.DB
}

type container struct {
	db *gorm.DB
}

func NewContainer(db *gorm.DB) Container {
	return &container{
		db: db,
	}
}

func (c *container) GetDatabase() *gorm.DB {
	return c.db
}

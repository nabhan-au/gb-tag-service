package repository

import (
	"github.com/GarnBarn/common-go/model"
	"gorm.io/gorm"
)

type Tag interface {
	GetAllTag(author string) (tags []model.Tag, err error)
	Create(tag *model.Tag) error
	Update(tag *model.Tag) error
	GetByID(id int) (*model.Tag, error)
	DeleteTag(tagID int) error
}

type tag struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) Tag {
	// Migrate the db
	db.AutoMigrate(&model.Tag{})

	return &tag{
		db: db,
	}
}

func (t *tag) GetAllTag(author string) (tags []model.Tag, err error) {
	res := t.db.Model(&tags).Where("author = ?", author).Find(&tags)
	if res.Error != nil {
		return tags, res.Error
	}
	return tags, nil
}

func (t *tag) GetByID(id int) (*model.Tag, error) {
	tag := model.Tag{}
	result := t.db.First(&tag, id)
	return &tag, result.Error
}

func (t *tag) Create(tag *model.Tag) error {
	result := t.db.Create(tag)
	return result.Error
}

func (t *tag) Update(tag *model.Tag) error {
	result := t.db.Save(tag)
	return result.Error
}
func (t *tag) DeleteTag(tagID int) error {
	result := t.db.Delete(&model.Tag{}, tagID)
	return result.Error
}

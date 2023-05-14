package models

import "time"

type Article struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"type:varchar(255);not null;required" json:"title"`
	CategoryID uint      `gorm:"not null; required" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Content    string    `gorm:"type:text;not null;required" json:"content"`
	Thumbnail  string    `gorm:"type:text;not null;required" json:"thumbnail"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ArticleResponse struct {
	ID         uint                    `json:"id"`
	Title      string                  `json:"title"`
	CategoryID uint                    `json:"-"`
	Category   ArticleCategoryResponse `gorm:"foreignKey:CategoryID" json:"category"`
	Content    string                  `json:"content"`
	Thumbnail  string                  `json:"thumbnail"`
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`
}

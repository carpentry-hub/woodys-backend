package models

import (
	"time"
	
	"github.com/lib/pq"
)

type Project struct {
	ID        	  int8 			`json:"id"`
	CreatedAt     time.Time		`json:"created_at"`
	Owner         int			`json:"owner"` // id del usuario que creo que proyecto
	Title         string		`json:"title"`
	UpdatedAt     time.Time		`json:"updated_at"`
	AverageRating float32		`json:"avarage_rating"`
	RatingCount   int			`json:"rating_count"` // cantidad de ratings del proyecto
	Materials     pq.StringArray `json:"materials" gorm:"type:varchar[]"`
	Tools         pq.StringArray `json:"tools" gorm:"type:varchar[]"`
	Description   string		`json:"description"`
	Style         pq.StringArray `json:"style" gorm:"type:varchar[]"`
	Enviroment    pq.StringArray `json:"enviroment" gorm:"type:text[]"`
	Portrait      string		`json:"portrait"`
	Images        pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Tutorial      string		`json:"tutorial"`
	TimeToBuild   int			`json:"time_to_build"`
}


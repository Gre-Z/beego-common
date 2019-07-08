package controllers

import "github.com/Gre-Z/common/jtime"

type BaseMode struct {
	ID        uint            `gorm:"primary_key" json:"id"`
	CreatedAt jtime.JsonTime  `json:"created_at"`
	UpdatedAt jtime.JsonTime  `json:"updated_at"`
	DeletedAt *jtime.JsonTime `sql:"index" json:"deleted_at"`
}

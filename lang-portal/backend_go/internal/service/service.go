package service

import (
	"database/sql"
)

// Services holds all service instances
type Services struct {
	Word  *WordService
	Group *GroupService
	Study *StudyService
}

// NewServices creates all services
func NewServices(db *sql.DB) *Services {
	return &Services{
		Word:  NewWordService(db),
		Group: NewGroupService(db),
		Study: NewStudyService(db),
	}
}

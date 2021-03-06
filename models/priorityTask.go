package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PriorityTask struct {
	gorm.Model

	Title       string
	Description string
	Complete    bool

	PersonID  int
	ProjectID int

	PriorityLevel PriorityLevelType
}

type PriorityLevelType uint

const (
	Undefined PriorityLevelType = iota
	ASAP
	SOON
	UPCOMING
	LATER
	end
)

// end is used to ensure type PriorityLevelType is restricted to range

// write unit test for this
func IsValidPriorityLevel(level PriorityLevelType) bool {
	return uint(level) < uint(end)
}

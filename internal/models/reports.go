package models

import "gorm.io/gorm"

type Report struct {
	gorm.Model
	ReporterID  uint
	ContentType ContentType
	ContentID   uint
	Reason      string `gorm:"type:text"`
	Status      ReportStatus
}

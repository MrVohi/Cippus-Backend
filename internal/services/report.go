package services

import (
	"cippus-backend/internal/models"

	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

func (s *ReportService) CreateReport(reporterID uint, contentType models.ContentType, contentID uint, reason string) (models.Report, error) {
	report := models.Report{
		ReporterID:  reporterID,
		ContentType: contentType,
		ContentID:   contentID,
		Reason:      reason,
		Status:      models.ReportPending,
	}
	result := s.db.Create(&report)
	return report, result.Error
}

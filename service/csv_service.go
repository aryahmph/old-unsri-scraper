package service

import (
	"old-unsri-scraper/entity"
)

type CSVService interface {
	WriteAllToCSV(students []entity.Student)
}

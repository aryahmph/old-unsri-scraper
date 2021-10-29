package service

import (
	"encoding/csv"
	"old-unsri-scraper/entity"
	"old-unsri-scraper/helper"
)

type CSVServiceImpl struct {
	Writer *csv.Writer
}

func NewCSVServiceImpl(writer *csv.Writer) *CSVServiceImpl {
	return &CSVServiceImpl{Writer: writer}
}

func (service *CSVServiceImpl) WriteAllToCSV(students []entity.Student) {
	var data [][]string
	for _, student := range students {
		data = append(data, []string{student.Name, student.NIM})
	}

	err := service.Writer.WriteAll(data)
	helper.LogIfError(err)
}

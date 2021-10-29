package service

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"old-unsri-scraper/entity"
	"old-unsri-scraper/helper"
	"sync"
)

type ScraperServiceImpl struct {
	csvService CSVService
}

func NewScraperServiceImpl(csvService CSVService) *ScraperServiceImpl {
	return &ScraperServiceImpl{csvService: csvService}
}

func (service *ScraperServiceImpl) Do(url string) {
	htmlResponse, err := http.Get(url)
	helper.LogIfError(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		helper.LogIfError(err)
	}(htmlResponse.Body)

	if htmlResponse.StatusCode != 200 {
		helper.LogIfError(errors.New("err : unexpected status code"))
	}

	document, err := goquery.NewDocumentFromReader(htmlResponse.Body)
	helper.LogIfError(err)

	mutex := sync.Mutex{}
	wg := new(sync.WaitGroup)
	var students []entity.Student

	document.Find(".mainContent-news-element table tbody").Children().Next().
		Children().Find("table tbody").Children().Each(func(i int, selection *goquery.Selection) {
		go func() {
			wg.Add(1)
			defer wg.Done()

			student := entity.Student{}
			if i != 0 && i != 1 {
				selection.Children().Each(func(i int, selection *goquery.Selection) {
					switch i {
					case 2:
						student.Name = helper.EscapeNewLineAndIndent(selection.Text())
						break
					case 3:
						student.NIM = helper.EscapeNewLineAndIndent(selection.Text())
						break
					}
				})
			}
			if student.Name != "" && student.NIM != "" {
				mutex.Lock()
				students = append(students, student)
				mutex.Unlock()
			}
		}()
	})

	wg.Wait()

	go func() {
		wg.Add(1)
		service.csvService.WriteAllToCSV(students)
		wg.Done()
	}()

	wg.Wait()
	log.Println("CSV OK")
}

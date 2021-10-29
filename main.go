package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"old-unsri-scraper/helper"
	"old-unsri-scraper/service"
	"os"
	"sync"
	"time"
)

func main() {
	log.Println("Running app...")
	start := time.Now()

	majorCodes := map[string]string{
		"Teknik Informatika REG": "9-10001-1",
		"Teknik Informatika BIL": "9-10002-2",
		"Sistem Komputer REG":    "9-10003-3",
	}
	wg := new(sync.WaitGroup)

	for key, value := range majorCodes {
		log.Println("Scraping", key)
		batchMin := 2018
		batchMax := 2021

		for ; batchMin <= batchMax; batchMin++ {
			go func(batchMin int, key string, value string) {
				wg.Add(1)
				defer wg.Done()

				file, err := os.Create(fmt.Sprintf("%s-%d.csv", key, batchMin))
				helper.LogIfError(err)

				csvWriter := csv.NewWriter(file)
				csvService := service.NewCSVServiceImpl(csvWriter)
				scraperService := service.NewScraperServiceImpl(csvService)

				scraperService.Do(fmt.Sprintf("http://old.unsri.ac.id/?act=daftar_mahasiswa&fak_prodi=%s&angkatan=%d",
					value, batchMin))
			}(batchMin, key, value)
		}
	}

	wg.Wait()
	log.Println("OK")
	log.Println("Execution time :", time.Since(start))
}

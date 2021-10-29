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
	waitGroup := new(sync.WaitGroup)

	for key, value := range majorCodes {
		batchMin := 2018
		batchMax := 2021

		for ; batchMin <= batchMax; batchMin++ {
			go func(batchMin int, key string, value string) {
				waitGroup.Add(1)
				defer waitGroup.Done()

				file, err := os.Create(fmt.Sprintf("%s-%d.csv", key, batchMin))
				helper.LogIfError(err)

				csvWriter := csv.NewWriter(file)
				csvService := service.NewCSVServiceImpl(csvWriter)
				scraperService := service.NewScraperServiceImpl(csvService)

				scraperService.Do(fmt.Sprintf("http://old.unsri.ac.id/?act=daftar_mahasiswa&fak_prodi=%s&angkatan=%d",
					value, batchMin), waitGroup)
			}(batchMin, key, value)
		}
	}

	waitGroup.Wait()
	// Sleep while because last CSV need extra time to write, IDK why
	time.Sleep(5 * time.Second)

	log.Println("OK")
	log.Println("Execution time :", time.Since(start))
}

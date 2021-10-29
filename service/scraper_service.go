package service

import "sync"

type ScraperService interface {
	Do(url string, waitGroup *sync.WaitGroup)
}

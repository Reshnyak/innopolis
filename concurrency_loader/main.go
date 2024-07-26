package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"sync"
	"time"
)

type Site struct {
	url string
}

const (
	maxWorkers    int           = 3
	timeoutRes    time.Duration = 100000 * time.Millisecond //ожидание результата
	timeoutClient time.Duration = 200 * time.Millisecond    //ожидание ответа на запрос от http.Client
)

var once sync.Once

func main() {
	urls := []string{
		"https://examle.com/file1.jpg",
		"https://examle.com/file2.jpg",
		"https://examle.com/file3.jpg",
		"https://examle.com/file4.jpg",
		"https://examle.com/file5.jpg",
		"https://examle.com/file6.jpg",
		"https://examle.com/file7.jpg",
		"https://examle.com/file8.jpg",
	}

	resultsSlice := downloadFiles(urls)
	checkErr := false
	for _, res := range resultsSlice {
		if res.err != nil {
			fmt.Println(res.err)
			once.Do(func() {
				checkErr = true
			})
			continue
		}
		fmt.Println(res)
	}
	if !checkErr {
		fmt.Println("All files downloaded successfully")
	}
}

// downloadFiles - принимает набор url для скачивания файлов.
// При успехе возвращает срез результатов: URL, СтатусКод, Имя файла; - иначе ошибку
func downloadFiles(urls []string) []Result {
	wg := new(sync.WaitGroup)
	jobs := make(chan Site)
	results := make(chan Result)
	workerResults := make([]chan Result, 0, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workerResults = append(workerResults, results)
		go worker(wg, i+1, jobs, workerResults[i])
	}

	wg.Add(len(urls))
	go func() {
		for _, url := range urls {
			jobs <- Site{url: url}
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	downloadRes := WaitResults(fanIn(workerResults))
	return downloadRes

}

// worker выполняет работу =) при завершении работ закрывает канал результатов,
// используя sync.Once, чтобы избежать закрытия закрытого канала другим воркером
func worker(wg *sync.WaitGroup, id int, jobs <-chan Site, results chan<- Result) {
	for job := range jobs {
		log.Printf("Worker %d downloading file from %s\n", id, job.url)
		results <- *process(&job)
		wg.Done()
	}

}

// Делаем запрос, получаем ответ и копируем в созданный файл
// Возвращаем результат. В случая возникновения ошибки кладем её в результат
func process(site *Site) *Result {
	client := http.Client{}

	resp, err := client.Get(site.url)
	if err != nil {
		return NewResult().AddUrl(site.url).AddError(fmt.Errorf("process client.Get:%s", err))
	}
	defer resp.Body.Close()

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return NewResult().AddUrl(site.url).AddError(fmt.Errorf("ParseMediaType client.Get:%s", err))
	}
	fileName := params["filename"]

	file, err := os.Create(fileName)
	if err != nil {
		return NewResult().AddUrl(site.url).AddError(fmt.Errorf("process os.Create:%s", err))
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return NewResult().AddUrl(site.url).AddError(fmt.Errorf("process io.Copy:%s", err))
	}
	err = file.Close()
	if err != nil {
		NewResult().AddUrl(site.url).AddError(fmt.Errorf("process file.Close:%s", err))
	}
	return NewResult().AddUrl(site.url).AddFileName(fileName).AddStatus(resp.StatusCode)
}

// fanIn - объединяет результаты из каналов воркеров в один канал
func fanIn(results []chan Result) <-chan Result {
	var wg sync.WaitGroup
	out := make(chan Result)
	output := func(c chan Result) {
		for res := range c {
			out <- res
		}
		wg.Done()
	}
	wg.Add(len(results))
	for _, res := range results {
		go output(res)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Ждем результаты - если долго нет, то выходим с тем что есть
func WaitResults(results <-chan Result) []Result {
	var resultSlice []Result
	for {
		select {
		case res, ok := <-results:
			if ok {
				resultSlice = append(resultSlice, res)
				continue
			}
			return resultSlice
		case <-time.After(timeoutRes):
			return append(resultSlice, *NewResult().AddError(fmt.Errorf("timeout  reached, aborting remaining downloads")))
		}
	}
}

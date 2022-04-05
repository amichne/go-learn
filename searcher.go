package main

import (
	"os"
	"sync"
)

var waitGroup = sync.WaitGroup{}
var channel = make(chan SearchResult)
var initialDirectory, _ = os.Getwd()

func main() {
	ProcessFiles(initialDirectory, &SearchQuery{
		Query: "ioutil",
		Regex: false,
	})
	go func() {
		waitGroup.Wait()
		close(channel)
	}()
	for elem := range channel {
		println("Emitted on channel for file: " + elem.Path)
	}
}

type SearchQuery struct {
	Query string
	Regex bool
}

type SearchResult struct {
	Path       string
	MatchFound bool
}

func ProcessFiles(path string, query *SearchQuery) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(files); i++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			if files[index].IsDir() {
				ProcessFiles(path+"/"+files[index].Name(), query)
			} else {
				SearchFile(query, path+"/"+files[index].Name())
			}
		}(i)
	}
}

func SearchFile(query *SearchQuery, path string) {
	os.ReadFile(path)
	channel <- SearchResult{
		Path:       path,
		MatchFound: false,
	}
}

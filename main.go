package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Name string `yaml:"name"`
	From string `yaml:"from"`
}

type Input struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
	Data []Data `yaml:"data"`
}

type BrickConfYaml struct {
	Version string  `yaml:"version"`
	Module  string  `yaml:"module"`
	Input   []Input `yaml:"input"`
}

const FILE_NUM = 9999
const FILE_SUFFIX = "file.yml"
const SAMPLE_FOLDER = "samples"

func Naive() []*BrickConfYaml {
	files := make([]string, FILE_NUM)
	toParse := make([]*[]byte, FILE_NUM)
	parsed := make([]*BrickConfYaml, FILE_NUM)

	for idx := range files {
		path := filepath.Join(SAMPLE_FOLDER, fmt.Sprintf("%v_%s", idx, FILE_SUFFIX))
		files[idx] = path
	}

	for idx, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}

		toParse[idx] = &b
	}

	for idx, b := range toParse {
		done := new(BrickConfYaml)

		err := yaml.Unmarshal(*b, done)
		if err != nil {
			log.Fatal(err)
		}
		parsed[idx] = done
	}

	return parsed
}

func WithGoKeyword() []*BrickConfYaml {
	var wgp sync.WaitGroup

	files := make([]string, FILE_NUM)
	toParse := make([]*[]byte, FILE_NUM)
	parsed := make([]*BrickConfYaml, FILE_NUM)

	for idx := range files {
		path := filepath.Join(SAMPLE_FOLDER, fmt.Sprintf("%v_%s", idx, FILE_SUFFIX))
		files[idx] = path
	}

	for idx, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}

		toParse[idx] = &b
	}

	for idx, b := range toParse {
		wgp.Add(1)

		go func(idx int, b *[]byte) {
			defer wgp.Done()

			done := new(BrickConfYaml)

			err := yaml.Unmarshal(*b, done)
			if err != nil {
				log.Fatal(err)
			}

			parsed[idx] = done
		}(idx, b)
	}

	wgp.Wait()

	return parsed
}

func WithGoroutines() []*BrickConfYaml {
	var wgf sync.WaitGroup
	var wgp sync.WaitGroup
	files := make([]string, FILE_NUM)
	toParse := make([]*[]byte, FILE_NUM)
	parsed := make([]*BrickConfYaml, FILE_NUM)

	for idx := range files {
		path := filepath.Join(SAMPLE_FOLDER, fmt.Sprintf("%v_%s", idx, FILE_SUFFIX))
		files[idx] = path
	}

	for idx, f := range files {
		wgf.Add(1)

		go func(idx int, f string) {
			defer wgf.Done()

			b, err := os.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}

			toParse[idx] = &b
		}(idx, f)
	}

	wgf.Wait()

	for idx, b := range toParse {
		wgp.Add(1)

		go func(idx int, b *[]byte) {
			defer wgp.Done()
			done := new(BrickConfYaml)

			err := yaml.Unmarshal(*b, done)
			if err != nil {
				log.Fatal(err)
			}
			parsed[idx] = done
		}(idx, b)
	}

	wgp.Wait()

	return parsed
}

type PathWithIndex struct {
	Path string
	Idx  int
}

type FileWithIndex struct {
	File *[]byte
	Idx  int
}

type ConfWithIndex struct {
	Conf *BrickConfYaml
	Idx  int
}

func WithChannels() []*BrickConfYaml {
	files := make([]string, FILE_NUM)
	confs := make([]*BrickConfYaml, FILE_NUM)

	toRead := make(chan *PathWithIndex, FILE_NUM)
	toParse := make(chan *FileWithIndex, FILE_NUM)
	parsed := make(chan *ConfWithIndex, FILE_NUM)

	for idx := range files {
		toRead <- &PathWithIndex{
			Path: filepath.Join(SAMPLE_FOLDER, fmt.Sprintf("%v_%s", idx, FILE_SUFFIX)),
			Idx:  idx,
		}
	}

	for n := FILE_NUM; n > 0; {
		select {
		case pathWithIndex := <-toRead:
			go func() {
				b, err := os.ReadFile(pathWithIndex.Path)
				if err != nil {
					log.Fatal(err)
				}

				toParse <- &FileWithIndex{
					File: &b,
					Idx:  pathWithIndex.Idx,
				}
			}()

		case fileWithIndex := <-toParse:
			go func() {
				done := new(BrickConfYaml)

				err := yaml.Unmarshal(*fileWithIndex.File, done)
				if err != nil {
					log.Fatal(err)
				}

				parsed <- &ConfWithIndex{
					Conf: done,
					Idx:  fileWithIndex.Idx,
				}
			}()
		case done := <-parsed:
			confs[done.Idx] = done.Conf
			n--
		}
	}

	return confs
}

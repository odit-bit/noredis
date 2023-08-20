package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func ReadConfig(path string) map[string]string {
	dir, err := os.Getwd()
	fmt.Println(dir, err)
	f, err := os.Open("./.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config := readConfigFile(f)
	for key, value := range config {
		err := os.Setenv(key, value)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}

func readConfigFile(rd io.Reader) map[string]string {
	bucket := map[string]string{}

	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.Split(scanner.Text(), ":")
		if len(text) != 2 {
			log.Fatal("conf file syntax error")
		}
		key := strings.ToUpper(strings.TrimSpace(text[0]))
		value := strings.TrimSpace(text[1])
		bucket[key] = value
	}

	return bucket
}

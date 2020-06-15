package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig() LociConfig {
	fname := ".loci.yml"
	if !fileExists(fname) {
		log.Fatal(".loci.yml not found in current folder.")
	}

	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config LociConfig
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

	return config
}

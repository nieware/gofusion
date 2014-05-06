package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Global Settings for the program
type GlobalSettings struct {
	Username string // username
	HiScore  uint32 // hiscore for user

	fileName string
}

// Constructor
func NewGlobalSettings(fileName string) *GlobalSettings {
	g := new(GlobalSettings)
	g.fileName = fileName
	return g
}

func (g *GlobalSettings) GetHiScore() uint32 {
	g.readFromFile()
	return g.HiScore
}

func (g *GlobalSettings) SetHiScore(v uint32) {
	g.HiScore = v
	g.writeToFile()
}

// get name of settings file
func (g *GlobalSettings) getFileName() string {
	return g.fileName
}

// read settings from file
func (g *GlobalSettings) readFromFile() {
	//fmt.Println(g.fileName)
	n, err := ioutil.ReadFile(g.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			g.writeToFile()
			return
		} else {
			panic(err)
		}
	}
	p := bytes.NewBuffer(n)
	dec := json.NewDecoder(p)
	err = dec.Decode(&g)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("read settings from file: %v\n", g)
}

// write settings to file
func (g *GlobalSettings) writeToFile() {
	fmt.Println("writeToFile")
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(g)
	err := ioutil.WriteFile(g.fileName, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

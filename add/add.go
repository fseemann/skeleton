package add

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile("\\${(.*?)}")

type projectDescriptor struct {
	Version     string      `json:"version"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Variables   []string    `json:"variables"`
	Structure   []structure `json:"structure"`
}

type structure struct {
	Dir      string `json:"dir"`
	File     string `json:"file"`
	Template string `json:"template"`
}

func Add() {
	fileContent, err := ioutil.ReadFile("src/github.com/manic/skeleton/templates/add/maven-domain-module.json")
	if err != nil {
		log.Fatal(err)
	}

	var pd projectDescriptor
	if err := json.Unmarshal(fileContent, &pd); err != nil {
		log.Fatal(err)
	}

	variables := readVariables(pd.Variables)
	for _, struc := range pd.Structure {
		actualDir := replaceVariables(struc.Dir, variables)
		t := template.Must(template.ParseFiles("src/github.com/manic/skeleton/templates/add/" + struc.Template))
		if err := os.MkdirAll(actualDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(actualDir + "/" + struc.File)
		if err != nil {
			log.Fatal(err)
		}

		writer := bufio.NewWriter(file)
		if err := t.Execute(writer, variables); err != nil {
			log.Fatal(err)
		}
		_ = writer.Flush()
		_ = file.Close()
	}
}

func replaceVariables(value string, variables map[string]string) string {
	submatch := regex.FindAllStringSubmatch(value, -1)
	if submatch == nil {
		return value
	} else {
		for _, v := range submatch {
			value = strings.Replace(value, v[0], variables[v[1]], 1)
		}
		return value
	}
}

func readVariables(variables []string) map[string]string {
	readVariables := make(map[string]string, len(variables))
	for {
		for _, v := range variables {
			fmt.Printf("Type %v: ", v)
			readVariables[v] = readLine()
		}

		fmt.Print("Are values correct?[y/n]: ")
		if line := readLine(); line == "y" {
			break
		}
	}

	return readVariables
}

func readLine() string {
	var readLine string
	if _, err := fmt.Scan(&readLine); err != nil {
		panic(err)
	}
	return readLine
}

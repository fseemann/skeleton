package add

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

type response struct {
	FileUrls []string `json:"fileUrls"`
}

func Add(args []string) {
	response := getTemplate([]byte(` { "groupId": "fseemann", "name": "maven-domain-module", "version": "latest" } `))
	files := createTmpFiles(response.FileUrls)

	filecontents, err := ioutil.ReadFile(files[0].Name())
	if err != nil {
		log.Fatal(err)
	}

	var pd projectDescriptor
	if err := json.Unmarshal(filecontents, &pd); err != nil {
		log.Fatal(err)
	}

	templateFiles := files[1:]
	variables := readVariables(pd.Variables)
	for i, struc := range pd.Structure {
		actualDir := replaceVariables(struc.Dir, variables)
		t := template.Must(template.ParseFiles(templateFiles[i].Name()))
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

	removeFiles(files)
}

func createTmpFiles(fileUrls []string) []*os.File {
	files := make([]*os.File, len(fileUrls))
	for i, url := range fileUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
		if err != nil {
			removeFiles(files)
			log.Fatal("Cannot create temporary file", err)
		}

		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			removeFiles(files)
			log.Fatal(err)
		}

		resp.Body.Close()
		files[i] = tmpFile
	}

	return files
}

func removeFiles(files []*os.File) {
	for _, f := range files {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}
}

func getTemplate(body []byte) response {
	resp, err := http.Post("https://p33dswbrne.execute-api.eu-central-1.amazonaws.com/develop/template", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseMarshaled response
	if err := json.Unmarshal(respBody, &responseMarshaled); err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return responseMarshaled
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

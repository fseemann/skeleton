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

const API_URL = "https://p33dswbrne.execute-api.eu-central-1.amazonaws.com/develop/template"
const SKELETON_FILE = "Skeletonfile.json"

var regex = regexp.MustCompile("\\${(.*?)}")

type skeletonfile struct {
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

type request struct {
	GroupId string `json:"groupId"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Add(args []string) {
	requestBody := createRequestBody(args)

	response := fetchTemplate(requestBody)
	templateFiles := createTmpFiles(response.FileUrls)
	defer removeFiles(templateFiles)

	skeletonfile := parseSkeletonfile(templateFiles[SKELETON_FILE].Name())
	variables := readVariables(skeletonfile.Variables)
	execute(skeletonfile, variables, templateFiles)
}

func execute(sf skeletonfile, variables map[string]string, templateFiles map[string]*os.File) {
	for _, struc := range sf.Structure {
		actualDir := replaceVariables(struc.Dir, variables)
		t := template.Must(template.ParseFiles(templateFiles[struc.Template].Name()))
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

func parseSkeletonfile(fileName string) skeletonfile {
	skeletonfileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var pd skeletonfile
	if err := json.Unmarshal(skeletonfileContents, &pd); err != nil {
		log.Fatal(err)
	}
	return pd
}

func createRequestBody(args []string) []byte {
	if len(args) != 1 {
		log.Fatalf("Provide a template identifier")
	}
	arg := args[0]
	splitVersion := strings.Split(arg, ":")

	version := "latest"
	if len(splitVersion) == 2 {
		version = splitVersion[1]
	}

	groupIdAndName := strings.Split(splitVersion[0], "/")
	if len(groupIdAndName) != 2 {
		log.Fatalf("Invalid template name.")
	}

	marshal, err := json.Marshal(request{
		GroupId: groupIdAndName[0],
		Name:    groupIdAndName[1],
		Version: version,
	})
	if err != nil {
		log.Fatalf("Could not create request body.", err)
	}

	return marshal
}

func createTmpFiles(fileUrls []string) map[string]*os.File {
	files := make(map[string]*os.File, len(fileUrls))
	for _, url := range fileUrls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		tmpFile, err := ioutil.TempFile(os.TempDir(), "skeleton-template-")
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

		lastIndex := strings.LastIndex(url, "/")
		filename := url[lastIndex+1:]
		files[filename] = tmpFile
	}

	return files
}

func removeFiles(files map[string]*os.File) {
	for _, f := range files {
		removeFile(f)
	}
}

func removeFile(f *os.File) {
	if f == nil {
		return
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
}

func fetchTemplate(body []byte) response {
	resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(body))
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

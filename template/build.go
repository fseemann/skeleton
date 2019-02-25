package template

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const zipName = "skeleton-image.zip"

type zipper struct {
	destFile *os.File
	srcDir   string
	writer   *zip.Writer
}

func Build(args []string) {
	zipFile, e := os.Create(zipName)
	if e != nil {
		log.Fatalf("Could not build zip zipFile.", e)
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	dir, e := os.Getwd()
	z := zipper{
		destFile: zipFile,
		srcDir:   dir,
		writer:   writer,
	}
	if e = z.zipDir(); e != nil {
		log.Fatalf("Could not zip working directory.", e)
	}
}

func (z *zipper) zipDir() error {
	err := filepath.Walk(z.srcDir, z.zipFile)
	return err
}

func (z *zipper) zipFile(path string, info os.FileInfo, err error) error {
	if !info.Mode().IsRegular() || info.Size() == 0 || info.Name() == z.destFile.Name() {
		return nil
	}

	file, e := os.Open(path)
	if e != nil {
		return e
	}
	defer file.Close()

	filename := strings.TrimPrefix(path, z.srcDir+"\\")
	writer, err := z.writer.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file.", err)
	}

	_, e = io.Copy(writer, file)
	return e
}

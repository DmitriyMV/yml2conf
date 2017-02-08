package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go4.org/sort"
	"gopkg.in/yaml.v2"
)

type FileStruct struct {
	Package string   `yaml:"package"`
	Imports []Import `yaml:"import"`
}

type Import struct {
	Package string
	Version string
	Repo    string `yaml:",omitempty"`
}

func convertFromYaml(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error during file reading %+v", err)
		os.Exit(-1)
	}

	var fileStruct FileStruct
	err = yaml.Unmarshal(data, &fileStruct)
	if err != nil {
		fmt.Printf("Error during file unmarshal %+v", err)
		os.Exit(-1)
	}

	sort.Slice(fileStruct.Imports, func(i, j int) bool { return fileStruct.Imports[i].Package < fileStruct.Imports[j].Package })
	fmt.Println(fileStruct.Package)
	for _, v := range fileStruct.Imports {
		fmt.Println(v.Package, v.Version, v.Repo)
	}
}

func convertFromConf(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open file, %+v", err)
		os.Exit(-1)
	}

	var fileStruct FileStruct

	reader := bufio.NewReader(file)
	splitter := regexp.MustCompile(`\t+| +`)
	packageNameLine := false // Only for vendor.conf of trash
	for {
		line, err := Readln(reader)
		if err != nil && err != io.EOF {
			fmt.Printf("Could not read line, %+v", err)
			os.Exit(-1)
		} else if err == io.EOF {
			break
		} else if packageNameLine {
			packageNameLine = false
			fileStruct.Package = line
			continue
		}

		// Detecting trash configuration file
		if strings.HasPrefix(line, "# package") {
			packageNameLine = true
			continue
		} else if strings.HasPrefix(line, "#") {
			// Comment section - ignoring it
			continue
		}

		lines := splitter.Split(line, -1)
		repo := ""
		if len(lines) < 2 {
			// TODO: Handle if regexp produced less than 2 results for line
			continue
		} else if len(lines) == 3 {
			repo = lines[2]
		} else if len(lines) > 3 {
			fmt.Printf("More then three objects on line, %+v", lines)
			os.Exit(-1)
		}

		fileStruct.Imports = append(fileStruct.Imports, Import{
			Package: lines[0],
			Version: lines[1],
			Repo:    repo,
		})
	}

	sort.Slice(fileStruct.Imports, func(i, j int) bool { return fileStruct.Imports[i].Package < fileStruct.Imports[j].Package })

	if fileStruct.Package == "" {
		fileStruct.Package = getPkg(filename)
	}

	data, err := yaml.Marshal(fileStruct)
	if err != nil {
		fmt.Printf("Could not marshal, %+v", err)
		os.Exit(-1)
	}

	fmt.Print(string(data))
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func getPkg(path string) string {
	gopathSrc := filepath.Join(os.Getenv("GOPATH"), "src")
	relPath, err := filepath.Rel(gopathSrc, path)
	if err != nil {
		println("Could not get relative path %+v", err)
		os.Exit(-1)
	} else if strings.Contains(relPath, "..") {
		println("Relative path contains ..", relPath)
		os.Exit(-1)
	}

	pkg := filepath.Dir(relPath)
	if pkg == "." {
		println("Could not get final directory ", relPath)
		os.Exit(-1)
	}

	return pkg
}

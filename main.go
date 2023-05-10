package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the directory path")
		return
	}
	dir := os.Args[1]

	functionsFile := "functions_to_rename.txt"
	functions, err := readFunctionList(functionsFile)
	if err != nil {
		fmt.Printf("Error reading function list: %v\n", err)
		return
	}

	functionMap := generateFunctionMap(functions)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			newContent := replaceFunctionNames(string(content), functionMap)

			err = ioutil.WriteFile(path, []byte(newContent), 0)
			if err != nil {
				return err
			}

			fmt.Printf("Processed file: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error traversing directory: %v\n", err)
	}

	fmt.Println("Done")
}

func readFunctionList(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	functions := strings.Split(string(content), "\n")
	var trimmedFunctions []string
	for _, function := range functions {
		function = strings.TrimSpace(function)
		if function != "" {
			trimmedFunctions = append(trimmedFunctions, function)
		}
	}

	return trimmedFunctions, nil
}

func generateFunctionMap(functions []string) map[string]string {
	functionMap := make(map[string]string)
	for _, function := range functions {
		newName := generateRandomName(25)
		functionMap[function] = newName
	}
	return functionMap
}

func replaceFunctionNames(content string, functionMap map[string]string) string {
	funcPattern := fmt.Sprintf(`\b(?:%s)\b`, strings.Join(getSortedKeys(functionMap), "|"))
	r := regexp.MustCompile(funcPattern)

	newContent := r.ReplaceAllStringFunc(content, func(match string) string {
		if strings.HasPrefix(match, ".") {
			return match
		}
		return functionMap[match]
	})

	return newContent
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomName(length int) string {
	return StringWithCharset(length, charset)
}

func getSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

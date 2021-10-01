package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var image = flag.String("image", "", "The image to change the tag for, E.G. 'snarlysodboxer/my-image'")
var newTag = flag.String("newTag", "", "The new tag to set for this image, E.G. '1.2.3'")

func main() {
	flag.Parse()
	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}
	filePaths := getFilePathsFromStdin()
	if len(filePaths) < 1 {
		log.Fatal("Error: no file paths passed to stdin")
	}

	for _, filePath := range filePaths {
		contents, mode, err := readFileAndPermissions(filePath)
		if err != nil {
			log.Fatal(err)
		}
		contents = updateContents(contents, image, newTag)
		err = ioutil.WriteFile(filePath, contents, mode)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Processed: %s\n", filePath)
	}

}

func validateFlags() error {
	if *image == "" {
		return errors.New("'-image' flag cannot be empty")
	}
	if *newTag == "" {
		return errors.New("'-newTag' flag cannot be empty")
	}

	return nil
}

func getFilePathsFromStdin() []string {
	filePaths := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) != 0 {
			filePaths = append(filePaths, text)
		} else {
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
	}

	return filePaths
}

func readFileAndPermissions(filePath string) ([]byte, os.FileMode, error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, 0644, err
	}
	file, err := os.Stat(filePath)
	if err != nil {
		return []byte{}, 0644, err
	}

	return contents, file.Mode(), nil
}

func updateContents(contents []byte, img, tag *string) []byte {
	safeImage := strings.ReplaceAll(*img, `/`, `\/`)
	regexString := fmt.Sprintf(`image:\s+%s:\S+`, safeImage)
	regex := regexp.MustCompile(regexString)
	return []byte(regex.ReplaceAllString(string(contents), fmt.Sprintf("image: %s:%s", *img, *tag)))
}

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
var searchRegex = flag.String("searchRegex", `image:\s+%s:\S+`, `The regex to use when searching for the image. The %s will be replaced by the escaped -image flag's value.`)
var replacementFormat = flag.String("replacementFormat", `image: %s:%s`, `The format string to use when creating the replacement string. The %s's will be replaced by the -image and -newTag flags' values.`)

func main() {
	flag.Parse()
	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}

	filePaths, err := getFilePathsFromStdin()
	if err != nil {
		log.Fatal(err)
	}

	for _, filePath := range filePaths {
		contents, mode, err := readFileAndPermissions(filePath)
		if err != nil {
			log.Fatal(err)
		}
		newContents := updateContents(contents, image, newTag)
		if string(newContents) != string(contents) {
			err = ioutil.WriteFile(filePath, newContents, mode)
			if err != nil {
				log.Fatal(err)
			}
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
	if !strings.Contains(*searchRegex, `%s`) {
		return errors.New("'-searchRegex' must contain '%s'")
	}
	if !strings.Contains(*replacementFormat, `%s`) {
		return errors.New("'-replacementFormat' must contain '%s' twice")
	}

	return nil
}

func getFilePathsFromStdin() ([]string, error) {
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
		return filePaths, err
	}
	if len(filePaths) < 1 {
		return filePaths, errors.New("Error: no file paths passed to stdin")
	}

	return filePaths, nil
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
	regexString := fmt.Sprintf(*searchRegex, safeImage)
	regex := regexp.MustCompile(regexString)
	// convert raw \n to actual \n
	*replacementFormat = strings.Replace(*replacementFormat, `\n`, "\n", -1)
	replacementString := fmt.Sprintf(*replacementFormat, *img, *tag)

	return []byte(regex.ReplaceAllString(string(contents), replacementString))
}

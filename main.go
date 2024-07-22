package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var nyCodesMap = map[string]bool{
	/*
		These are the file codes that are associated with NY. This can be cross
		checked on the website [https://www.anthem.com/machine-readable-file/search/].

		The NY code for the file can change but the 4 char code for the state will not
	*/
	"39B0": true,
	"71A0": true,
	"72A0": true,
	"42B0": true,
}

var validLinks = map[string]bool{}

func main() {
	buff := make([]byte, 0, 16)
	networkFiles := InNetworkFiles{}
	reader, err := os.Open("2024-07-01_anthem_index.json.gz")
	if err != nil {
		log.Println("unable to open file... is is located in the root folder of the project?")
		log.Fatal(err)
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		fmt.Println("error with creating a reader")
		log.Fatal(err)
	}
	defer archive.Close()

	scanner := bufio.NewScanner(archive)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		buff = append(buff, scanner.Bytes()...)
		if len(buff) == 16 {
			if string(buff) == "in_network_files" {
				captureLinks(scanner, &networkFiles)
				// clear the struct and buffer
				networkFiles = InNetworkFiles{}
				buff = buff[:0]
			} else {
				//sliding window
				buff = append(buff[:0], buff[1:]...)
			}
		}
	}

	var link Links
	link = make(Links, 1)
	for k := range validLinks {
		link[0].Link = k
		link = append(link, link[0])
	}
	//there is a duplicate entry in the slice at the beginning, this deletes that entry
	link = append(link[:0], link[1:]...)
	bytes, _ := jsonMarshal(link)

	resultsFile, _ := os.Create("results.json")
	os.WriteFile("results.json", bytes, 0644)
	resultsFile.Close()
}

type InNetworkFiles []struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}

type Links []struct {
	Link string `json:"link"`
}

func captureLinks(scanner *bufio.Scanner, obj *InNetworkFiles) {
	defer safeExit("hmm, found a link that is does not pertain to in_network_files")
	buf := new(bytes.Buffer)

	idx := 0
	for scanner.Scan() {
		//to prevent a leading `:` from being added to the buffer
		if idx > 1 {
			buf.WriteByte(scanner.Bytes()[0])
			if scanner.Text() == "]" {
				break
			}
		}
		idx++
	}

	json.NewDecoder(buf).Decode(obj)

	// find location
	var arr []string
	for _, val := range *obj {
		arr = strings.Split(val.Location, "_")
		if ok := nyCodesMap[arr[2]]; ok {
			validLinks[val.Location] = true
		}
	}
}

// to recover from painicing
func safeExit(s string) {
	if r := recover(); r != nil {
		log.Println(s)
	}
}

/*
to prevent escaping the `&`.
curtesy of [https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and]
*/
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

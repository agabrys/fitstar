package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	step1()
}

func step1() {
	body, _ := fetchPage(`https://fit-star.easy2book.de/`)
	log.Println(findFormToken(body))
	log.Print(findTrainings(`FIT STAR München-Neuhausen`, body))
}

func fetchPage(url string) (string, []*http.Cookie) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body), resp.Cookies()
}

func findFormToken(html string) (string, error) {
	regex := regexp.MustCompile(`.*<input type="hidden" name="(\w+)" id="formtoken" value="1">.*`)
	groups := regex.FindStringSubmatch(html)
	if len(groups) == 0 {
		return ``, errors.New(`Couldn't find the form token`)
	}
	return groups[1], nil
}

func findTrainings(gymName, html string) (map[string]string, error) {
	tokenStart := fmt.Sprintf(`<optgroup id="\d+" label="%v">`, gymName)
	regex := regexp.MustCompile(tokenStart)
	indexes := regex.FindStringIndex(html)
	if len(indexes) == 0 {
		return nil, errors.New(fmt.Sprintf(`Couldn't find the start token: %v`, tokenStart))
	}

	text := html[indexes[1]:]
	tokenStop := `</optgroup>`
	index := strings.Index(text, tokenStop)
	if index == -1 {
		return nil, errors.New(fmt.Sprintf(`Couldn't find the stop token "%v" in %v`, tokenStop, text))
	}

	text = text[:index]

	parts := strings.Split(text, `</option>`)
	if len(parts) == 0 {
		return nil, errors.New(`No options have been found`)
	}
	parts = parts[:len(parts)-1]

	regex = regexp.MustCompile(`<option value="(\d+)">(.*)`)
	events := make(map[string]string)
	for _, row := range parts {
		groups := regex.FindStringSubmatch(row)
		if len(groups) == 0 {
			return nil, errors.New(`Couldn't find any valid options`)
		}
		events[groups[2]] = groups[1]
	}

	return events, nil
}

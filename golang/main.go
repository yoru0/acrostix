package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Data struct {
	XCoordinate int
	YCoordinate int
	Character   string
}

type Datas []Data

func main() {
	url := "https://docs.google.com/document/d/e/2PACX-1vSZ1vDD85PCR1d5QC2XwbXClC1Kuh3a4u0y3VbTvTFQI53erafhUkGot24ulET8ZRqFSzYoi3pLTGwM/pub"
	data, err := decodeFromUrl(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	data = data[3:]
	dt := insertData(data)

	mapData(dt)
}

func getArraySize(dt Datas) (int, int) {
	x, y := 0, 0
	for i := range dt {
		if x < dt[i].XCoordinate {
			x = dt[i].XCoordinate
		}
		if y < dt[i].YCoordinate {
			y = dt[i].YCoordinate
		}
	}
	return x + 1, y + 1
}

func mapData(dt Datas) {
	x, y := getArraySize(dt)
	arr := make([][]string, y)
	for i := 0; i < y; i++ {
		arr[i] = make([]string, x)
		for j := 0; j < x; j++ {
			arr[i][j] = " "
		}
	}

	for i := range dt {
		arr[dt[i].YCoordinate][dt[i].XCoordinate] = dt[i].Character
	}

	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			fmt.Print(arr[i][j])
		}
		fmt.Println()
	}

}

func insertData(d []string) Datas {
	dataLen := len(d)
	var dt Datas
	for i := 0; i < dataLen; i += 3 {
		dt = append(dt, Data{
			XCoordinate: stringToInt(d[i]),
			YCoordinate: stringToInt(d[i+2]),
			Character:   d[i+1],
		})
	}
	return dt
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return i
}

func decodeFromUrl(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseHTML(resp.Body)
}

func parseHTML(r io.Reader) ([]string, error) {
	var result []string
	z := html.NewTokenizer(r)
	inTable := false

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return result, nil
			}
			return nil, z.Err()

		case html.StartTagToken:
			t := z.Token()
			if t.Data == "table" {
				inTable = true
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "table" {
				inTable = false
			}

		case html.TextToken:
			if inTable {
				text := strings.TrimSpace(string(z.Text()))
				if text != "" {
					result = append(result, text)
				}
			}
		}
	}
}

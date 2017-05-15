package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	decode()
}

func decode() {
	content, err := ioutil.ReadFile("/Users/gb123/github.com/saleemjaffer/BitTorrentClient/src/main/test.torrent")
	if err != nil {
		log.Fatal(err)
	}

	//content := "l3:qqqe"
	//var err error

	switch string(string(content)[0]) {
	case "d":
		_, _, err := ParseDict(string(content), 0)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(d)
	case "l":
		_, _, err = ParseList(string(content), 0)
		if err != nil {
			log.Fatal(err)
		}
	case "i":
		_, _, err = ParseInt(string(content), 0)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("not yet handled strings")
	}
}

func ParseDict(input string, index int) (map[string]interface{}, int, error) {
	dict := make(map[string]interface{}, 0)
	initialIndex := index
	for {
		if string(input[index]) == "e" {
			break // end of dict
		}

		var key string
		var val interface{}
		var err error

		// This is to handle the first invocation of ParseDict. If this check
		// is not present you will end up doing infinite recursion
		if string(input[index]) == "d" && initialIndex == index {
			index++
		}

		key, index, err = parseString(input, index)
		if err != nil {
			log.Fatal(err)
		}

		// check what the type the value is
		switch string(string(input[index])) {
		case "l":
			val, index, err = ParseList(input, index)
			if err != nil {
				log.Fatal(err)
			}
			dict[key] = val
		case "d":
			val, index, err = ParseDict(input, index)
			if err != nil {
				log.Fatal(err)
			}
			dict[key] = val
		case "i":
			val, index, err = ParseInt(input, index)
			if err != nil {
				log.Fatal(err)
			}
			dict[key] = val
		default:
			_, err := strconv.Atoi(string(string(input[index])))
			if err != nil {
				log.Fatal(err)
			}
			val, index, err = parseString(input, index)
			dict[key] = val

		}
	}

	return dict, index + 1, nil
}

func ParseList(input string, index int) ([]interface{}, int, error) {
	list := make([]interface{}, 0)
	initialIndex := index
	for {
		if string(input[index]) == "e" {
			break // end of list
		}

		if string(input[index]) == "l" && initialIndex == index {
			index++
		}

		var err error
		var s interface{}
		var d map[string]interface{}
		l := make([]interface{}, 0)
		var i int

		switch string(string(input[index])) {
		case "l":
			l, index, err = ParseList(input, index)
			if err != nil {
				log.Fatal(err)
			}

			list = append(list, l)
		case "d":
			d, index, err = ParseDict(input, index)
			if err != nil {
				log.Fatal(err)
			}

			list = append(list, d)
		case "i":
			i, index, err = ParseInt(input, index)
			if err != nil {
				log.Fatal(err)
			}

			list = append(list, i)
		default:
			_, err := strconv.Atoi(string(string(input[index])))
			if err != nil {
				log.Fatal(err)
			}

			s, index, err = parseString(input, index)
			list = append(list, s)

		}

	}

	return list, index + 1, nil
}

func parseString(input string, index int) (string, int, error) {
	stringLen := 0
	initialIndex := index
	for {
		len, err := strconv.Atoi(string(input[initialIndex : index+1]))
		if len != 0 {
			stringLen = len
		}

		if stringLen == 0 && err != nil {
			return "", 0, errors.New(err.Error())
		} else if string(input[index]) == ":" {
			index++
			break
		}
		index++
	}

	return string(input[index : index+stringLen]), index + stringLen, nil

}

func ParseInt(input string, index int) (int, int, error) {
	if string(input[index]) != "i" {
		return 0, 0, errors.New("Invalid input")
	}
	index++

	var t string = ""

	for {
		if string(input[index]) != "e" {
			t = t + string(input[index])
			index++
		} else {
			break
		}

	}

	integer, err := strconv.Atoi(t)
	if err != nil {
		log.Fatal(err)
	}

	return integer, index + 1, nil

}

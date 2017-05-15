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

		if string(input[index+1]) == "l" ||
			(string(input[index]) == "l" && initialIndex != index) {
			var ll interface{}
			var err error

			ll, index, err = ParseList(input, index+1)
			if err != nil {
				log.Fatal(err)
			}
			list = append(list, ll)
		} else {
			var err error
			var str interface{}
			var dict map[string]interface{}
			var integer string

			if string(input[index]) == "l" {
				index++
			}
			if string(input[index]) == "d" {
				dict, index, err = ParseDict(input, index)
				if err != nil {
					panic("fucked up dict")
				}

				list = append(list, dict)

			} else if string(input[index]) == "i" {
				integer, index, err = ParseInt(input, index)
				if err != nil {
					panic("fucked up integer")
				}

				list = append(list, integer)

			} else {
				str, index, err = parseString(input, index)
				if err != nil {
					panic("fucked up string")
				}

				list = append(list, str)
			}

		}
	}

	//fmt.Println(list)

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

		//fmt.Println(stringLen)
		//fmt.Println(string(input[index]))

		if stringLen == 0 && err != nil {
			return "", 0, errors.New(err.Error())
		} else if string(input[index]) == ":" {
			index++
			break
		}
		index++
	}
	//fmt.Println(string(input[index : index+stringLen]))

	return string(input[index : index+stringLen]), index + stringLen, nil

}

func ParseInt(input string, index int) (string, int, error) {
	if string(input[index]) != "i" {
		panic("fucked up inside parse int")
	}
	index++
	var integer string = ""

	for {
		if string(input[index]) != "e" {
			integer = integer + string(input[index])
			index++
		} else {
			break
		}

	}

	return integer, index + 1, nil

}

package main

import (
	"errors"
	"fmt"
	"strconv"
	//"io/ioutil"
	"io/ioutil"
)

func main() {
	//s, newIndex, err := parseString("4:hell", 0)
	//if err != nil {
	//	panic("fuck up main")
	//}
	//fmt.Println(s)
	//fmt.Println(newIndex)

	//s, newIndex, err := ParseList("l4:hell6:saleeml3:hey4:home9:beautifule4:fucke", 0)
	content, err := ioutil.ReadFile("/Users/gb123/github.com/saleemjaffer/BitTorrentClient/src/main/test.torrent")

	if err != nil {
		panic(err)
	}
	//fmt.Println(string(content))
	s, _, err := ParseDict(string(content), 0)
	//s, _, err := ParseDict("d13:announce-listll2:wreee", 0)
	if err != nil {
		panic("fuck up main")
	}
	fmt.Println("success")
	fmt.Println(s)

}

func ParseDict(input string, index int) (map[string]interface{}, int, error) {
	dict := make(map[string]interface{}, 0)
	initialIndex := index
	for {
		if string(input[index]) == "e" {
			// end of dict
			break
		}

		var key string
		var val interface{}
		var err error
		if string(input[index]) == "d" && initialIndex == index {
			index++
		}

		key, index, err = parseString(input, index)
		if err != nil {
			panic("fucked up string")
		}

		// check what the type the value is
		_, err = strconv.Atoi(string(input[index]))
		if err != nil {
			if string(input[index]) == "l" {
				//fmt.Println("FUCK")
				val, index, err = ParseList(input, index)
				if err != nil {
					panic("")
				}
				dict[key] = val
			} else if string(input[index]) == "d" {
				val, index, err = ParseDict(input, index)
				if err != nil {
					panic("")
				}
				dict[key] = val
			} else {
				val, index, err = ParseInt(input, index)
				if err != nil {
					panic("fuck")
				}
				dict[key] = val
			}
		} else {
			val, index, err = parseString(input, index)
			dict[key] = val
		}

	}

	return dict, index + 1, nil
}

func ParseList(input string, index int) ([]interface{}, int, error) {
	//fmt.Println(index)
	//fmt.Println(string(input[18]))
	list := make([]interface{}, 0)
	initialIndex := index
	for {
		if string(input[index]) == "e" {
			// end of list
			break
		}

		if string(input[index+1]) == "l" || (string(input[index]) == "l" && initialIndex != index) {
			var ll interface{}
			var err error
			//fmt.Println(string(input[index]))
			ll, index, err = ParseList(input, index+1)
			if err != nil {
				panic("fucked")
			}
			list = append(list, ll)
		} else {
			var err error
			var str interface{}
			var dict map[string]interface{}
			var integer string
			//fmt.Println(string(input[index]))
			if string(input[index]) == "l" {
				index++
			}
			if string(input[index]) == "d" {
				dict, index, err = ParseDict(input, index)
				if err != nil {
					panic("fucked up dict")
				}

				list = append(list, dict)

			} else if string(input[index]) == "i"{
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
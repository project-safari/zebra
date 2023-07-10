package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var ipInUse []int

func contains(all []int, one int) bool {
	for _, each := range all {
		eachStr := fmt.Sprint(each)
		str := fmt.Sprint(one)
		if eachStr == str {
			return true
		}
	}

	return false
}

func genNumericValue() int {
	// numeric possibilities for the decimal value of the IP.
	possibilities := "0123456789"

	// the length of decimal, numeric value IP can be anywhere between 6 and 12.
	min := 6
	max := 12

	// generate the size randomly for each IP address.
	size := rand.Intn(max-min) + min

GENERATE:
	// create random IP values from the possible numbers in the possibilities string.
	b := make([]byte, size)
	for i := range b {
		b[i] = possibilities[rand.Intn(len(possibilities))]
	}

	// get the IP together as a decimal.
	decimalIP, _ := strconv.Atoi((string(b)))

	if contains(ipInUse, decimalIP) {
		goto GENERATE
	} else {
		ipInUse = append(ipInUse, decimalIP)
	}

	return decimalIP
}

func convertToHex(integerIP int) string {
	// hexIP := strconv.FormatInt(255, integerIP)
	hexIP := fmt.Sprintf("0x%x", integerIP)
	hexIP = strings.TrimLeft(hexIP, "0x")
	return hexIP
}

func getFinalIP(hexIP string) []string {
	// finalHex := ""

	l := len(hexIP)
	manyBytes := (l / 2)

	if (l % 2) != 0 {
		manyBytes = manyBytes + 1
	}

	arrForIP := make([]string, manyBytes)

	for i := 0; i <= len(hexIP); i++ {
		if i+2 <= len(hexIP) {
			piece := hexIP[i : i+2]
			arrForIP = append(arrForIP, piece)
		}
	}

	for i, j := 0, len(arrForIP)-1; i < j; i, j = i+1, j-1 {
		arrForIP[i], arrForIP[j] = arrForIP[j], arrForIP[i]
	}

	// finalHex = strings.Join(arrForIP, "")

	return arrForIP
}

func hexToIP(hexVal []string) string {
	var theIP []string
	start := 0
	end := 2
	var piece string
	var finalIP string

	for i := 0; i < (len(hexVal))/2; i++ {
		if end <= len(hexVal)-2 && start <= len(hexVal)-4 {
			theIP = hexVal[start:end]
			piece = strings.Join(theIP, "")
			decimal, err := strconv.ParseInt(piece, 16, 64)
			//in case of any error.
			if err != nil {
				panic(err)
			}
			finalIP = finalIP + fmt.Sprint(decimal) + "."

			fmt.Println(theIP)
		} else {
			theIP = hexVal[len(hexVal)-2 : len(hexVal)-0]
			piece = strings.Join(theIP, "")
			decimal, err := strconv.ParseInt(piece, 16, 64)
			//in case of any error.
			if err != nil {
				panic(err)
			}
			finalIP = finalIP + fmt.Sprint(decimal)

			fmt.Println(theIP)
		}

		start = start + 2
		end = end + 2

	}

	return ""
}

func main() {
	fmt.Println(getFinalIP(convertToHex(genNumericValue())))
}

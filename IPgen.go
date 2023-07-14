package zebra

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// nolint: gochecknoglobals // needed for later use throughout the code.
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
	maxVal := 4294967296

	// the length of decimal, numeric value IP can be anywhere between 4 and 10.
	min := 4
	max := 10

	// generate the size randomly for each IP address.
	size := rand.Intn(max-min) + min // nolint:all.

GENERATE:
	// create random IP values from the possible numbers in the possibilities string.
	b := make([]byte, size)

	for i := range b {
		b[i] = possibilities[rand.Intn(len(possibilities))] // nolint:all.
	}

	// get the IP together as a decimal.
	decimalIP, _ := strconv.Atoi((string(b)))

	if contains(ipInUse, decimalIP) || decimalIP > maxVal {
		goto GENERATE
	} else {
		ipInUse = append(ipInUse, decimalIP)
	}

	return decimalIP
}

func convertToHex(integerIP int) string {
	hexIP := fmt.Sprintf("0x%x", integerIP)

	// hex value w/o the hex-specific prefix.
	hexIP = strings.TrimLeft(hexIP, "0x")

	return hexIP
}

func genByteNum(l int) int {
	parity := 2
	manyBytes := (l / parity)

	// calculate the needed number of bytes, given the size of the initial integer.
	if (l % parity) != 0 {
		manyBytes++
	}

	return manyBytes
}

func getFinal(hexIP string) []string {
	l := len(hexIP)

	manyBytes := genByteNum(l)

	slider := 2
	start := 0

	var piece string

	// array to hold the IP pairs.
	var arrForIP []string

	// situation 1: uneven bits.
	if l%2 != 0 {
		for i := 0; i < manyBytes-1; i++ {
			if slider <= l {
				piece = hexIP[start:slider]
			}

			arrForIP = append(arrForIP, piece)

			fmt.Println(arrForIP)

			start += 2
			slider += 2
		}

		this := string(hexIP[l-1])
		arrForIP = append(arrForIP, this)
	} else { // situation 2: even number of bits.
		for i := 0; i < manyBytes; i++ {
			if slider <= l {
				piece = hexIP[start:slider]
			}

			arrForIP = append(arrForIP, piece)

			fmt.Println(arrForIP)

			start += 2
			slider += 2
		}
	}

	for i, j := 0, len(arrForIP)-1; i < j; i, j = i+1, j-1 {
		arrForIP[i], arrForIP[j] = arrForIP[j], arrForIP[i]
	}

	return arrForIP
}

func theIP(hexVal []string) string {
	var finalIP string

	base := 16
	size := 32

	// get the IP value for each pair in the hex string.
	for i := 0; i < (len(hexVal)); i++ {
		decimal, err := strconv.ParseInt(hexVal[i], base, size)

		if err != nil {
			fmt.Println(err)
		}

		if i <= len(hexVal)-2 {
			finalIP += fmt.Sprint(decimal) + "."
		} else {
			finalIP += fmt.Sprint(decimal)
		}
	}

	return finalIP
}

func AssignIP() net.IP {
	resIP := theIP(getFinal(convertToHex(genNumericValue())))

	return net.IP(resIP)
}

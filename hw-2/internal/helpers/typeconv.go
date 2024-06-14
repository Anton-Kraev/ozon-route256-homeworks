package helpers

import (
	"strconv"
	"strings"
)

func StrToUint64Arr(str string) ([]uint64, error) {
	var arr []uint64

	for _, strNum := range strings.Split(str, ",") {
		num, err := strconv.ParseUint(strNum, 10, 64)
		if err != nil {
			return []uint64{}, err
		}

		arr = append(arr, num)
	}

	return arr, nil
}

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

func TypedSliceToInterfaceSlice[T any](slice []T) []interface{} {
	res := make([]interface{}, len(slice))
	for i, v := range slice {
		res[i] = v
	}

	return res
}

package utils

import "go-novel/app/models"

func ArrayReverse(input []string) []string {
	length := len(input)
	reversed := make([]string, length)

	for i := 0; i < length; i++ {
		reversed[i] = input[length-i-1]
	}

	return reversed
}

func InArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1
	if len(array) <= 0 {
		return
	}
	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}

func InArrayInt(val int64, array []int64) (exists bool, index int) {
	exists = false
	index = -1
	if len(array) <= 0 {
		return
	}
	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}

func CategoryEquiv(categorys []*models.CategoryReg, category string) (local int64) {
	for _, value := range categorys {
		if value.Target == category {
			return value.Local
		}
	}
	return
}

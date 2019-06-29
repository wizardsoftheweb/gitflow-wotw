package main

import (
	"errors"
	"fmt"
)

func RemoveStringElementFromStringSlice(slice_to_parse []string, element_to_remove string, must_exist ...bool) ([]string, error) {
	cut_index := 0
	for _, array_element := range slice_to_parse {
		if element_to_remove != array_element {
			slice_to_parse[cut_index] = array_element
			cut_index++
		}
	}
	if 0 == len(must_exist) || !must_exist[0] {
		return slice_to_parse[:cut_index], nil
	}
	return nil, errors.New(fmt.Sprintf("'%s' was not found in the list", element_to_remove))
}

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

func CheckError(err error) {
	if nil != err {
		logrus.Fatal(err)
	}
}

func RemoveStringElementFromStringSlice(slice_to_parse []string, element_to_remove string, must_exist ...bool) ([]string, error) {
	logrus.Debug("RemoveStringElementFromStringSlice")
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

func CollectInput() {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}

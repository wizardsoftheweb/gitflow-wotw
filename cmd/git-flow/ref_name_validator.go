package gitflow

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
)

type ValidationTarget int

const (
	RefNameValidation ValidationTarget = iota
	PrefixNameValidation
	TagNameValidation
)

var (
	// https://stackoverflow.com/a/12093994
	// https://regex101.com/r/E2TCqU/3/tests
	GitReferenceRestrictionsPattern = regexp.MustCompile(`/\.|\.\.|\/\/+|[\000-\037\177 \\~^:?*[]+|(.(lock)?|/)$|^[^/]+$`)
	GitReferenceNotAt               = regexp.MustCompile(`^@$`)
	GitReferenceNoLeadingDots       = regexp.MustCompile(`^\.`)
	GitReferenceNoAtBracket         = regexp.MustCompile(`@\{`)
	GitReferenceNoSlashDot          = regexp.MustCompile(`/\.`)
	GitReferenceNoMultipleDot       = regexp.MustCompile(`\.\.+`)
	GitReferenceNoMultipleSlash     = regexp.MustCompile(`//+`)
	GitReferenceNoSpecialChars      = regexp.MustCompile(`[\000-\037\177 \~^:?*[]+`)
	GitReferenceNoDotLockEnd        = regexp.MustCompile(`\.(lock)?$`)
	GitReferenceNoSlashEnd          = regexp.MustCompile(`/$`)
	GitReferenceMustContainSlash    = regexp.MustCompile(`^[^/]+$`)
)

func PrefixValidator(refName string) error {
	if GitReferenceNotAt.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot only be '@': %s", refName))
	}
	if GitReferenceNoLeadingDots.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot lead with '.': %s", refName))
	}
	if GitReferenceNoAtBracket.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot contain '@{': %s", refName))
	}
	if GitReferenceNoSlashDot.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot contain '/.': %s", refName))
	}
	if GitReferenceNoMultipleDot.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '.': %s", refName))
	}
	if GitReferenceNoMultipleSlash.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '/': %s", refName))
	}
	if GitReferenceNoSpecialChars.MatchString(refName) {
		return errors.New(fmt.Sprintf("Contains an unallowed special character: %s", refName))
	}
	if GitReferenceNoDotLockEnd.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot end in '.', '.lock': %s", refName))
	}
	return nil
}

func ValidateRefName(refName string) error {
	if GitReferenceMustContainSlash.MatchString(refName) {
		return errors.New(fmt.Sprintf("Must contain at least one '/': %s", refName))
	}
	if GitReferenceNoSlashEnd.MatchString(refName) {
		return errors.New(fmt.Sprintf("Cannot end in '/': %s", refName))
	}
	return PrefixValidator(refName)
}

func ValidateBranchName(refName string) error {
	return ValidateRefName(fmt.Sprintf("refs/heads/%s", refName))
}

func ValidatePrefixName(refName string) error {
	return PrefixValidator(fmt.Sprintf("refs/heads/%s", refName))
}

func ValidateTagPrefix(tagPrefix string) error {
	return ValidatePrefixName(tagPrefix)
}

func PromptForInput(inputType ValidationTarget, promptMessage string, defaultValue string) string {
	var validator func(string) error
	switch inputType {
	case TagNameValidation:
		validator = ValidateTagPrefix
		break
	case PrefixNameValidation:
		validator = ValidatePrefixName
		break
	default:
		validator = ValidateBranchName
	}
	prompt := promptui.Prompt{
		Label:    promptMessage,
		Validate: validator,
		Default:  defaultValue,
	}
	result, err := prompt.Run()
	if nil != err {
		fmt.Printf("Prompt failed %v", err)
		return PromptForInput(inputType, promptMessage, defaultValue)
	}
	return result
}

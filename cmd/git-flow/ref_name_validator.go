package gitflow

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
)

type ValidationTarget int

const (
	REF_NAME_VALIDATION ValidationTarget = iota
	PREFIX_NAME_VALIDATION
	TAG_NAME_VALIDATION
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

func PrefixValidator(ref_name string) error {
	if GitReferenceNotAt.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot only be '@': %s", ref_name))
	}
	if GitReferenceNoLeadingDots.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot lead with '.': %s", ref_name))
	}
	if GitReferenceNoAtBracket.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain '@{': %s", ref_name))
	}
	if GitReferenceNoSlashDot.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain '/.': %s", ref_name))
	}
	if GitReferenceNoMultipleDot.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '.': %s", ref_name))
	}
	if GitReferenceNoMultipleSlash.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot contain multiple consecutive '/': %s", ref_name))
	}
	if GitReferenceNoSpecialChars.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Contains an unallowed special character: %s", ref_name))
	}
	if GitReferenceNoDotLockEnd.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot end in '.', '.lock': %s", ref_name))
	}
	return nil
}

func ValidateRefName(ref_name string) error {
	PrefixValidator(ref_name)
	if GitReferenceMustContainSlash.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Must contain at least one '/': %s", ref_name))
	}
	if GitReferenceNoSlashEnd.MatchString(ref_name) {
		return errors.New(fmt.Sprintf("Cannot end in '/': %s", ref_name))
	}
	return nil
}

func ValidateBranchName(ref_name string) error {
	return ValidateRefName(fmt.Sprintf("refs/heads/%s", ref_name))
}

func ValidatePrefixName(ref_name string) error {
	return PrefixValidator(fmt.Sprintf("refs/heads/%s", ref_name))
}

func ValidateTagPrefix(tag_prefix string) error {
	return ValidatePrefixName(tag_prefix)
}

func PromptForInput(inputType ValidationTarget, promptMessage string, defaultValue string) string {
	var validator func(string) error
	switch inputType {
	case TAG_NAME_VALIDATION:
		validator = ValidateTagPrefix
		break
	case PREFIX_NAME_VALIDATION:
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

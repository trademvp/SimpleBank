package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	isValidPassword = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from(必须包含) %d-%d characters(字符)", minLength, maxLength)
	}
	return nil
}

func ValidateUserName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore(必须只包含小写字母、数字或下划线)")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces(必须只包含字母或空格)")
	}
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}

	if !isValidPassword(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore(必须只包含小写字母、数字或下划线)")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 6, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address(不是有效的电子邮件地址)")
	}
	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer(必须是一个正整数)")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}

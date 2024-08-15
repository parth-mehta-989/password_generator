package password_generator

import (
	"strings"

	consts "github.com/parth_mehta_989/password_generator/src/constants"
	"github.com/parth_mehta_989/password_generator/src/utils/random"
	"github.com/pkg/errors"
)

type passBuilder struct {
	strings.Builder
}

type PasswordGenerator struct {
	// contains unexported fields

	condition           Conditions
	builder             passBuilder
	allowedSpecialChars *string
}

type Conditions struct {
	MinUppercase   int
	MinLowercase   int
	MinNumber      int
	MinSpecialChar int
	MinLength      int
	MaxLength      int
}

// NewGenerator creates a new instance of the PasswordGenerator struct.
//
/*
Parameters:
 	* cond: The minimum condition for the password.
	* allowedSpecialChars: A pointer to a string containing the allowed special characters. If nil, will take from `!@#$`
Returns:
	* A pointer to the newly created PasswordGenerator.
*/
func NewGenerator(conditions Conditions, allowedSpecialChars *string) *PasswordGenerator {
	if conditions.MinLength == 0 {
		conditions.MinLength = consts.DefaultMinLength
	}
	if conditions.MinLength > conditions.MaxLength {
		conditions.MaxLength = conditions.MinLength
	}

	return &PasswordGenerator{condition: conditions, allowedSpecialChars: allowedSpecialChars}
}

// Generate generates a password that meets the minimum condition set in the PasswordGenerator.
//
// It initializes the passBuilder, satisfies the minimum condition by adding the required number of characters to the password,
// generates a random password length, and adds random characters to the password until the desired length is reached.
//
// Returns a pointer to the generated password as a string and an error if the password generation fails.
func (p *PasswordGenerator) Generate() (*string, error) {
	p.builder = passBuilder{}
	if err := p.satisfyMinimumCondition(); err != nil {
		return nil, errors.Wrapf(err, "error satisfying minimum condition of one CAP, one lowercase, one number, and one special char")
	}

	specialCharList := consts.SpecialChars
	if p.allowedSpecialChars != nil {
		specialCharList = *p.allowedSpecialChars
	}
	eligibleChars := strings.Join([]string{consts.UppercaseLetters, consts.LowercaseLetters, consts.Numbers, specialCharList}, "")
	passLength, err := random.RandBetween(p.condition.MinLength, p.condition.MaxLength)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating password length")
	}
	for p.builder.Len() < int(passLength) {
		if err := p.addOneChar(eligibleChars); err != nil {
			return nil, errors.Wrap(err, "error adding random character")
		}
	}
	pass := p.builder.String()
	return &pass, nil
}

// addOneChar adds a random character from the given string to the password.
//
// Parameters:
//   - letters: The string containing the characters to choose from.
//
// Returns:
//   - An error if the operation fails.
func (p *PasswordGenerator) addOneChar(letters string) error {
	index, err := random.RandomNumber(len(letters))
	if err != nil {
		return err
	}
	p.builder.Grow(1)
	p.builder.WriteByte(letters[index])
	return nil
}

// satisfyMinimumCondition ensures that the generated password meets the minimum condition set in the PasswordGenerator.
//
// It iterates over the character types (uppercase letters, lowercase letters, consts.Numbers, and special characters) and adds the required number of each type to the password.
//
// Returns an error if the operation fails.
func (p *PasswordGenerator) satisfyMinimumCondition() error {
	characterTypes := []struct {
		chars string
		count int
	}{
		{consts.UppercaseLetters, p.condition.MinUppercase},
		{consts.LowercaseLetters, p.condition.MinLowercase},
		{consts.Numbers, p.condition.MinNumber},
		{p.getSpecialChars(), p.condition.MinSpecialChar},
	}

	for _, charType := range characterTypes {
		for i := 0; i < charType.count; i++ {
			if err := p.addOneChar(charType.chars); err != nil {
				return errors.Wrapf(err, "error adding minimum %s characters", charType.chars)
			}
		}
	}

	return nil
}

// getSpecialChars returns the special characters allowed in the PasswordGenerator.
//
// If AllowedSpecialChars is not nil, it returns the value of AllowedSpecialChars.
// Otherwise, it returns the default value of specialChars.
//
// Returns:
//   - string: The special characters allowed in the PasswordGenerator.
func (p *PasswordGenerator) getSpecialChars() string {
	if p.allowedSpecialChars != nil {
		return *p.allowedSpecialChars
	}
	return consts.SpecialChars
}

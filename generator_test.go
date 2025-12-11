package password_generator

import (
	"regexp"
	"strings"
	"testing"
)

// TestNewGenerator tests the NewGenerator function with various conditions.
func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name                string
		inputConditions     Conditions
		inputAllowedSpecial *string
		expectedMinLength   int
		expectedMaxLength   int
		expectedSpecialChar string
	}{
		{
			name: "Default lengths and special chars",
			inputConditions: Conditions{
				MinUppercase:   1,
				MinLowercase:   1,
				MinNumber:      1,
				MinSpecialChar: 1,
			},
			inputAllowedSpecial: nil,
			expectedMinLength:   defaultMinLength,
			expectedMaxLength:   defaultMaxLength,
			expectedSpecialChar: specialChars,
		},
		{
			name: "Custom lengths",
			inputConditions: Conditions{
				MinUppercase:   1,
				MinLowercase:   1,
				MinNumber:      1,
				MinSpecialChar: 1,
				MinLength:      10,
				MaxLength:      20,
			},
			inputAllowedSpecial: nil,
			expectedMinLength:   10,
			expectedMaxLength:   20,
			expectedSpecialChar: specialChars,
		},
		{
			name: "MinLength greater than MaxLength",
			inputConditions: Conditions{
				MinUppercase:   1,
				MinLowercase:   1,
				MinNumber:      1,
				MinSpecialChar: 1,
				MinLength:      20,
				MaxLength:      10,
			},
			inputAllowedSpecial: nil,
			expectedMinLength:   20,
			expectedMaxLength:   21, // MinLength + 1
			expectedSpecialChar: specialChars,
		},
		{
			name: "Custom allowed special characters",
			inputConditions: Conditions{
				MinUppercase:   1,
				MinLowercase:   1,
				MinNumber:      1,
				MinSpecialChar: 1,
			},
			inputAllowedSpecial: func() *string { s := "!@#$"; return &s }(),
			expectedMinLength:   defaultMinLength,
			expectedMaxLength:   defaultMaxLength,
			expectedSpecialChar: "!@#$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewGenerator(tt.inputConditions, tt.inputAllowedSpecial)

			if pg.condition.MinLength != tt.expectedMinLength {
				t.Errorf("NewGenerator() MinLength = %v, want %v", pg.condition.MinLength, tt.expectedMinLength)
			}
			if pg.condition.MaxLength != tt.expectedMaxLength {
				t.Errorf("NewGenerator() MaxLength = %v, want %v", pg.condition.MaxLength, tt.expectedMaxLength)
			}
			if pg.getSpecialChars() != tt.expectedSpecialChar {
				t.Errorf("NewGenerator() getSpecialChars = %v, want %v", pg.getSpecialChars(), tt.expectedSpecialChar)
			}
		})
	}
}

// TestGenerate_DefaultConditions tests the Generate method with default conditions.
func TestGenerate_DefaultConditions(t *testing.T) {
	cond := Conditions{
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumber:      1,
		MinSpecialChar: 1,
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) < defaultMinLength || len(pw) > defaultMaxLength {
		t.Errorf("Generated password length %d is not within default range [%d, %d]", len(pw), defaultMinLength, defaultMaxLength)
	}

	// Verify minimum conditions
	if !regexp.MustCompile(`[A-Z]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain uppercase letters", pw)
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain lowercase letters", pw)
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain numbers", pw)
	}
	if !regexp.MustCompile(`[!@#$]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain special characters", pw)
	}
	for _, char := range pw {
		charStr := string(char)
		if !strings.ContainsAny(charStr, uppercaseLetters+lowercaseLetters+numbers+specialChars) {
			t.Errorf("Password '%s' contains disallowed character '%s'", pw, charStr)
		}
	}
}

// TestGenerate_CustomConditions tests the Generate method with custom conditions.
func TestGenerate_CustomConditions(t *testing.T) {
	tests := []struct {
		name        string
		conditions  Conditions
		expectedLen func(int) bool // function to check if length is valid
	}{
		{
			name: "All minimums set to 2",
			conditions: Conditions{
				MinUppercase:   2,
				MinLowercase:   2,
				MinNumber:      2,
				MinSpecialChar: 2,
				MinLength:      10,
				MaxLength:      15,
			},
			expectedLen: func(l int) bool { return l >= 10 && l <= 15 },
		},
		{
			name: "Minimum length only",
			conditions: Conditions{
				MinLength: 12,
				MaxLength: 12,
			},
			expectedLen: func(l int) bool { return l == 12 },
		},
		{
			name: "Only uppercase and numbers",
			conditions: Conditions{
				MinUppercase: 3,
				MinNumber:    3,
				MinLength:    8,
				MaxLength:    10,
			},
			expectedLen: func(l int) bool { return l >= 8 && l <= 10 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewGenerator(tt.conditions, nil)
			password, err := pg.Generate()
			if err != nil {
				t.Fatalf("Generate() failed: %v", err)
			}
			if password == nil {
				t.Fatal("Generate() returned nil password")
			}

			pw := *password
			if !tt.expectedLen(len(pw)) {
				t.Errorf("Generated password length %d is not within expected range", len(pw))
			}

			// Verify minimum conditions
			upperCount := 0
			lowerCount := 0
			numberCount := 0
			specialCount := 0

			for _, char := range pw {
				if strings.ContainsRune(uppercaseLetters, char) {
					upperCount++
				} else if strings.ContainsRune(lowercaseLetters, char) {
					lowerCount++
				} else if strings.ContainsRune(numbers, char) {
					numberCount++
				} else if strings.ContainsRune(specialChars, char) {
					specialCount++
				} else {
					t.Errorf("Password '%s' contains disallowed character '%s'", pw, string(char))
				}
			}

			if upperCount < tt.conditions.MinUppercase {
				t.Errorf("Password '%s' has %d uppercase, expected at least %d", pw, upperCount, tt.conditions.MinUppercase)
			}
			if lowerCount < tt.conditions.MinLowercase {
				t.Errorf("Password '%s' has %d lowercase, expected at least %d", pw, lowerCount, tt.conditions.MinLowercase)
			}
			if numberCount < tt.conditions.MinNumber {
				t.Errorf("Password '%s' has %d numbers, expected at least %d", pw, numberCount, tt.conditions.MinNumber)
			}
			if specialCount < tt.conditions.MinSpecialChar {
				t.Errorf("Password '%s' has %d special chars, expected at least %d", pw, specialCount, tt.conditions.MinSpecialChar)
			}
		})
	}
}

// TestGenerate_CustomSpecialChars tests the Generate method with custom allowed special characters.
func TestGenerate_CustomSpecialChars(t *testing.T) {
	customSpecial := "~!@#$%^&*"
	cond := Conditions{
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumber:      1,
		MinSpecialChar: 1,
		MinLength:      10,
		MaxLength:      15,
	}
	pg := NewGenerator(cond, &customSpecial)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	// Verify that special characters are from the custom set
	foundCustomSpecial := false
	for _, char := range pw {
		if strings.ContainsRune(customSpecial, char) {
			foundCustomSpecial = true
		} else if !strings.ContainsRune(uppercaseLetters+lowercaseLetters+numbers, char) {
			t.Errorf("Password '%s' contains disallowed special character '%s'", pw, string(char))
		}
	}
	if !foundCustomSpecial {
		t.Errorf("Password '%s' does not contain any of the custom special characters '%s'", pw, customSpecial)
	}

	// Also verify general conditions again
	if !regexp.MustCompile(`[A-Z]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain uppercase letters", pw)
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain lowercase letters", pw)
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(pw) {
		t.Errorf("Password '%s' does not contain numbers", pw)
	}
}

// TestGenerate_NoConditions tests generating password with no minimum conditions for character types.
func TestGenerate_NoConditions(t *testing.T) {
	cond := Conditions{
		MinLength: 10,
		MaxLength: 10,
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) != 10 {
		t.Errorf("Generated password length %d is not 10", len(pw))
	}
	// With no minimum conditions, the password should still only contain eligible characters
	eligible := uppercaseLetters + lowercaseLetters + numbers + specialChars
	for _, char := range pw {
		if !strings.ContainsRune(eligible, char) {
			t.Errorf("Password '%s' contains disallowed character '%s'", pw, string(char))
		}
	}
}

// TestGenerate_MinLengthOnly tests generating password with only minimum length.
func TestGenerate_MinLengthOnly(t *testing.T) {
	cond := Conditions{
		MinLength: 10,
		MaxLength: 10,
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) != 10 {
		t.Errorf("Generated password length %d is not 10", len(pw))
	}
}

// TestGenerate_ZeroMinMaxLengths tests generating password with zero min/max lengths.
func TestGenerate_ZeroMinMaxLengths(t *testing.T) {
	cond := Conditions{
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumber:      1,
		MinSpecialChar: 1,
		MinLength:      0, // Should default to 8
		MaxLength:      0, // Should default to 15
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) < defaultMinLength || len(pw) > defaultMaxLength {
		t.Errorf("Generated password length %d is not within default range [%d, %d]", len(pw), defaultMinLength, defaultMaxLength)
	}
}

// TestGenerate_MinLengthGreaterThanMaxLength tests scenario where MinLength > MaxLength.
func TestGenerate_MinLengthGreaterThanMaxLength(t *testing.T) {
	cond := Conditions{
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumber:      1,
		MinSpecialChar: 1,
		MinLength:      15,
		MaxLength:      10, // Invalid, should be adjusted to 16 (MinLength + 1)
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) < 15 || len(pw) > 16 { // Expected range after adjustment
		t.Errorf("Generated password length %d is not within adjusted range [%d, %d]", len(pw), 15, 16)
	}
}

// TestGenerate_HighMinimums tests generating password with high minimum character requirements.
func TestGenerate_HighMinimums(t *testing.T) {
	cond := Conditions{
		MinUppercase:   5,
		MinLowercase:   5,
		MinNumber:      5,
		MinSpecialChar: 5,
		MinLength:      20,
		MaxLength:      25,
	}
	pg := NewGenerator(cond, nil)

	password, err := pg.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if password == nil {
		t.Fatal("Generate() returned nil password")
	}

	pw := *password
	if len(pw) < cond.MinLength || len(pw) > cond.MaxLength {
		t.Errorf("Generated password length %d is not within range [%d, %d]", len(pw), cond.MinLength, cond.MaxLength)
	}

	upperCount := 0
	lowerCount := 0
	numberCount := 0
	specialCount := 0

	for _, char := range pw {
		if strings.ContainsRune(uppercaseLetters, char) {
			upperCount++
		} else if strings.ContainsRune(lowercaseLetters, char) {
			lowerCount++
		} else if strings.ContainsRune(numbers, char) {
			numberCount++
		} else if strings.ContainsRune(specialChars, char) {
			specialCount++
		}
	}

	if upperCount < cond.MinUppercase {
		t.Errorf("Password '%s' has %d uppercase, expected at least %d", pw, upperCount, cond.MinUppercase)
	}
	if lowerCount < cond.MinLowercase {
		t.Errorf("Password '%s' has %d lowercase, expected at least %d", pw, lowerCount, cond.MinLowercase)
	}
	if numberCount < cond.MinNumber {
		t.Errorf("Password '%s' has %d numbers, expected at least %d", pw, numberCount, cond.MinNumber)
	}
	if specialCount < cond.MinSpecialChar {
		t.Errorf("Password '%s' has %d special chars, expected at least %d", pw, specialCount, cond.MinSpecialChar)
	}
}

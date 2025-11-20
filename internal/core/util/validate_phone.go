package util

import "regexp"

// ValidatePhone validates a phone number in Brazilian or American format
func ValidatePhone(phone string) bool {
	// Validate phone number in Brazilian or American format
	// Brazilian format: +55 followed by 2-digit area code + 8-9 digit number (10-11 digits total)
	// American format: +1 followed by 10-digit number
	brazilianRegex := `^\+55\d{10,11}$`
	americanRegex := `^\+1\d{10}$`

	brazilianMatch, _ := regexp.MatchString(brazilianRegex, phone)
	americanMatch, _ := regexp.MatchString(americanRegex, phone)

	return brazilianMatch || americanMatch
}

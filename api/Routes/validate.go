/*
  Form Validation Functions -- used by routes
  --NOTE: I'm pretty bad at Regexp stuff, trying to get better,
          but here's the source for the Username one:
          http://stackoverflow.com/a/12019115/6137013
 */
package routes

import (
    "regexp"
    "unicode"
)

// function to validate villanova email addresses
func validateEmail(email string) bool {
 	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)
 	return regex.MatchString(email)
}

// function to validate password
func validatePassword(password string) bool {
  var sevenPlus bool = len([]rune(password)) > 7
  var chars bool = false
  var nums bool = false
  for _, c := range password {
    if unicode.IsLetter(c) {
      chars = true
    }
    if unicode.IsNumber(c) {
      nums = true
    }
    if nums && chars {
      break
    }
  }

  //has to have all requirements
  return sevenPlus && chars && nums
}

// function to validate username field, group name field -- 15 length
func validateGeneral(try string) bool {
  regex := regexp.MustCompile(`[^a-zA-Z0-9]`)
  return !regex.MatchString(try) && len(string) < 15
}

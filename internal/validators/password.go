package validators

import "regexp"

const passwordRegexp = `[a-zA-Z0-9<>%!.,'\/@#$*()]{7,30}`

func IsPasswordValid(login string) (bool, error) {
	return regexp.MatchString(passwordRegexp, login)
}

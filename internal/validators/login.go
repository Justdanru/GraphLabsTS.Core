package validators

import "regexp"

const loginRegexp = `[a-zA-Z0-9]{5,20}`

func IsLoginValid(login string) (bool, error) {
	return regexp.MatchString(loginRegexp, login)
}

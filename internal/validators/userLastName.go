package validators

import "regexp"

const userLastNameRegexp = `[ёЁА-я\-]{0,50}`

func IsUserLastNameValid(userLastName string) (bool, error) {
	return regexp.MatchString(userLastNameRegexp, userLastName)
}

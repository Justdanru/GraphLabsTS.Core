package validators

import "regexp"

const userNameRegexp = `[ёЁА-я\-]{1,50}`

func IsUserNameValid(userName string) (bool, error) {
	return regexp.MatchString(userNameRegexp, userName)
}

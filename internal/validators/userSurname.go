package validators

import "regexp"

const userSurnameRegexp = `[ёЁА-я\-]{1,50}`

func IsUserSurnameValid(userSurname string) (bool, error) {
	return regexp.MatchString(userSurnameRegexp, userSurname)
}

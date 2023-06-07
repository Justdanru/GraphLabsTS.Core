package validators

import "regexp"

const groupNameRegexp = `[a-zA-ZёЁА-я0-9\-]{1,15}`

func IsGroupNameValid(groupName string) (bool, error) {
	return regexp.MatchString(groupNameRegexp, groupName)
}

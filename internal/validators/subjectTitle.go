package validators

import "regexp"

const subjectTitleRegexp = `[a-zA-ZёЁА-я0-9\-()\s]{1,150}`

func IsSubjectTitleValid(subjectTitle string) (bool, error) {
	return regexp.MatchString(subjectTitleRegexp, subjectTitle)
}

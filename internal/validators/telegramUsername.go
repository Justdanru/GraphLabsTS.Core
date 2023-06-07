package validators

import "regexp"

const telegramUsernameRegexp = `[a-zA-Z0-9_]{5,32}`

func IsTelegramUsernameValid(telegramUsername string) (bool, error) {
	return regexp.MatchString(telegramUsernameRegexp, telegramUsername)
}

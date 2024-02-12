package common

func StringEmptyOrNil(value *string) bool {
	if value != nil {
		if *value != "" {
			return false
		}
	}
	return true
}

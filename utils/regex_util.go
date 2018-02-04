package utils

import "regexp"

type Regex_Util struct {
}

func (r *Regex_Util) IsNumber(text string) bool {
	var partten = "^\\d+$"
	re := regexp.MustCompile(partten)
	return re.Match([]byte(text))
}

var RegexUtil = Regex_Util{}

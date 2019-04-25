package bnet

import "strings"

func parseTag(tag string) map[string]string {
	out := make(map[string]string)
	if tag == "-" || tag == "" {
		return out
	}

	tags := strings.Split(tag, ",")

	for i := 0; i < len(tags); i++ {
		keyVal := strings.Split(tags[i], "-")
		if len(keyVal) == 1 {
			out[keyVal[0]] = ""
		} else {
			out[keyVal[0]] = keyVal[1]
		}
	}

	return out
}

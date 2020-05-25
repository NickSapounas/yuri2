package util

import (
	"fmt"
	"mime"
	"strconv"
	"strings"
)

func GetExtension(fileName string) string {
	iExt := strings.LastIndex(fileName, ".")
	if iExt < 0 {
		return ""
	}

	return fileName[iExt+1:]
}

func GetMimeType(name string) string {
	return mime.TypeByExtension(fmt.Sprintf(".%s", GetExtension(name)))
}

func StringArrayContains(arr []string, v string) bool {
	for _, vi := range arr {
		if vi == v {
			return true
		}
	}

	return false
}

func AtoiDef(str string, def int) int {
	if str == "" {
		return def
	}

	if i, err := strconv.Atoi(str); err == nil {
		return i
	}

	return def
}

package lib

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Helper function
func Filter[T any](slice []T, test func(T) bool) []T {
	var ret = make([]T, 0)
	for _, v := range slice {
		if test(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Check if a file or directory is hidden. Also ignore OS files that definitely shouldn't be served publicly
func IsHidden(e os.DirEntry) bool {
	return strings.HasPrefix(e.Name(), ".") ||
		strings.HasPrefix(e.Name(), "_") ||
		regexp.MustCompile(`(?i)desktop\.ini`).MatchString(e.Name())
}

func Includes(entry string, ignoreList []string) bool {
	for _, item := range ignoreList {
		if item == entry {
			return true
		}
	}
	return false
}

// Log error and move on to the next entry. I don't know if you can continue a loop from within here so just PLEASE use 'continue' right after it in the error check!!!
func LogError(err error) {
	redBg := color.New(color.BgRed).SprintFunc()
	if err != nil {
		log.Printf("%s %s", redBg(" ERROR "), err)
	}
}

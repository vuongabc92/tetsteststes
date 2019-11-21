package helpers

import (
	"github.com/vuongabc92/octocv/lang"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"strings"
	"time"
)

func FilterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func IsMongoNoDocumentError(e error) bool {
	return e == mongo.ErrNoDocuments
}

func Trans(key interface{}, a ...interface{}) string {
	return lang.Trans.Sprintf(key, a)
}

func Now() time.Time {
	return time.Now().UTC()
}

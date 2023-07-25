package main

import (
	"fmt"
	"regexp"
)

func main() {
	var matcher = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[?|.]*`)

	matches := matcher.FindStringSubmatch(" /swagger/index.html")
	for _, match := range matches {
		fmt.Println(match)
	}
}

package net

import (
	"net/url"

	"github.com/liov/hoper/go/v2/utils/log"
)

// RelativeURLToAbsoluteURL 相对URL转绝对URL
func RelativeURLToAbsoluteURL(curURL string, baseURL string) (string, error) {
	curURLData, err := url.Parse(curURL)
	if err != nil {
		log.Error(err)
		return "", err
	}
	baseURLData, err := url.Parse(baseURL)
	if err != nil {
		log.Error(err)
		return "", err
	}
	curURLData = baseURLData.ResolveReference(curURLData)
	return curURLData.String(), nil
}

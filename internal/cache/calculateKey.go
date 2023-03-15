package cache

import (
	"bytes"
	"github.com/cespare/xxhash/v2"
	"sort"
)

type CalculateKeyData struct {
	ReqBody *bytes.Buffer
	Path    string
	Headers map[string][]string
}

func CalculateKey(input CalculateKeyData) (key uint64, err error) {
	hasher := xxhash.New()
	_, err = hasher.WriteString(input.Path)
	if err != nil {
		return
	}
	_, err = hasher.Write(input.ReqBody.Bytes())
	if err != nil {
		return
	}

	headers := make([]string, 0, len(input.Headers))
	tempValue := ""
	for key, values := range input.Headers {
		tempValue = key
		for i := range values {
			tempValue += values[i]
		}
		headers = append(headers, tempValue)
	}
	sort.Strings(headers)

	headerString := ""
	for i := range headers {
		headerString += headers[i]
	}

	_, err = hasher.WriteString(headerString)
	if err != nil {
		return
	}

	key = hasher.Sum64()
	return
}

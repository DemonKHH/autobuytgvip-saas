package main

import (
	"encoding/base64"
	"fmt"
	"regexp"
)

var payload = "te6ccgEBAgEANAABTgAAAABUZWxlZ3JhbSBQcmVtaXVtIGZvciAxIHllYXIgCgpSZWYjNgEAEFdGQ3pxbnNm"
var tonCommentFormats = map[int]string{
	3:  "Telegram Premium for 3 months \n\nRef#%s",
	6:  "Telegram Premium for 6 months \n\nRef#%s",
	12: "Telegram Premium for 1 year \n\nRef#%s",
}

func removeInvalidChars4(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-\_\.\[\]\(\)\{\}]`) //  Keep these characters
	return re.ReplaceAllString(s, "")
}
func main() {
	decodeBytes, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", string(decodeBytes))
	fragmentRef := removeInvalidChars4("Ref#6\x01\x10WFCzqnsf")
	fmt.Sprintf(tonCommentFormats[12], fragmentRef)
}

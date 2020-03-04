package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unsafe"

	"github.com/winterssy/gjson"
)

var (
	re = regexp.MustCompile(`[\\/:*?"<>|]`)
)

func TrimInvalidFilePathChars(path string) string {
	return strings.TrimSpace(re.ReplaceAllString(path, " "))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func PrettyJSON(v interface{}) string {
	s, err := gjson.EncodeToString(v, func(enc *gjson.Encoder) {
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
	})
	if err != nil {
		return "{}"
	}
	return s
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	text, _ := reader.ReadString('\n')
	input := strings.TrimSpace(text)
	if input == "" {
		return Input(prompt)
	}
	return input
}

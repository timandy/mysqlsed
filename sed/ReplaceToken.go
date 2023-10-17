package sed

import (
	"fmt"
	"strings"
)

type ReplaceToken struct {
	Source string
	Target string
}

func ResolveReplaceToken(value string) ReplaceToken {
	if !strings.Contains(value, ">>>") {
		msg := fmt.Sprintf("value: %v must contains '>>>'", value)
		panic(msg)
	}
	values := strings.Split(value, ">>>")
	if len(values) != 2 {
		msg := fmt.Sprintf("value: %v has too many sub values", value)
		panic(msg)
	}
	return ReplaceToken{Source: values[0], Target: values[1]}
}

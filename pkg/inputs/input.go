package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ReplaceAll(scanner.Text(), "\"", "")
}

func InputWithText(text string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(text)
	scanner.Scan()
	return strings.ReplaceAll(scanner.Text(), "\"", "")
}

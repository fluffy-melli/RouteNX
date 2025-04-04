package logger

import (
	"fmt"
	"os"
	"regexp"
)

type BlockLogger struct {
	OriginIP  string `json:"origin_ip"`
	ForwordIP string `json:"forword_ip"`
	Host      string `json:"host"`
	Time      string `json:"time"`
}

type Logger struct {
	Block []BlockLogger `json:"block"`
}

func (l *Logger) AddBlockLog(block BlockLogger) {
	l.Block = append(l.Block, block)
}

func NewLogger() *Logger {
	return &Logger{
		Block: make([]BlockLogger, 0),
	}
}

var colorMap = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"reset":  "\033[0m",
}

func ColorTags(input string) string {
	re := regexp.MustCompile(`\{(\w+)\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		colorName := re.ReplaceAllString(match, "$1")
		if colorCode, exists := colorMap[colorName]; exists {
			return colorCode
		}
		return match
	})
}

var Info = "{green}✔{reset}"
var Error = "{red}✘{reset}"
var Warning = "{yellow}⚠{reset}"

func INFO(format string, a ...any) {
	message := fmt.Sprintf("[%s] %s\n", Info, fmt.Sprintf(format, a...))
	fmt.Print(ColorTags(message))
}

func WARNING(format string, a ...any) {
	message := fmt.Sprintf("[%s] %s\n", Warning, fmt.Sprintf(format, a...))
	fmt.Print(ColorTags(message))
}

func ERROR(format string, a ...any) {
	message := fmt.Sprintf("[%s] %s\n", Error, fmt.Sprintf(format, a...))
	fmt.Print(ColorTags(message))
	os.Exit(1)
}

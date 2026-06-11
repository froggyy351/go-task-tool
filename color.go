package main

import "os"

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[92m"
)

// NO_COLOR 環境変数が設定されている場合は色を無効にする（標準仕様: no-color.org）。
func colorEnabled() bool {
	return os.Getenv("NO_COLOR") == ""
}

func green(s string) string {
	if !colorEnabled() {
		return s
	}
	return colorGreen + s + colorReset
}

func red(s string) string {
	if !colorEnabled() {
		return s
	}
	return colorRed + s + colorReset
}

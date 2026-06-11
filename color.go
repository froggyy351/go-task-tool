package main

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[92m"
)

func green(s string) string { return colorGreen + s + colorReset }
func red(s string) string   { return colorRed + s + colorReset }

package logger

import "fmt"

func Debug(args ...any) {
	fmt.Println(args)
}

func Error(args ...any) {
	fmt.Println(args)
}
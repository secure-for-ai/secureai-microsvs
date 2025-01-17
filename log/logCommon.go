package log

import (
	"log"
	_ "unsafe"
)

// ================= Log.Printf =================
func Errorf(msg string, v ...any) {
	log.Printf("[ERROR] " + msg, v...)
}

func Fatalf(msg string, v ...any) {
	log.Fatalf("[FATAL] " + msg, v...)
}

//go:linkname Print log.Print
func Print(v ...any)

// ================= Log.Println =================

func Errorln(v ...any) {
	log.Println(append([]any{"[ERROR]"}, v...)...)
}

func Fatalln(v ...any) {
	log.Fatalln(append([]any{"[FATAL]"}, v...)...)
}

//go:linkname Printf log.Printf
func Printf(format string, v ...any)

// ================= Log.Print =================
func Error(v ...any) {
	log.Print(append([]any{"[ERROR] "}, v...)...)
}

func Fatal(v ...any) {
	log.Fatal(append([]any{"[FATAL] "}, v...)...)
}

//go:linkname Println log.Println
func Println(v ...any)
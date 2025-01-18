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

func Print(v ...any) {
	log.Print(v...)
}

// ================= Log.Println =================

func Errorln(v ...any) {
	log.Println(append([]any{"[ERROR]"}, v...)...)
}

func Fatalln(v ...any) {
	log.Fatalln(append([]any{"[FATAL]"}, v...)...)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

// ================= Log.Print =================
func Error(v ...any) {
	log.Print(append([]any{"[ERROR] "}, v...)...)
}

func Fatal(v ...any) {
	log.Fatal(append([]any{"[FATAL] "}, v...)...)
}

func Println(v ...any) {
	log.Println(v...)
}
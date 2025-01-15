//go:build debug

package log

import (
	"log"
)

// ================= Log.Printf =================
func Debugf(format string, v ...any) {
	log.Printf("[DEBUG] " + format, v...)
}

func Warningf(format string, v ...any) {
	log.Printf("[Warning] " + format, v...)
}

func Infof(msg string, v ...any) {
	log.Printf("[INFO] " + msg, v...)
}

func Errorf(msg string, v ...any) {
	log.Printf("[ERROR] " + msg, v...)
}

// ================= Log.Println =================
func Debugln(v ...any) {
	log.Println(append([]any{"[DEBUG]"}, v...)...)
}

func Warningln(v ...any) {
	log.Println(append([]any{"[WARNING]"}, v...)...)
}

func Infoln(v ...any) {
	log.Println(append([]any{"[INFO]"}, v...)...)
}

func Errorln(v ...any) {
	log.Println(append([]any{"[ERROR]"}, v...)...)
}

// ================= Log.Print =================
func Debug(v ...any) {
	log.Print(append([]any{"[DEBUG] "}, v...)...)
}

func Warning(v ...any) {
	log.Print(append([]any{"[WARNING] "}, v...)...)
}

func Info(v ...any) {
	log.Print(append([]any{"[INFO] "}, v...)...)
}

func Error(v ...any) {
	log.Print(append([]any{"[ERROR] "}, v...)...)
}
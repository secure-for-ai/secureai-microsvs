//go:build !debug && warning


package log

import (
	"log"
)

// ================= Log.Printf =================
func Debugf(format string, v ...any) {}

func Warningf(format string, v ...any) {
	log.Printf("[Warning] " + format, v...)
}

func Infof(msg string, v ...any) {
	log.Printf("[INFO] " + msg, v...)
}

// ================= Log.Println =================
func Debugln(v ...any) {}

func Warningln(v ...any) {
	log.Println(append([]any{"[WARNING]"}, v...)...)
}

func Infoln(v ...any) {
	log.Println(append([]any{"[INFO]"}, v...)...)
}

// ================= Log.Print =================
func Debug(v ...any) {}

func Warning(v ...any) {
	log.Print(append([]any{"[WARNING] "}, v...)...)
}

func Info(v ...any) {
	log.Print(append([]any{"[INFO] "}, v...)...)
}

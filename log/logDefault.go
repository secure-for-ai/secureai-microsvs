//go:build !debug && !warning && !info
package log

import(
	"log"
)

// ================= Log.Printf =================
func Debugf(format string, v ...any) {}

func Warningf(format string, v ...any) {}

func Infof(msg string, v ...any) {}

func Errorf(msg string, v ...any) {
	log.Printf("[ERROR] " + msg, v...)
}

// ================= Log.Println =================
func Debugln(v ...any) {}

func Warningln(v ...any) {}

func Infoln(v ...any) {}

func Errorln(v ...any) {
	log.Println(append([]any{"[ERROR]"}, v...)...)
}

// ================= Log.Print =================
func Debug(v ...any) {}

func Warning(v ...any) {}

func Info(v ...any) {}

func Error(v ...any) {
	log.Print(append([]any{"[ERROR]"}, v...)...)
}
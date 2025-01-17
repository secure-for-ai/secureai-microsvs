//go:build !debug && !warning && !info
package log

// ================= Log.Printf =================
func Debugf(format string, v ...any) {}

func Warningf(format string, v ...any) {}

func Infof(msg string, v ...any) {}

// ================= Log.Println =================
func Debugln(v ...any) {}

func Warningln(v ...any) {}

func Infoln(v ...any) {}

// ================= Log.Print =================
func Debug(v ...any) {}

func Warning(v ...any) {}

func Info(v ...any) {}

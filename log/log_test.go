package log_test

import (
	log2 "log"
	"github.com/secure-for-ai/secureai-microsvs/log"
	"testing"
)

func TestLog(t *testing.T) {
	log2.Print("================= Log.Printf =================")
	log.Debugf("int %d", 1)
	log.Warningf("int %d", 1)
	log.Infof("int %d", 1)
	log.Errorf("int %d", 1)

	log2.Print("================= Log.Println =================")
	log.Debugln("int", 1)
	log.Warningln("int", 1)
	log.Infoln("int", 1)
	log.Errorln("int", 1)

	log2.Print("================= Log.Print =================")
	log.Debug("int", 1)
	log.Warning("int", 1)
	log.Info("int", 1)
	log.Error("int", 1)
}

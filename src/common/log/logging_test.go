package log

import (
	"common/config"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

const peer = "peer"
const network = "network"
const test = "testCmd"

func TestGetLogger(t *testing.T) {
	log1 := GetLogger(peer)
	fmt.Println("should be print debug/error level.... 2 lines :")
	assert.NotNil(t, log1)
	log1.Debugf("testCmd log from %s", peer)
	log1.Errorf("testCmd log from %s", peer)

	log2 := GetLogger(network)
	fmt.Println("should be print info level.... 1 lines :")
	assert.NotNil(t, log2)
	log2.Infof("testCmd log from %s", network)

	assert.NotEqual(t, log1, log2, "Not pass! should be different logger instance")

	log3 := GetLogger(peer)
	assert.NotNil(t, log3)
	assert.Equal(t, log1, log3, "Not pass! should be same logger instance")
	fmt.Println("should be print error level.... 1 lines :")
	log3.Errorf("testCmd log from %s", peer)
}

func loadConfig() *logConfig {
	var list map[string]interface{}
	list = make(map[string]interface{})
	list["peer"] = "debug"
	list["network"] = "info"
	return &logConfig{
		formatter:  "[%time%] [%lvl%] [%module%] [%file%] → %msg%",
		isDump:     false,
		dumpPath:   "D:/work/zc/trunk/01_Project/zc/data",
		dumpFormat: "text",
		timestamp:  time.RFC3339,
		moduleList: list,
	}
}
func (c *logConfig) setDump(bool bool) {
	c.isDump = bool
}
func (c *logConfig) setPath(path string) {
	c.dumpPath = path
}

func TestInitLog(t *testing.T) {
	os.Setenv("zc_path", "D:/work/zc/trunk/01_Project/zc/data")
	config.Initialize()
	conf := loadConfig()
	conf.setDump(false)
	for k, v := range conf.moduleList {
		conf.module = k
		conf.level, _ = logrus.ParseLevel(v.(string))
		resetLogInstance(conf)
	}
	assert.True(t, modules.Check(peer), "Not Pass!! fail to get logger instance")
	assert.True(t, modules.Check(network), "Not Pass!! fail to get logger instance")
	log1 := modules.Get(peer).(*logrus.Entry)
	log1.Debug("1111")
	log1.Info("2222")
	log1.Warn("3333")
	log1.Error("4444")
	//log1.Fatal("5555")
}

func TestDump(t *testing.T) {
	os.Setenv("zc_path", "D:/work/zc/trunk/01_Project/zc/data")
	config.Initialize()
	conf := loadConfig()
	conf.setDump(true)
	conf.setPath("./")
	conf.module = peer
	conf.level = logrus.DebugLevel
	resetLogInstance(conf)
	assert.True(t, modules.Check(peer), "Not Pass!! fail to get logger instance")
	log1 := modules.Get(peer).(*logrus.Entry)
	assert.NotNil(t, log1)
	log1.Debug("1111")
	log1.Info("2222")
	log1.Warn("3333")
	log1.Error("4444")
	assert.FileExists(t, "./all.log.201804030800", "Not Pass!! fail to dump to file")
	assert.FileExists(t, "./error.log.201804030800", "Not Pass!! fail to dump to file")
}
func TestInitLog2(t *testing.T) {
	os.Setenv("zc_path", "D:/work/zc/trunk/01_Project/zc/data")
	InitLog()
	assert.True(t, modules.Check(peer), "Not Pass!! fail to get logger instance")
	log1 := modules.Get(peer).(*logrus.Entry)
	log1.Debug("1111")
	log1.Info("2222")
	log1.Warn("3333")
	log1.Error("4444")
}
func TestBanch(t *testing.T) {
	os.Setenv("zc_path", "D:/work/zc/trunk/01_Project/zc/data")
	config.Initialize()
	InitLog()
	count := 100000
	log1 := GetLogger("peer")
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log1.Debug("test1111111111")
			}()
		}
	}()
	log2 := GetLogger("network")

	go func() {
		defer wg.Done()
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log2.Debug("test2222222222")
			}()
		}
	}()
	log3 := GetLogger("vm")

	go func() {
		defer wg.Done()
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log3.Debug("test33333333")
			}()
		}
	}()
	log4 := GetLogger("db")

	go func() {
		defer wg.Done()
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log4.Debug("test4444444444")
			}()
		}
	}()
	log5 := GetLogger("chain")
	go func() {
		defer wg.Done()
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log5.Debug("test5555555")
			}()
		}
	}()

	wg.Wait()
}

func TestOne(t *testing.T) {

	os.Setenv("zc_home", "D:/work/zc/trunk/01_Project/zc/src/zchome/conf")
	config.Initialize()
	InitLog()
	log1 := GetLogger("peer")
	log1.Debug("abc")
	log1.Debug("abc")
}

type Te struct {
	S []string
}

func Test(t *testing.T) {
	var tes = Te{
		S: []string{"test.example.com"},
	}
	byte, _ := json.Marshal(tes)
	fmt.Println(string(byte))
}

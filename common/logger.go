package common

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/MoonighT/GarenaHack/common/logfile"
)

var (
	infoLogger   *log.Logger
	loggerOutput *logfile.Logfile
	logLevel     int
)

func LoggerInit(filename string, rotateInterval, maxLogFileSize, maxLogFileCount int64, level int) {
	output, err := logfile.Open(filename, rotateInterval, maxLogFileSize, maxLogFileCount)
	if err != nil {
		log.Fatal(err)
	}
	logLevel = level
	infoLogger = log.New(output, "", log.LstdFlags|log.Lmicroseconds)
}

func LoggerInitBuf(filename string, rotateInterval, maxLogFileSize, maxLogFileCount, bufSize int64, level int) {
	output, err := logfile.OpenBuf(filename, rotateInterval, maxLogFileSize, maxLogFileCount, bufSize)
	if err != nil {
		log.Fatal(err)
	}
	logLevel = level
	infoLogger = log.New(output, "", log.LstdFlags|log.Lmicroseconds)
	loggerOutput = output.(*logfile.Logfile)
}

func addPrefix(prefix string, v []interface{}) []interface{} {
	var v2 []interface{}
	v2 = append(v2, prefix)
	v2 = append(v2, v...)
	return v2
}
func LogVerbose(v ...interface{}) {
	if logLevel >= 2 && infoLogger != nil {
		infoLogger.Print("[VERBOSE]", fmt.Sprint(v...))
	}
}

func getFuncNameWithoutPackage(name string) string {
	pos := strings.LastIndex(name, ".")
	if pos >= 0 {
		return name[pos+1:]
	}
	return name
}

func getActualCaller() (file string, line int, ok bool) {
	// Get func name of caller in this file
	cpc, _, _, ok := runtime.Caller(1)
	if !ok {
		return
	}
	callerFuncPtr := runtime.FuncForPC(cpc)
	if callerFuncPtr == nil {
		ok = false
		return
	}
	// Get lowest caller func info whose name
	// not same as the caller in this file
	var pc uintptr
	for callLevel := 2; callLevel < 5; callLevel++ {
		pc, file, line, ok = runtime.Caller(callLevel)
		if !ok {
			return
		}
		funcPtr := runtime.FuncForPC(pc)
		if funcPtr == nil {
			ok = false
			return
		}
		if getFuncNameWithoutPackage(funcPtr.Name()) !=
			getFuncNameWithoutPackage(callerFuncPtr.Name()) {
			return
		}
	}
	ok = false
	return
}

func LogVerbosef(format string, v ...interface{}) {
	if logLevel >= 2 {
		file, line, ok := getActualCaller()
		if ok {
			arg := make([]interface{}, 0)
			arg = append(arg, path.Base(file), line)
			arg = append(arg, v...)
			infoLogger.Printf("[VERBOSE][%s:%d]"+format, arg...)
		} else {
			infoLogger.Printf("[VERBOSE]"+format, v...)
		}
	}
}

func LogDetail(v ...interface{}) {
	if logLevel >= 1 && infoLogger != nil {
		infoLogger.Print("[DETAIL]", fmt.Sprint(v...))
	}
}

func LogDetailf(format string, v ...interface{}) {
	if infoLogger == nil {
		return
	}
	if logLevel >= 1 {
		file, line, ok := getActualCaller()
		if ok {
			arg := make([]interface{}, 0)
			arg = append(arg, path.Base(file), line)
			arg = append(arg, v...)
			infoLogger.Printf("[DETAIL][%s:%d]"+format, arg...)
		} else {
			infoLogger.Printf("[DETAIL]"+format, v...)
		}
	}
}

func LogInfo(v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Print("[INFO]", fmt.Sprint(v...))
	}
}

func LogInfof(format string, v ...interface{}) {
	if infoLogger == nil {
		return
	}
	file, line, ok := getActualCaller()
	if ok {
		arg := make([]interface{}, 0)
		arg = append(arg, path.Base(file), line)
		arg = append(arg, v...)
		infoLogger.Printf("[INFO][%s:%d]"+format, arg...)
	} else {
		infoLogger.Printf("[INFO]"+format, v...)
	}
}

func LogWarning(v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Print("[WARNING]", fmt.Sprint(v...))
	}
}

func LogWarningf(format string, v ...interface{}) {
	if infoLogger == nil {
		return
	}
	file, line, ok := getActualCaller()
	if ok {
		arg := make([]interface{}, 0)
		arg = append(arg, path.Base(file), line)
		arg = append(arg, v...)
		infoLogger.Printf("[WARNING][%s:%d]"+format, arg...)
	} else {
		infoLogger.Printf("[WARNING]"+format, v...)
	}
}
func LogError(v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Print("[ERROR]", fmt.Sprint(v...))
	}
}
func LogErrorf(format string, v ...interface{}) {
	if infoLogger == nil {
		return
	}
	file, line, ok := getActualCaller()
	if ok {
		arg := make([]interface{}, 0)
		arg = append(arg, path.Base(file), line)
		arg = append(arg, v...)
		infoLogger.Printf("[ERROR][%s:%d]"+format, arg...)
	} else {
		infoLogger.Printf("[ERROR]"+format, v...)
	}
}
func LogFatal(v ...interface{}) {
	fmt.Fprintf(os.Stderr, "[FATAL]"+fmt.Sprint(v...)+"\n")
	if infoLogger != nil {
		infoLogger.Fatal("[FATAL]", fmt.Sprint(v...))
	}
}

func Logf(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

func LogFlush() {
	if loggerOutput != nil {
		loggerOutput.BufHandle.Flush()
	}
}

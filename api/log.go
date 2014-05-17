package webl

import (
  "io"
  "io/ioutil"
  "log"
  "os" 
  "code.google.com/p/go.net/websocket"
)

var (
  TRACE   *log.Logger
  INFO    *log.Logger
  WARN    *log.Logger
  ERROR   *log.Logger
)

func InitLogging(isQuiet bool, isVerbose bool, isTimestamped bool, ws *websocket.Conn) {
  var traceHandle, infoHandle, warnHandle, errorHandle io.Writer
  var traceFormat, infoFormat, warnFormat, errorFormat int
  var tracePrefix, infoPrefix, warnPrefix, errorPrefix string

  traceHandle = ioutil.Discard
  if ws == nil {
    infoHandle = os.Stdout
    warnHandle = os.Stdout
    errorHandle = os.Stderr
  } else {
    infoHandle = ws
    warnHandle = ws
    errorHandle = ws
  }

  tracePrefix = "TRACE: "
  warnPrefix = ""
  infoPrefix = ""
  errorPrefix = "ERROR: "

  if isQuiet {
    infoHandle = ioutil.Discard
    warnHandle = ioutil.Discard
  } else if isVerbose {
    infoFormat = log.Ldate|log.Ltime
    if ws == nil {
      traceHandle = os.Stdout
    } else {
      traceHandle = ws
    }
  }

  if isTimestamped {
    traceFormat = log.Ldate|log.Ltime
    infoFormat = log.Ldate|log.Ltime
    warnFormat = log.Ldate|log.Ltime
    errorFormat = log.Ldate|log.Ltime
    warnPrefix = "WARNING: "
    infoPrefix = "INFO: "
  } else {
    traceFormat = 0
    infoFormat = 0
    warnFormat = 0
    errorFormat = 0
  }

  TRACE = log.New(traceHandle,  tracePrefix,   traceFormat)
  INFO = log.New(infoHandle,    infoPrefix,    infoFormat)
  WARN = log.New(warnHandle, warnPrefix, warnFormat)
  ERROR = log.New(errorHandle,  errorPrefix,   errorFormat)
}

func FailOnError(err error) {
  if err != nil {
    ERROR.Println("Error occurred, cannot proceed:", err)
    os.Exit(1)
  }
}



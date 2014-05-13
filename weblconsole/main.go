package main

import (
  "flag"
  "github.com/aforward/webl/api"
)

func main() {
  isVerbose     := flag.Bool("verbose",          false,   "Turn on as musch debugging information as possible")
  isQuiet       := flag.Bool("quiet",            false,   "Turn off all but the most important logging")
  isTimestamped := flag.Bool("timestamped",      false,   "Should outputs be timestamped")
  isVersion     := flag.Bool("version",          false,   "Output the version of this app")
  startingUrl   := flag.String("url",            "",      "Specify which URL to crawl (e.g. a4word.com)")
  redisServer   := flag.String("redis",          ":6379", "Specify the redis server (e.g. 127.0.0.1:6379)")
  redisPassword := flag.String("redis-password", "",      "Specify the redis server password")

  flag.Parse()

  webl.InitLogging(*isQuiet, *isVerbose, *isTimestamped)
  showVersion()
  if *isVersion {
    return
  }

  if *startingUrl == "" {
    showMissingUrl()
    return
  }

  webl.Pool = webl.NewPool(*redisServer, *redisPassword)
  webl.Crawl(*startingUrl)
}

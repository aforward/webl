package main

import (
  "io/ioutil"
  "net/http"
  "text/template"
  "fmt"
  "log"
  "github.com/realistschuckle/gohaml"
  "github.com/bmizerany/pat"
  "github.com/garyburd/redigo/redis"
  "flag"
  "github.com/aforward/webl/api"
  "gopkg.in/fatih/set.v0"
  "code.google.com/p/go.net/websocket"
)

//-----------
// HELPERS
//-----------

var webViews = initWebViews() 

func initWebViews() map[string]*template.Template {
  webViews := make(map[string]*template.Template)
  webViews["new"] = makeHtmlView("new")
  webViews["url"] = makeHtmlView("url")
  return webViews
}

func makeHtmlView(viewName string) *template.Template {
  t, _ := template.ParseFiles(fmt.Sprintf("app/views/%s.html", viewName))
  return t
}

func makeHamlView(viewName string) *template.Template {
  var scope = make(map[string]interface{})
  scope["lang"] = "HAML"
  content, _ := ioutil.ReadFile(fmt.Sprintf("app/views/%s.haml", viewName))
  engine, _ := gohaml.NewEngine(string(content))
  output := engine.Render(scope)
  return template.Must(template.New("").Parse(output))
}

func v(name string, r *http.Request) string {
  return r.URL.Query().Get(fmt.Sprintf(":%s", name)) 
}

//-----------
// VIEWS
//-----------

type Graph struct {
  Edges []Edge 
}

//-----------
// ROUTES
//-----------

func getRoot(w http.ResponseWriter, r *http.Request) {
  webViews["new"].Execute(w,nil)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
  edges := make([]Edge,10)
  root := webl.LoadDomain(v("url",r),false)
  edges = flatten(edges,&root,set.New())
  webViews["url"].Execute(w, Graph{ Edges: edges })
}

func postUrl(w http.ResponseWriter, r *http.Request) {
  webViews["url"].Execute(w, webl.Resource{ Name: "Test", Url: "TEST" })
}

func getStatic(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, r.URL.Path[1:])
}

func doCrawl(ws *websocket.Conn) {
  webl.InitLogging(*isQuiet, *isVerbose, *isTimestamped, ws)
  var url string
  websocket.Message.Receive(ws, &url)
  webl.Crawl(url)
  websocket.Message.Send(ws, "exit")
}

//-----------
// WEB SERVER
//-----------

var (
  pool *redis.Pool
  isQuiet *bool
  isVerbose *bool
  isTimestamped *bool
)

func main() {
  isVerbose     = flag.Bool("verbose",          false,   "Turn on as musch debugging information as possible")
  isQuiet       = flag.Bool("quiet",            false,   "Turn off all but the most important logging")
  isTimestamped = flag.Bool("timestamped",      false,   "Should outputs be timestamped")
  isVersion     := flag.Bool("version",          false,   "Output the version of this app")
  redisServer   := flag.String("redis",          ":6379", "Specify the redis server (default 127.0.0.1:6379)")
  redisPassword := flag.String("redis-password", "",      "Specify the redis server password")
  port          := flag.String("port",           "4005",  "Specify the web server port (default 4005)")

  flag.Parse()
  webl.InitLogging(*isQuiet, *isVerbose, *isTimestamped, nil)

  showVersion()
  if *isVersion {
    return
  }

  webl.Pool = webl.NewPool(*redisServer, *redisPassword)

  m := pat.New()
  m.Get("/static/", http.HandlerFunc(getStatic))
  m.Get("/favicon.ico", http.HandlerFunc(getStatic))
  m.Get("/u/:url", http.HandlerFunc(getUrl))
  m.Post("/", http.HandlerFunc(postUrl))
  m.Get("/", http.HandlerFunc(getRoot))

  fmt.Println(fmt.Sprintf("Starting server, accessible at http://localhost:%s", *port))
  http.Handle("/crawl", websocket.Handler(doCrawl))
  http.Handle("/", m)
  err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
  fmt.Println("Server shutting down, goodbye!")
}

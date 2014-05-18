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
var host string

func initWebViews() map[string]*template.Template {
  webViews := make(map[string]*template.Template)
  webViews["index"] = makeHtmlView("index")
  webViews["list"] = makeHtmlView("list")
  webViews["url"] = makeHtmlView("url")
  webViews["details"] = makeHtmlView("details")
  return webViews
}

func makeHtmlView(viewName string) *template.Template {
  t, _ := template.ParseFiles(
    fmt.Sprintf("app/views/%s.html", viewName),
    "app/views/_html_header.html",
    "app/views/_html_footer.html",
    "app/views/_output.html",
  )
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

type HtmlHeader struct {
  Title string
}

type HtmlFooter struct {
  CustomJs string
  AppHost string
}

type Page struct {
  HtmlHeader *HtmlHeader
  HtmlFooter *HtmlFooter
}

type ListPage struct {
  HtmlHeader *HtmlHeader
  HtmlFooter *HtmlFooter
  AllDomains []*webl.Resource
}

type DetailsPage struct {
  HtmlHeader *HtmlHeader
  HtmlFooter *HtmlFooter
  Domain *webl.Resource
}

//-----------
// ROUTES
//-----------

func homepage(w http.ResponseWriter, r *http.Request) {
  page := Page{ HtmlHeader: &HtmlHeader{ Title: "webl" }, HtmlFooter: &HtmlFooter{ AppHost: host } }
  webViews["index"].Execute(w,page)
}

func listSitemaps(w http.ResponseWriter, r *http.Request) {
  page := ListPage{ 
    HtmlHeader: &HtmlHeader{ Title: "webl -- listing sitemaps" }, 
    HtmlFooter: &HtmlFooter{ AppHost: host }, 
    AllDomains: webl.ListDomains(),
  }
  webViews["list"].Execute(w,page)
}

func details(w http.ResponseWriter, r *http.Request) {
  domain := webl.LoadDomain(v("url",r),false)
  page := DetailsPage{ 
    HtmlHeader: &HtmlHeader{ Title: fmt.Sprintf("%s sitemap (using webl)",domain.Name) }, 
    HtmlFooter: &HtmlFooter{ AppHost: host }, 
    Domain: &domain,
  }
  webViews["details"].Execute(w, page)
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

func deleteDomain(w http.ResponseWriter, r *http.Request) {
  url := v("url",r)
  webl.DeleteDomain(url)
  fmt.Fprintf(w, url)
}

func deleteAllDomains(w http.ResponseWriter, r *http.Request) {
  webl.DeleteAllDomains()
  fmt.Fprintf(w, "true")
}


func checkDomain(w http.ResponseWriter, r *http.Request) {
  domain := webl.LoadDomain(v("url",r),false)
  if domain.LastAnalyzed != "" {
    fmt.Fprintf(w, "true")
  } else {
    fmt.Fprintf(w, "false")
  }
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

  host = fmt.Sprintf("localhost:%s", *port)

  webl.Pool = webl.NewPool(*redisServer, *redisPassword)

  m := pat.New()
  m.Get("/static/", http.HandlerFunc(getStatic))
  m.Get("/favicon.ico", http.HandlerFunc(getStatic))
  m.Get("/u/:url", http.HandlerFunc(getUrl))
  m.Get("/details/:url", http.HandlerFunc(details))
  m.Get("/", http.HandlerFunc(homepage))
  m.Get("/list", http.HandlerFunc(listSitemaps))
  m.Get("/kill", http.HandlerFunc(deleteAllDomains))
  m.Post("/delete/:url", http.HandlerFunc(deleteDomain))
  m.Get("/exists/:url", http.HandlerFunc(checkDomain))
  m.Get("/exists/", http.HandlerFunc(checkDomain))

  fmt.Println(fmt.Sprintf("Starting server, accessible at http://%s", host))
  http.Handle("/ws/crawl", websocket.Handler(doCrawl))
  http.Handle("/", m)
  err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
  fmt.Println("Server shutting down, goodbye!")
}

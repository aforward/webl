package webl

import (
  "strings"
  "fmt"
  "net/url"
  "code.google.com/p/go.net/html"
  "gopkg.in/fatih/set.v0"
  "github.com/temoto/robotstxt-go"
  "net/http"
)

//----------------
// DATA STRUCTURES
//----------------

type Resource struct {
  Name string
  LastAnalyzed string
  Url string
  Status string
  StatusCode int
  LastModified string
  Type string
  Links []*Resource
  Assets []*Resource
}

func (resource *Resource) FriendlyName() string {
  return ToFriendlyName(resource.Url)
}

func (resource *Resource) FriendlyType() string {
  return ToFriendlyType(resource.Type)
}

func (resource *Resource) FriendlyStatus() string {
  return ToFriendlyStatus(resource.Status,resource.StatusCode)
}

type Graph struct {
  Edges []Edge 
}

type Edge struct {
  FromName string
  ToName string
}

//----------------
// PUBLIC
//----------------

func Version() string {
  return "0.0.1"
}

func ToFriendlyName(raw_url string) (name string) {
  u, _ := url.Parse(raw_url)

  if (u.Path == "" || u.Path == "/") {
    name = strings.Split(u.Host, ":")[0]
  } else {
    name = u.Path
  }
  return
}

func ToFriendlyType(raw_type string) string {
  if raw_type == "" {
    return ""
  }

  all := strings.Split(raw_type, "/")
  name := all[len(all) - 1]
  switch name {
  case "x-javascript":
    return "js"
  case "msword":
    return "doc"
  case "x-shockwave-flash":
    return "flash"
  }
  return name
}

func ToFriendlyStatus(status string, code int) string {
  switch status {
  case "missing":
    return "--"
  }
  return fmt.Sprintf("%d",code)
}

func toDomain(provided string) (domain string) {
  providedUrl, _ := url.Parse(provided)
  if providedUrl.Host == "" {
    domain = provided
  } else {
    domain = strings.Split(providedUrl.Host, ":")[0]  
  }
  return
}

func IsWebpage(contentType string) (bool) {
  return contentType == "text/html" || strings.Contains(contentType,"text/html")
}

func CreateGraph(domain *Resource) Graph {
  edges := make([]Edge,10)
  edges = flattenEdges(edges,domain,set.New())
  return Graph{ Edges: edges }
}

//----------------
// HEL-]PERS
//----------------

func canRobotsAccess(input string, allRobots map[string]*robotstxt.RobotsData) (canAccess bool) {
  canAccess = true
  robotsUrl := toRobotsUrl(input)
  inputPath := toPath(input)

  if robot,ok := allRobots[robotsUrl]; ok {
    if robot == nil {
      return
    }
    canAccess = robot.TestAgent(inputPath, "WeblBot")
  } else {
    allRobots[robotsUrl] = nil
    TRACE.Println(fmt.Sprintf("Loading %s",robotsUrl))
    resp, err := http.Get(robotsUrl)
    if resp != nil && resp.Body != nil {
      defer resp.Body.Close()  
    }
    if err != nil {
      return
    }
    if resp.StatusCode != 200 {
      TRACE.Println(fmt.Sprintf("Unable to access %s, assuming full access.",robotsUrl)) 
      return
    } else {
      robot, err := robotstxt.FromResponse(resp)
      if err != nil {
        return
      }
      allRobots[robotsUrl] = robot
      canAccess = robot.TestAgent(inputPath, "WeblBot")
      TRACE.Println(fmt.Sprintf("Access to %s via %s (ok? %t)",inputPath,robotsUrl,canAccess))  
    }
  }
  return
}

func toRobotsUrl(provided string) (robotsUrl string) {
  providedUrl, _ := url.Parse(provided)
  robotsUrl = fmt.Sprintf("%s://%s/robots.txt",providedUrl.Scheme,providedUrl.Host)
  return
}

func toPath(input string) (path string) {
  inputUrl, _ := url.Parse(input)
  path = inputUrl.Path
  if path == "" {
    path = "/"  
  }
  return
}

func toUrl(provided string, path string) (processedUrl string) {

  if provided == path {
    path = ""
  }

  path_url, _ := url.Parse(path)
  providedUrl, _ := url.Parse(provided)

  if path_url.Scheme != "" {
    processedUrl = fmt.Sprintf("%s://%s",path_url.Scheme,path_url.Host)
    if path_url.Path != "/" {
      processedUrl = fmt.Sprintf("%s%s",processedUrl,path_url.Path)
    }
    return
  }

  if path == "#" || path == "/" || path == "?" {
    path = ""
  }

  processedUrl = provided
  if providedUrl.Scheme == "" {
    processedUrl = fmt.Sprintf("http://%s", processedUrl)
  }

  is_missing_protocol := strings.HasPrefix(path,"//")
  if is_missing_protocol {
    processedUrl = fmt.Sprintf("http:%s", path)
    path = ""
  }

  is_absolute_path := strings.HasPrefix(path,"/")

  if path == "" && strings.HasSuffix(processedUrl,"/") {
    processedUrl = processedUrl[0:len(processedUrl)-1]
  }


  if !strings.HasSuffix(processedUrl,"/") && path != "" && !is_absolute_path {
    processedUrl = fmt.Sprintf("%s/", processedUrl)
  } else if is_absolute_path {
    u, _ := url.Parse(processedUrl)
    processedUrl = processedUrl[0:len(processedUrl)-len(u.Path)]
  }

  processedUrl = fmt.Sprintf("%s%s", processedUrl, path)
  return 
}

func shouldProcessUrl(provided string, current string) (bool) {
  providedUrl, _ := url.Parse(toUrl(provided,""))
  currentUrl, _ := url.Parse(toUrl(current,""))
  provided_domain := strings.Split(providedUrl.Host, ":")[0]
  current_domain := strings.Split(currentUrl.Host, ":")[0]
  return provided_domain == current_domain
}

func resourcePath(token html.Token) (string) {
  for _,attr := range token.Attr {
    switch attr.Key {
    case "href","src":
      return strings.TrimSpace(attr.Val)
    }
  }
  return ""
}

func initCapacity(maxOutstanding int) (sem chan int) {
  sem = make(chan int, maxOutstanding)
  for i := 0; i < maxOutstanding; i++ {
    sem <- 1
  }
  return
}

func flattenEdges(edges []Edge, node *Resource, alreadyProcessed *set.Set) []Edge {
  // root.Links = append(root.Links,node.Links...)
  for _,link := range node.Links {
    if IsWebpage(link.Type) {
      r2l := fmt.Sprintf("%s -> %s", node.Name, link.Name) 
      l2r := fmt.Sprintf("%s -> %s", link.Name, node.Name) 

      if link.StatusCode == 200 && !alreadyProcessed.Has(r2l) && !alreadyProcessed.Has(l2r) {
        edges = append(edges,Edge{ FromName: node.Name, ToName: link.Name })
      }

      alreadyProcessed.Add(r2l)
      alreadyProcessed.Add(l2r)
      
      if !alreadyProcessed.Has(link.Url) {
        alreadyProcessed.Add(link.Url)
        edges = flattenEdges(edges,link,alreadyProcessed)
      }
    }
  } 
  return edges
}



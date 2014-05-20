package webl

import (
  "strings"
  "fmt"
  "net/url"
  "code.google.com/p/go.net/html"
  "gopkg.in/fatih/set.v0"
)

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
  provided_url, _ := url.Parse(provided)
  if provided_url.Host == "" {
    domain = provided
  } else {
    domain = strings.Split(provided_url.Host, ":")[0]  
  }
  return
}

func toUrl(provided string, path string) (processed_url string) {
  path_url, _ := url.Parse(path)
  provided_url, _ := url.Parse(provided)

  if path_url.Scheme != "" {
    processed_url = fmt.Sprintf("%s://%s",path_url.Scheme,path_url.Host)
    if path_url.Path != "/" {
      processed_url = fmt.Sprintf("%s%s",processed_url,path_url.Path)
    }
    return
  }

  if path == "#" || path == "/" || path == "?" {
    path = ""
  }

  processed_url = provided
  if provided_url.Scheme == "" {
    processed_url = fmt.Sprintf("http://%s", processed_url)
  }

  is_absolute_path := strings.HasPrefix(path,"/")

  if path == "" && strings.HasSuffix(processed_url,"/") {
    processed_url = processed_url[0:len(processed_url)-1]
  }


  if !strings.HasSuffix(processed_url,"/") && path != "" && !is_absolute_path {
    processed_url = fmt.Sprintf("%s/", processed_url)
  } else if is_absolute_path {
    u, _ := url.Parse(processed_url)
    processed_url = processed_url[0:len(processed_url)-len(u.Path)]
  }

  processed_url = fmt.Sprintf("%s%s", processed_url, path)
  return 
}

func shouldProcessUrl(provided string, current string) (bool) {
  provided_url, _ := url.Parse(toUrl(provided,""))
  currentUrl, _ := url.Parse(toUrl(current,""))
  provided_domain := strings.Split(provided_url.Host, ":")[0]
  current_domain := strings.Split(currentUrl.Host, ":")[0]
  return provided_domain == current_domain
}

func IsWebpage(contentType string) (bool) {
  return contentType == "text/html" || strings.Contains(contentType,"text/html")
}

func resourcePath(token html.Token) (string) {
  for _,attr := range token.Attr {
    switch attr.Key {
    case "href","src":
      return attr.Val
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

func CreateGraph(domain *Resource) Graph {
  edges := make([]Edge,10)
  edges = flattenEdges(edges,domain,set.New())
  return Graph{ Edges: edges }
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



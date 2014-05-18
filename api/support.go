package webl

import (
  "strings"
  "fmt"
  "net/url"
  "code.google.com/p/go.net/html"
)

func Version() string {
  return "0.0.1"
}

func toFriendlyName(raw_url string) (name string) {
  u, _ := url.Parse(raw_url)

  if (u.Path == "" || u.Path == "/") {
    name = strings.Split(u.Host, ":")[0]
  } else {
    name = u.Path
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

func isWebpage(contentType string) (bool) {
  return contentType == "text/html"
}

func resource_path(token html.Token) (string) {
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



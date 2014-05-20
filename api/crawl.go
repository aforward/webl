package webl

import (
  "fmt"
  "sync"
  "net/http"
  "code.google.com/p/go.net/html"
  "io"
  "gopkg.in/fatih/set.v0"
  "io/ioutil"
)

//----------------
// PUBLIC
//----------------

func Crawl(input string, saveDir string) bool {
  domainName := toDomain(input)
  url := toUrl(domainName,input)

  if (domainName == "") {
    WARN.Println("No domain provided, nothing to crawl.")
    return false
  }

  lastAnalyzed := friendlyNow()
  INFO.Println(fmt.Sprintf("About to crawl (%s): %s", lastAnalyzed, domainName))

  httpLimitChannel := initCapacity(4)
  var wg sync.WaitGroup
  wg.Add(1)

  alreadyProcessed := set.New()
  name := ToFriendlyName(url)

  AddDomain(&Resource{ Name: name, Url: url, LastAnalyzed: lastAnalyzed })
  go fetchResource(name, url, alreadyProcessed, httpLimitChannel, &wg)
  TRACE.Println("Wait...")
  wg.Wait();
  TRACE.Println("Done waiting.")

  savedResource := LoadDomain(domainName,true)
  WriteSitemap(savedResource, fmt.Sprintf("%s/%s.xml", saveDir,domainName))
  INFO.Println(fmt.Sprintf("Done crawing: %s", domainName))
  return true
}

//----------------
// HELPERS
//----------------

func fetchResource(domainName string, currentUrl string, alreadyProcessed *set.Set, httpLimitChannel chan int, wg *sync.WaitGroup) {
  defer wg.Done()
  if alreadyProcessed.Has(currentUrl) {
    TRACE.Println(fmt.Sprintf("Duplicate (skipping): %s", currentUrl))
  } else if shouldProcessUrl(domainName,currentUrl) {
    saveResource(&Resource{ Name: ToFriendlyName(currentUrl), Url: currentUrl })    
    alreadyProcessed.Add(currentUrl)
    TRACE.Println(fmt.Sprintf("Fetch: %s", currentUrl))

    <-httpLimitChannel
    resp, err := http.Get(currentUrl)
    httpLimitChannel <- 1

    should_close_resp := true

    if err != nil {
      WARN.Println(fmt.Sprintf("UNABLE TO FETCH %s, due to %s", currentUrl, err))
    } else {
      contentType := resp.Header.Get("Content-Type")
      lastModified := resp.Header.Get("Last-Modified")
      TRACE.Println(fmt.Sprintf("Done Fetch (%s %s): %s",contentType, resp.Status, currentUrl))
      saveResource(&Resource{ Name: ToFriendlyName(currentUrl), Url: currentUrl, Type: contentType, Status: resp.Status, StatusCode: resp.StatusCode, LastModified: lastModified })
      if IsWebpage(contentType) {
        if (!shouldProcessUrl(domainName,resp.Request.URL.String())) {
          TRACE.Println(fmt.Sprintf("Not following %s, as we redirected to a URL we should not process %s", currentUrl,resp.Request.URL.String()))
        } else if resp.StatusCode != 200 {
          WARN.Println(fmt.Sprintf("Not analyzing due to status code (%s): %s", resp.Status, currentUrl))
        } else {
          should_close_resp = false
          wg.Add(1);
          go analyzeResource(domainName, currentUrl, resp, alreadyProcessed, httpLimitChannel, wg)
        }
      }
    }
    if should_close_resp {
      if resp == nil {
        return
      }
      defer io.Copy(ioutil.Discard, resp.Body)
      if resp.Body == nil {
        return
      }
      defer resp.Body.Close()
    }
  } else {
    TRACE.Println(fmt.Sprintf("Skipping: %s", currentUrl))
  }
}

func analyzeResource(domainName string, currentUrl string, resp *http.Response, alreadyProcessed *set.Set, httpLimitChannel chan int, wg *sync.WaitGroup) {
  defer wg.Done()
  defer resp.Body.Close()
  defer io.Copy(ioutil.Discard, resp.Body)

  INFO.Println(fmt.Sprintf("Analyze (%s): %s", resp.Status, currentUrl))
  tokenizer := html.NewTokenizer(resp.Body)
  for { 
    token_type := tokenizer.Next() 
    if token_type == html.ErrorToken {
      if tokenizer.Err() != io.EOF {
        WARN.Println(fmt.Sprintf("HTML error found in %s due to ", currentUrl, tokenizer.Err()))  
      }
      return     
    }       
    token := tokenizer.Token()
    switch token_type {
    case html.StartTagToken, html.SelfClosingTagToken: // <tag>
      path := resourcePath(token)
      if path != "" {
        wg.Add(1)
        nextUrl := toUrl(domainName,path)
        saveEdge(domainName,currentUrl,nextUrl)
        go fetchResource(domainName,nextUrl,alreadyProcessed,httpLimitChannel,wg)
      }
    }
  }
  TRACE.Println(fmt.Sprintf("Done Analyze: %s", currentUrl))  
}

package webl

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
)

var (
  Pool *redis.Pool
)

type Resource struct {
  Name string
  LastAnalyzed string
  Url string
  Status string
  StatusCode int
  LastModified string
  Type string
  Links []Resource
  Assets []Resource
}

func NewPool(server, password string) *redis.Pool {
  return &redis.Pool{
    MaxIdle: 3,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) {
      c, err := redis.Dial("tcp", server)
      if err != nil {
        return nil, err
      }
      if (len(password) > 0) {
        if _, err := c.Do("AUTH", password); err != nil {
          c.Close()
          return nil, err
        }
      }
      return c, err
    },
    TestOnBorrow: func(c redis.Conn, t time.Time) error {
      _, err := c.Do("PING")
      return err
    },
  }
}

func friendlyNow() string {
  return time.Now().Format("2006-01-02 15:04:05")  
}

func addDomain(domain *Resource) {
  conn := Pool.Get()
  defer conn.Close()

  conn.Do("SADD","domains",domain.Name)
  saveResource(domain)
}

func DeleteDomain(domainName string) {
  conn := Pool.Get()
  defer conn.Close()

  conn.Do("SREM","domains",domainName)
  deleteResource(toUrl(domainName,""))
}

func DeleteAllDomains() {
  conn := Pool.Get()
  defer conn.Close()
  deleteKeys(conn,"resources:::*")
  deleteKeys(conn,"edges:::*")
  conn.Do("DEL","domains")
}

func ListDomains() (domains []*Resource) {
  conn := Pool.Get()
  defer conn.Close()

  members, err :=  redis.Strings(conn.Do("SMEMBERS","domains"))
  if members == nil {
    domains = make([]*Resource,0)
  } else {
    FailOnError(err)
    domains = make([]*Resource,len(members))
    for i, domain := range members {
      resource := LoadDomain(domain,true)
      domains[i] = &resource
    }
  }
  return
}

func saveResource(resource *Resource) {
  conn := Pool.Get()
  defer conn.Close()

  if (resource.LastAnalyzed == "") {
    resource.LastAnalyzed = friendlyNow()
  }

  key := fmt.Sprintf("resources:::%s",resource.Url)
  conn.Do("HSET",key,"name",resource.Name)
  conn.Do("HSET",key,"lastanalyzed",resource.LastAnalyzed)
  conn.Do("HSET",key,"url",resource.Url)
  conn.Do("HSET",key,"status",resource.Status)
  conn.Do("HSET",key,"statuscode",resource.StatusCode)
  conn.Do("HSET",key,"lastmodified",resource.LastModified)
  conn.Do("HSET",key,"type",resource.Type)
}

func saveEdge(domainName string, fromUrl string, toUrl string) {
  conn := Pool.Get()
  defer conn.Close()
  TRACE.Println(fmt.Sprintf("Saving edge between %s --> %s", fromUrl, toUrl))
  conn.Do("SADD",fmt.Sprintf("edges:::%s",fromUrl),toUrl)
}

func deleteResource(url string) {
  conn := Pool.Get()
  defer conn.Close()
  deleteKeys(conn, fmt.Sprintf("edges:::%s*",url))
  deleteKeys(conn, fmt.Sprintf("resources:::%s*",url))
}

func deleteKeys(conn redis.Conn, keyFilter string) {
  keys, err := redis.Strings(conn.Do("KEYS", keyFilter))
  FailOnError(err)
  for _, k := range keys {
    conn.Do("DEL",k)
  }
}

func LoadDomain(domain string, isBasic bool) Resource {
  allResources := make(map[string]Resource)
  return loadResource(toUrl(domain,""),isBasic,allResources)
}

func loadResource(url string, isBasic bool, allResources map[string]Resource) (resource Resource) {
  conn := Pool.Get()
  defer conn.Close()

  key := fmt.Sprintf("resources:::%s",url)
  if ok, _ := redis.Bool(conn.Do("EXISTS",key)); ok {
    var r struct {
      Name string           `redis:"name"`
      LastAnalyzed string   `redis:"lastanalyzed"`
      Url string            `redis:"url"`
      Status string         `redis:"status"`
      StatusCode int        `redis:"statuscode"`
      LastModified string   `redis:"lastmodified"`
      Type string           `redis:"type"`
    }
    values, err := redis.Values(conn.Do("HGETALL", key))
    FailOnError(err)

    err = redis.ScanStruct(values, &r)
    FailOnError(err)

    resource = Resource{ Name: r.Name, LastAnalyzed: r.LastAnalyzed, Url: r.Url, Status: r.Status, StatusCode: r.StatusCode, Type: r.Type }

    allResources[r.Url] = resource

    linksKey := fmt.Sprintf("edges:::%s",r.Url)
    members, err :=  redis.Strings(conn.Do("SMEMBERS",linksKey))
    resource.Links = make([]Resource,len(members))
    var possibleLink Resource
    for _, link := range members {
      if linkResource,ok := allResources[link]; ok {
        TRACE.Println(fmt.Sprintf("Reused link between %s --> %s", r.Url, link))
        possibleLink = linkResource
      } else {
        TRACE.Println(fmt.Sprintf("Looking up link between %s --> %s", r.Url, link))
        possibleLink = loadResource(link,isBasic,allResources)
        allResources[possibleLink.Url] = possibleLink
      }
      if (possibleLink.Url == "") {
        continue
      }
      if (!isBasic || isWebpage(possibleLink.Type)) {
        resource.Links = append(resource.Links, possibleLink)
      }
    }
  } else {
    resource = Resource{ Name: toFriendlyName(url), Url: url, Status: "missing" }
  }
  return
}




package webl

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
)

var (
  Pool *redis.Pool
)

//----------------
// PUBLIC
//----------------

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

func AddDomain(domain *Resource) {
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
      domains[i] = LoadDomain(domain,true)
    }
  }
  return
}

func LoadDomain(domain string, isBasic bool) *Resource {
  resource := LoadResource(toUrl(domain,""), isBasic)
  resource.Name = domain
  return resource
}

func LoadResource(domain string, isBasic bool) *Resource {
  return findResource(toUrl(domain,""), true, isBasic, make(map[string]*Resource))
}

//----------------
// HELPERS
//----------------

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
  TRACE.Println(fmt.Sprintf("Saving edge between %s -> %s", fromUrl, toUrl))
  conn.Do("SADD",fmt.Sprintf("edges:::%s",fromUrl),toUrl)
}

func deleteResource(url string) {
  conn := Pool.Get()
  defer conn.Close()
  deleteKeys(conn, fmt.Sprintf("edges:::%s*",url))
  deleteKeys(conn, fmt.Sprintf("resources:::%s*",url))
}

func findResource(url string, isRoot bool, isBasic bool, allResources map[string]*Resource) (resource *Resource) {
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

    resource = &Resource{ Name: r.Name, LastAnalyzed: r.LastAnalyzed, Url: r.Url, Status: r.Status, StatusCode: r.StatusCode, Type: r.Type }
    allResources[r.Url] = resource
    if !isRoot && isBasic {
      return
    }

    members, err :=  redis.Strings(conn.Do("SMEMBERS",fmt.Sprintf("edges:::%s",r.Url)))
    resource.Links = make([]*Resource,0)
    resource.Assets = make([]*Resource,0)
    var possibleLink *Resource
    for _, link := range members {
      if linkResource,ok := allResources[link]; ok {
        TRACE.Println(fmt.Sprintf("Reused link between %s -> %s", r.Url, link))
        possibleLink = linkResource
      } else {
        TRACE.Println(fmt.Sprintf("Looking up link between %s -> %s", r.Url, link))
        possibleLink = findResource(link,false,isBasic,allResources)
        allResources[possibleLink.Url] = possibleLink
      }
      if (possibleLink.Url == "") {
        continue
      }

      if IsWebpage(possibleLink.Type) || possibleLink.Type == "" {
        resource.Links = append(resource.Links, possibleLink)
      } else {
        resource.Assets = append(resource.Assets, possibleLink)
      }

    }
  } else {
    resource = &Resource{ Name: ToFriendlyName(url), Url: url, Status: "missing" }
  }
  return
}

func deleteKeys(conn redis.Conn, keyFilter string) {
  keys, err := redis.Strings(conn.Do("KEYS", keyFilter))
  FailOnError(err)
  for _, k := range keys {
    conn.Do("DEL",k)
  }
}

func friendlyNow() string {
  return time.Now().Format("2006-01-02 15:04:05")  
}



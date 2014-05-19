package webl

import (
  "encoding/xml"
  "io/ioutil"
  "os"
  "path"
  "fmt"
  "gopkg.in/fatih/set.v0"
)

//----------------
// STRUCTURES
//----------------

type Sitemap struct {
  XMLName        xml.Name     `xml:"urlset"`
  Namespace      string       `xml:"xmlns,attr"`
  Schema         string       `xml:"xmlns:xsi,attr"`
  SchemaLocation string       `xml:"xsi:schemaLocation,attr"`
  Urls           []UrlItem
}

type UrlItem struct {
  XMLName    xml.Name `xml:"url"`
  Loc        string   `xml:"loc"`
  StatusCode int      `xml:"-"`
  LastMod    string   `xml:"lastmod,omitempty"`
  Priority   float32  `xml:"priority,omitempty"`
  ChangeFreq string   `xml:"changefreq,omitempty"`
}

func (item *UrlItem) FriendlyName() string {
  return ToFriendlyName(item.Loc)
}

func (item *UrlItem) Assets() []*Resource {
  resource := LoadResource(item.Loc,true)
  if resource.Assets == nil {
    return make([]*Resource,0)
  } else {
    return resource.Assets  
  }
}

func (item *UrlItem) Links() []*Resource {
  resource := LoadResource(item.Loc,true)
  if resource.Links == nil {
    return make([]*Resource,0)
  } else {
    return resource.Links  
  }
}

//----------------
// PUBLIC
//----------------

func InitSitemap() Sitemap {
  return Sitemap{ 
      Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9",
      Schema: "http://www.w3.org/2001/XMLSchema-instance",
      SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
    }
}

func GenerateSitemap(domain *Resource, onlyValidUrls bool) *Sitemap {
  TRACE.Println(fmt.Sprintf("Generating sitemap for: %s", domain.Name))
  sitemap := InitSitemap()
  alreadyProcessed := set.New()
  alreadyProcessed.Add(domain.Url)
  sitemap.Urls = []UrlItem{UrlItem{Loc: domain.Url, StatusCode: domain.StatusCode, LastMod: domain.LastModified }}
  appendSitemapChildren(&sitemap,domain,onlyValidUrls,alreadyProcessed)
  return &sitemap
}

func WriteSitemap(domain *Resource, filename string) (sitemap *Sitemap) {
  sitemap = GenerateSitemap(domain,true)
  TRACE.Println(fmt.Sprintf("Saving sitemap to: %s", filename))
  urlSetAsXml, err := xml.MarshalIndent(sitemap, "", "  ")
  FailOnError(err)
  os.MkdirAll(path.Dir(filename), 0755)
  err = ioutil.WriteFile(filename, urlSetAsXml, 0744)
  FailOnError(err)
  return
}

//----------------
// HELPERS
//----------------

func appendSitemapChildren(urlSet *Sitemap, resource *Resource, onlyValidUrls bool, alreadyProcessed *set.Set) {
  for _,link := range resource.Links {
    if alreadyProcessed.Has(link.Url) {
      continue
    }
    alreadyProcessed.Add(link.Url)
    if (link.Url == "") {
      continue
    } else if (link.Status == "missing") {
      TRACE.Println(fmt.Sprintf("Skipping invalid link: %s", link.Url))
    } else if (!IsWebpage(link.Type)) {
      TRACE.Println(fmt.Sprintf("Skipping non-webpage (%s %s): %s", link.Type, link.Status, link.Url))
    } else if (onlyValidUrls && (link.StatusCode == 404 || link.StatusCode >= 500)) {
      TRACE.Println(fmt.Sprintf("Skipping invalid resource (%s %s): %s", link.Type, link.Status, link.Url))
    } else {
      TRACE.Println(fmt.Sprintf("Adding to sitemap (%s %s): %s", link.Type, link.Status, link.Url))
      urlSet.Urls = append(urlSet.Urls, UrlItem{Loc: link.Url, StatusCode: link.StatusCode, LastMod: link.LastModified})
      appendSitemapChildren(urlSet,link,onlyValidUrls,alreadyProcessed)
    }
  }
}

package webl

import (
  "encoding/xml"
  "io/ioutil"
  "os"
  "path"
  "fmt"
  "gopkg.in/fatih/set.v0"
)

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
  StatusCode int
  LastMod    string   `xml:"lastmod,omitempty"`
  Priority   float32  `xml:"priority,omitempty"`
  ChangeFreq string   `xml:"changefreq,omitempty"`
}

func (item *UrlItem) FriendlyName() string {
  return ToFriendlyName(item.Loc)
}

func (item *UrlItem) Assets() []Resource {
  resource := LoadResource(item.Loc,true)
  return resource.Assets
}

func (item *UrlItem) Links() []Resource {
  resource := LoadResource(item.Loc,true)
  return resource.Links
}


func WriteSitemap(domain *Resource, filename string) (urlSet *Sitemap) {
  urlSet = GenerateSitemap(domain,true)
  TRACE.Println(fmt.Sprintf("Saving sitemap to: %s", filename))
  urlSetAsXml, err := xml.MarshalIndent(urlSet, "", "  ")
  FailOnError(err)
  os.MkdirAll(path.Dir(filename), 0755)
  err = ioutil.WriteFile(filename, urlSetAsXml, 0744)
  FailOnError(err)
  return
}

func GenerateSitemap(domain *Resource, onlyValidUrls bool) (urlSet *Sitemap) {
  TRACE.Println(fmt.Sprintf("Generating sitemap for: %s", domain.Name))
  urlSet = initSitemap()
  alreadyProcessed := set.New()
  alreadyProcessed.Add(domain.Url)
  urlSet.Urls = []UrlItem{UrlItem{Loc: domain.Url, StatusCode: domain.StatusCode, LastMod: domain.LastModified }}
  appendSitemapChildren(urlSet,domain,onlyValidUrls,alreadyProcessed)
  return
}

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
      appendSitemapChildren(urlSet,&link,onlyValidUrls,alreadyProcessed)
    }
  }
}

func initSitemap() *Sitemap {
  return &Sitemap{ 
      Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9",
      Schema: "http://www.w3.org/2001/XMLSchema-instance",
      SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
    }
}
package webl

import (
  "encoding/xml"
  "io/ioutil"
  "os"
  "path"
  "fmt"
  "gopkg.in/fatih/set.v0"
)

type UrlSet struct {
  XMLName        xml.Name     `xml:"urlset"`
  Namespace      string       `xml:"xmlns,attr"`
  Schema         string       `xml:"xmlns:xsi,attr"`
  SchemaLocation string       `xml:"xsi:schemaLocation,attr"`
  Urls           []UrlItem
}

type UrlItem struct {
  XMLName    xml.Name `xml:"url"`
  Loc        string   `xml:"loc"`

  LastMod    string   `xml:"lastmod,omitempty"`
  Priority   float32  `xml:"priority,omitempty"`
  ChangeFreq string   `xml:"changefreq,omitempty"`
}

func writeSitemap(domain *Resource, filename string) (urlSet *UrlSet) {
  urlSet = generateSitemap(domain)
  TRACE.Println(fmt.Sprintf("Saving sitemap to: %s", filename))
  urlSetAsXml, err := xml.MarshalIndent(urlSet, "", "  ")
  FailOnError(err)
  os.MkdirAll(path.Dir(filename), 0755)
  err = ioutil.WriteFile(filename, urlSetAsXml, 0744)
  FailOnError(err)
  return
}

func generateSitemap(domain *Resource) (urlSet *UrlSet) {
  TRACE.Println(fmt.Sprintf("Generating sitemap for: %s", domain.Name))
  urlSet = initUrlSet()
  alreadyProcessed := set.New()
  urlSet.Urls = []UrlItem{UrlItem{Loc: domain.Url, LastMod: domain.LastModified }}
  appendSitemapChildren(urlSet,domain,alreadyProcessed)
  return
}

func appendSitemapChildren(urlSet *UrlSet, resource *Resource, alreadyProcessed *set.Set) {
  for _,link := range resource.Links {
    if alreadyProcessed.Has(link.Url) {
      continue
    }
    alreadyProcessed.Add(link.Url)
    if (link.Url == "") {
      continue
    } else if (link.Status == "missing") {
      TRACE.Println(fmt.Sprintf("Skipping invalid link: %s", link.Url))
    } else if (!isWebpage(link.Type)) {
      TRACE.Println(fmt.Sprintf("Skipping non-webpage (%s %s): %s", link.Type, link.Status, link.Url))
    } else if (link.StatusCode == 404 || link.StatusCode >= 500) {
      TRACE.Println(fmt.Sprintf("Skipping invalid resource (%s %s): %s", link.Type, link.Status, link.Url))
    } else {
      TRACE.Println(fmt.Sprintf("Adding to sitemap (%s %s): %s", link.Type, link.Status, link.Url))
      urlSet.Urls = append(urlSet.Urls, UrlItem{Loc: link.Url, LastMod: link.LastModified})
      appendSitemapChildren(urlSet,&link,alreadyProcessed)
    }
  }
}

func initUrlSet() *UrlSet {
  return &UrlSet{ 
      Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9",
      Schema: "http://www.w3.org/2001/XMLSchema-instance",
      SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
    }
}
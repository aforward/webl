package webl

import (
  . "gopkg.in/check.v1"
  "io/ioutil"
  "os"
)

//------
// InitSitemap
//------

func (s *MySuite) Test_InitSitemap(c *C) {
  sitemap := InitSitemap()

  c.Check(sitemap.Namespace,Equals,"http://www.sitemaps.org/schemas/sitemap/0.9")
  c.Check(sitemap.Schema,Equals,"http://www.w3.org/2001/XMLSchema-instance")
  c.Check(sitemap.SchemaLocation,Equals,"http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd")
}

//------
// FriendlyName
//------

func (s *MySuite) Test_FriendlyName(c *C) {
  r := Resource{ Name: "a4word.com", Url: "http://a4word.com/one.php" }
  sitemap := GenerateSitemap(&r,true)
  c.Check(sitemap.Urls[0].FriendlyName(),Equals,"/one.php")
}

//------
// Assets / Links
//------

func (s *MySuite) Test_AssetsLinks_empty(c *C) {
  Pool = NewPool(":6379","")
  DeleteDomain("a4word.com")

  r := Resource{ Name: "a4word.com", Url: "http://a4word.com/" }
  sitemap := GenerateSitemap(&r,true)
  c.Check(len(sitemap.Urls[0].Assets()),Equals,0)
  c.Check(len(sitemap.Urls[0].Links()),Equals,0)
}

func (s *MySuite) Test_Assets_some(c *C) {
  Pool = NewPool(":6379","")
  DeleteDomain("a4word.com")

  r := &Resource{ Name: "a4word.com", Url: "http://a4word.com", Type: "text/html" }
  child1 := &Resource{ Name: "http://a4word.com/one.php", Url: "http://a4word.com/one.php", Type: "text/html" }
  child2 := &Resource{ Name: "http://a4word.com/two.js", Url: "http://a4word.com/two.js", Type: "application/js" }
  r.Links = []*Resource{child1,child2}
  saveResource(r)
  saveResource(child1)
  saveResource(child2)
  saveEdge("a4word.com", "http://a4word.com", "http://a4word.com/one.php")
  saveEdge("a4word.com", "http://a4word.com", "http://a4word.com/two.js")

  sitemap := GenerateSitemap(r,true)
  c.Check(len(sitemap.Urls[0].Links()),Equals,1)
  c.Check(len(sitemap.Urls[0].Assets()),Equals,1)

  c.Check(sitemap.Urls[0].Links()[0].Name,Equals,child1.Name)
  c.Check(sitemap.Urls[0].Assets()[0].Name,Equals,child2.Name)
}

//------
// GenerateSitemap
//------

func (s *MySuite) Test_GenerateSitemap_empty(c *C) {
  r := Resource{ Name: "a4word.com", Url: "http://a4word.com/" }
  sitemap := GenerateSitemap(&r,true)
  c.Check(sitemap.Urls[0],Equals,UrlItem{Loc: "http://a4word.com/"})
}

func (s *MySuite) Test_GenerateSitemap_multiple(c *C) {
  r := &Resource{ Name: "a4word.com", Url: "http://a4word.com/", Type: "text/html" }
  child := &Resource{ Name: "http://a4word.com/one.php", Url: "http://a4word.com/one.php", Type: "text/html" }
  grandchild := &Resource{ Name: "http://a4word.com/two.php", Url: "http://a4word.com/two.php", Type: "text/html" }

  child.Links = []*Resource{grandchild}
  r.Links = []*Resource{child}

  sitemap := GenerateSitemap(r,false)

  c.Check(len(sitemap.Urls),Equals,3)
  c.Check(sitemap.Urls[0],Equals,UrlItem{Loc: "http://a4word.com/"})
  c.Check(sitemap.Urls[1],Equals,UrlItem{Loc: "http://a4word.com/one.php"})
  c.Check(sitemap.Urls[2],Equals,UrlItem{Loc: "http://a4word.com/two.php"})
}

func (s *MySuite) Test_GenerateSitemap_ignoreInvalid(c *C) {
  r := &Resource{ Name: "a4word.com", Url: "http://a4word.com/", Type: "text/html", StatusCode: 200 }
  child := &Resource{ Name: "http://a4word.com/one.php", Url: "http://a4word.com/one.php", Type: "text/html", StatusCode: 200 }
  grandchild1 := &Resource{ Name: "http://a4word.com/two.php", Url: "http://a4word.com/two.php", Type: "text/html", StatusCode: 500 }
  grandchild2 := &Resource{ Name: "http://a4word.com/three.php", Url: "http://a4word.com/three.php", Type: "application/js", StatusCode: 200 }

  child.Links = []*Resource{grandchild1,grandchild2}
  r.Links = []*Resource{child}

  sitemap := GenerateSitemap(r,true)

  c.Check(len(sitemap.Urls),Equals,2)
  c.Check(sitemap.Urls[0],Equals,UrlItem{Loc: "http://a4word.com/", StatusCode: 200})
  c.Check(sitemap.Urls[1],Equals,UrlItem{Loc: "http://a4word.com/one.php", StatusCode: 200})
}

//------
// WriteSitemap
//------

func (s *MySuite) Test_WriteSitemap_empty(c *C) {
  os.Remove("./tmp/test_sitemap.xml")

  r := Resource{ Name: "a4word.com", StatusCode: 400, Url: "http://a4word.com/" }
  WriteSitemap(&r,"./tmp/test_sitemap.xml")

  expectedData, _ := ioutil.ReadFile("./sampledata/test_sitemap.xml")
  savedData, _ := ioutil.ReadFile("./tmp/test_sitemap.xml")
  c.Check(string(savedData),Equals,string(expectedData))
}

func (s *MySuite) Test_WriteSitemap_multiple(c *C) {
  r := &Resource{ Name: "a4word.com", Url: "http://a4word.com/", Type: "text/html" }
  child := &Resource{ Name: "http://a4word.com/one.php", Url: "http://a4word.com/one.php", Type: "text/html" }
  grandchild := &Resource{ Name: "http://a4word.com/two.php", Url: "http://a4word.com/two.php", Type: "text/html" }

  child.Links = []*Resource{grandchild}
  r.Links = []*Resource{child}

  os.Remove("./tmp/test_sitemap2.xml")
  WriteSitemap(r,"./tmp/test_sitemap2.xml")

  expectedData, _ := ioutil.ReadFile("./sampledata/test_sitemap2.xml")
  savedData, _ := ioutil.ReadFile("./tmp/test_sitemap2.xml")
  c.Check(string(savedData),Equals,string(expectedData))
}
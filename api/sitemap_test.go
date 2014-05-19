package webl

import (
  . "gopkg.in/check.v1"
  // "io/ioutil"
  // "os"
)

//------
// version
//------

func (s *MySuite) Test_GenerateSitemap_empty(c *C) {
  r := Resource{ Name: "a4word.com", Url: "http://a4word.com/" }
  urlSet := GenerateSitemap(&r,true)
  c.Check(urlSet.Urls[0],Equals,UrlItem{Loc: "http://a4word.com/"})
}

// TODO: fix me
// func (s *MySuite) Test_GenerateSitemap_multiple(c *C) {
//   r := Resource{ Name: "a4word.com", Url: "http://a4word.com/" }
//   child := Resource{ Name: "http://a4word.com/one.php", Url: "http://a4word.com/one.php" }
//   grandchild := Resource{ Name: "http://a4word.com/two.php", Url: "http://a4word.com/two.php" }

//   r.Links = []Resource{child}
//   child.Links = []Resource{grandchild}

//   urlSet := GenerateSitemap(r)
//   c.Check(urlSet.Urls[0],Equals,UrlItem{Loc: "http://a4word.com/"})
//   c.Check(urlSet.Urls[1],Equals,UrlItem{Loc: "http://a4word.com/one.php"})
//   c.Check(urlSet.Urls[2],Equals,UrlItem{Loc: "http://a4word.com/two.php"})
// }

// func (s *MySuite) Test_WriteSitemap_empty(c *C) {
//   os.Remove("./tmp/a4word.com.sitemap.xml")

//   r := Resource{ Name: "a4word.com", Url: "http://a4word.com/" }
//   WriteSitemap(r,"./tmp/a4word.com.sitemap.xml")

//   expectedData, _ := ioutil.ReadFile("./docs/test_output_simple.xml")
//   savedData, _ := ioutil.ReadFile("./tmp/a4word.com.sitemap.xml")
//   c.Check(string(savedData),Equals,string(expectedData))
// }


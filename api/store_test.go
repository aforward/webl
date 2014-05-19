package webl

import (
  // "testing"
  . "gopkg.in/check.v1"
)

//------
// version
//------

func (s *MySuite) Test_saveResource(c *C) {
  Pool = NewPool(":6379","")

  r := Resource{ Name: "a", Url: "http://a", Status: "404", Type: "html" }
  saveResource(&r)
  sameR := LoadDomain("a",true)

  c.Check(sameR.Name,Equals,"a")
  c.Check(sameR.Url,Equals,"http://a")
  c.Check(sameR.Status,Equals,"404")
  c.Check(sameR.Type,Equals,"html")  
}

func (s *MySuite) Test_saveEdge(c *C) {
  Pool = NewPool(":6379","")

  saveResource(&Resource{ Name: "a4word.com", Url: "http://a4word.com", Status: "200", Type: "text/html" })
  saveResource(&Resource{ Name: "/links.php", Url: "http://a4word.com/links.php", Status: "200", Type: "text/html" })

  saveEdge("a4word.com","http://a4word.com","http://a4word.com/links.php")

  sameR := LoadDomain("a4word.com",true)

  c.Check(sameR.Name,Equals,"a4word.com")
  c.Check(sameR.Url,Equals,"http://a4word.com")
  c.Check(sameR.Status,Equals,"200")
  c.Check(sameR.Type,Equals,"text/html")  
}

func (s *MySuite) Test_deleteResource(c *C) {
  Pool = NewPool(":6379","")

  r := Resource{ Name: "aa", Url: "b", Status: "404", Type: "html" }
  saveResource(&r)
  deleteResource("aa")

  sameR := LoadDomain("aa",true)
  c.Check(sameR.Name,Equals,"aa")
  c.Check(sameR.Url,Equals,"http://aa")
  c.Check(sameR.Status,Equals,"missing")
  c.Check(sameR.Type,Equals,"")  
}

func (s *MySuite) Test_AddDomain(c *C) {
  Pool = NewPool(":6379","")
  DeleteAllDomains()
  c.Check(len(ListDomains()),Equals,0)

  r := Resource{ Name: "a4word.com", Url: "http://a4word.com" }
  AddDomain(&r)
  all := ListDomains()
  c.Check(len(all),Equals,1)
  c.Check(all[0].Name,Equals,"a4word.com")
}

func (s *MySuite) Test_RemoveDomain(c *C) {
  Pool = NewPool(":6379","")
  DeleteAllDomains()
  c.Check(len(ListDomains()),Equals,0)

  r := Resource{ Name: "a4word.com", Url: "http://a4word.com" }
  AddDomain(&r)
  all := ListDomains()
  c.Check(len(all),Equals,1)
  
  DeleteDomain("a4word.com")
  c.Check(len(ListDomains()),Equals,0)
}

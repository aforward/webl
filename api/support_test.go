package webl

import (
  "testing"
  . "gopkg.in/check.v1"
  "code.google.com/p/go.net/html"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { 
  InitLogging(false, false, false, nil)
  TestingT(t) 
}

type MySuite struct{}
var _ = Suite(&MySuite{})

//------
// version
//------

func (s *MySuite) Test_version_should_be_set(c *C) {
  c.Check(Version(),Equals,"0.0.1")
}

//------
// toFriendlyName
//------

func (s *MySuite) Test_toFriendlyName_domainName(c *C) {
  c.Check(toFriendlyName("http://a4word.com"),Equals,"a4word.com")
}

func (s *MySuite) Test_toFriendlyName_domainName2(c *C) {
  c.Check(toFriendlyName("http://a4word.com/"),Equals,"a4word.com")
}

func (s *MySuite) Test_toFriendlyName_path(c *C) {
  c.Check(toFriendlyName("http://a4word.com/a"),Equals,"/a")
}

func (s *MySuite) Test_toFriendlyName_path2(c *C) {
  c.Check(toFriendlyName("http://a4word.com/a/b/c.txt"),Equals,"/a/b/c.txt")
}

//------
// toUrl
//------

func (s *MySuite) Test_toUrl_as_is(c *C) {
  c.Check(toUrl("http://a4word.com",""),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_drop_slash(c *C) {
  c.Check(toUrl("http://a4word.com/",""),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_drop_hashtag(c *C) {
  c.Check(toUrl("http://a4word.com","http://a4word.com/#"),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_drop_question_mark(c *C) {
  c.Check(toUrl("http://a4word.com","http://a4word.com/?"),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_add_http(c *C) {
  c.Check(toUrl("a4word.com",""),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_append_path(c *C) {
  c.Check(toUrl("a4word.com","/abc"),Equals,"http://a4word.com/abc")
}

func (s *MySuite) Test_toUrl_hashtag(c *C) {
  c.Check(toUrl("a4word.com","#"),Equals,"http://a4word.com")
}

func (s *MySuite) Test_toUrl_no_slash(c *C) {
  c.Check(toUrl("a4word.com","x"),Equals,"http://a4word.com/x")
}

func (s *MySuite) Test_toUrl_no_slash_sub_dir(c *C) {
  c.Check(toUrl("a4word.com/b","x"),Equals,"http://a4word.com/b/x")
}

func (s *MySuite) Test_toUrl_slash_sub_dir(c *C) {
  c.Check(toUrl("a4word.com/b","/x"),Equals,"http://a4word.com/x")
}

func (s *MySuite) Test_toUrl_somewhere_else(c *C) {
  c.Check(toUrl("a4word.com/b","http://a5word.com/a"),Equals,"http://a5word.com/a")
}

func (s *MySuite) Test_toUrl_https_url(c *C) {
  c.Check(toUrl("https://a4word.com","x"),Equals,"https://a4word.com/x")
}

func (s *MySuite) Test_toUrl_ftp_url(c *C) {
  c.Check(toUrl("ftp://a4word.com","x"),Equals,"ftp://a4word.com/x")
}

func (s *MySuite) Test_toUrl_https_path(c *C) {
  c.Check(toUrl("a4word.com","https://a5word.com"),Equals,"https://a5word.com")
}

//------
// IsWebpage
//------

func (s *MySuite) Test_IsWebpage_html(c *C) {
  c.Check(IsWebpage("text/html"),Equals,true)
  c.Check(IsWebpage("text/garble"),Equals,false)
}

//------
// resource_path
//------

func (s *MySuite) Test_should_resource_path_empty(c *C) {
  c.Check(resource_path(html.Token{}),Equals,"")
}


func (s *MySuite) Test_should_resource_path_href(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "href", "/x"})
  c.Check(resource_path(t),Equals,"/x")
}

func (s *MySuite) Test_should_resource_path_src(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "src", "/y"})
  c.Check(resource_path(t),Equals,"/y")
}

func (s *MySuite) Test_should_resource_path_garble(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "garble", "/y"})
  c.Check(resource_path(t),Equals,"")
}


//------
// should_process
//------

func (s *MySuite) Test_should_process_no_protocol(c *C) {
  c.Check(shouldProcessUrl("a4word.com","a4word.com"),Equals,true)
  c.Check(shouldProcessUrl("a4word.com","blah.com"),Equals,false)
}

func (s *MySuite) Test_should_process_has_protocol(c *C) {
  c.Check(shouldProcessUrl("http://a4word.com","http://a4word.com"),Equals,true)
  c.Check(shouldProcessUrl("http://a4word.com","http://blah.com"),Equals,false)
}

func (s *MySuite) Test_should_process_mismatched_protocol(c *C) {
  c.Check(shouldProcessUrl("https://a4word.com","http://a4word.com"),Equals,true)
  c.Check(shouldProcessUrl("a4word.com","http://a4word.com"),Equals,true)
  c.Check(shouldProcessUrl("http://a4word.com","ftp://a4word.com"),Equals,true)
}

func (s *MySuite) Test_should_process_exclude_subdomains(c *C) {
  c.Check(shouldProcessUrl("a4word.com","http://notokay.a4word.com"),Equals,false)
}

func (s *MySuite) Test_should_process_include_subdomains_if_at_that_level(c *C) {
  c.Check(shouldProcessUrl("aha.a4word.com","http://aha.a4word.com"),Equals,true)
  c.Check(shouldProcessUrl("aha.a4word.com","http://sub.aha.a4word.com"),Equals,false)
}


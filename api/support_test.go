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
// ToFriendlyName
//------

func (s *MySuite) Test_ToFriendlyName_domainName(c *C) {
  c.Check(ToFriendlyName("http://a4word.com"),Equals,"a4word.com")
}

func (s *MySuite) Test_ToFriendlyName_domainName2(c *C) {
  c.Check(ToFriendlyName("http://a4word.com/"),Equals,"a4word.com")
}

func (s *MySuite) Test_ToFriendlyName_path(c *C) {
  c.Check(ToFriendlyName("http://a4word.com/a"),Equals,"/a")
}

func (s *MySuite) Test_ToFriendlyName_path2(c *C) {
  c.Check(ToFriendlyName("http://a4word.com/a/b/c.txt"),Equals,"/a/b/c.txt")
}

//-------
// ToFriendlyType
//-------

func (s *MySuite) Test_ToFriendlyType_empty(c *C) {
  c.Check(ToFriendlyType(""),Equals,"")
}

func (s *MySuite) Test_ToFriendlyType_lastSlash(c *C) {
  c.Check(ToFriendlyType("a/b/c"),Equals,"c")
}

func (s *MySuite) Test_ToFriendlyType_examples(c *C) {
  c.Check(ToFriendlyType("application/x-javascript"),Equals,"js")
  c.Check(ToFriendlyType("text/css"),Equals,"css")
  c.Check(ToFriendlyType("image/png"),Equals,"png")
  c.Check(ToFriendlyType("application/msword"),Equals,"doc")
  c.Check(ToFriendlyType("application/x-shockwave-flash"),Equals,"flash")
}

//-----
// toDomain
//-----

func (s *MySuite) Test_toDomain_empty(c *C) {
  c.Check(toDomain(""),Equals,"")
}

func (s *MySuite) Test_toDomain_asIs(c *C) {
  c.Check(toDomain("a4word.com"),Equals,"a4word.com")
  c.Check(toDomain("a4word"),Equals,"a4word")
  c.Check(toDomain("www.a4word.com"),Equals,"www.a4word.com")
  c.Check(toDomain("sub.a4word.com"),Equals,"sub.a4word.com")
}

func (s *MySuite) Test_toDomain_stripUrlComponents(c *C) {
  c.Check(toDomain("http://a4word.com"),Equals,"a4word.com")
  c.Check(toDomain("https://a4word.com/"),Equals,"a4word.com")
  c.Check(toDomain("git://a4word.com/"),Equals,"a4word.com")

  c.Check(toDomain("http://sub.a4word.com/"),Equals,"sub.a4word.com")
  c.Check(toDomain("http://www.a4word.com/"),Equals,"www.a4word.com")
  c.Check(toDomain("http://www.a4word.com/one/two.php"),Equals,"www.a4word.com")
}

//------
// toUrl
//------

func (s *MySuite) Test_toUrl_rootDomain(c *C) {
  c.Check(toUrl("http://a4word.com",""),Equals,"http://a4word.com")
  c.Check(toUrl("http://a4word.com","/"),Equals,"http://a4word.com")
  c.Check(toUrl("http://a4word.com","#"),Equals,"http://a4word.com")
  c.Check(toUrl("http://a4word.com","?"),Equals,"http://a4word.com")
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

func (s *MySuite) Test_IsWebpage_htmlWithCharset(c *C) {
  c.Check(IsWebpage("text/html; charset=utf-8"),Equals,true)
  c.Check(IsWebpage("charset=utf-8; text/html"),Equals,true)
  c.Check(IsWebpage("charset=utf-8; text/html garble"),Equals,true)
}

//------
// resourcePath
//------

func (s *MySuite) Test_should_resourcePath_empty(c *C) {
  c.Check(resourcePath(html.Token{}),Equals,"")
}


func (s *MySuite) Test_should_resourcePath_href(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "href", "/x"})
  c.Check(resourcePath(t),Equals,"/x")
}

func (s *MySuite) Test_should_resourcePath_src(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "src", "/y"})
  c.Check(resourcePath(t),Equals,"/y")
}

func (s *MySuite) Test_should_resourcePath_garble(c *C) {
  t := html.Token{}
  t.Attr = append(t.Attr,html.Attribute{"", "garble", "/y"})
  c.Check(resourcePath(t),Equals,"")
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


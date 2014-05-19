webl
======

A web crawler written in [Go](http://golang.org/).

Installation
------------

The webcrawler uses [Redis](http://redis.io), to store results.  Please install it and ensure it is running before starting.

Install command line tool, weblconsole, and the web server, weblui, using the "go get" command:

    go get github.com/aforward/webl/weblconsole
    go get github.com/aforward/webl/weblui

The the installation
  
    cd $GOPATH/src
    go test github.com/aforward/webl/api

Now install it

    cd $GOPATH/src
    go install github.com/aforward/webl/weblconsole
    go install github.com/aforward/webl/weblui
    
You should see an application in your bin directory

    ls -la $GOPATH/bin | grep webl

Command Line
------------

From the command line you can 

    weblconsole -url=a4word.com -verbose

You can start the webserver using

    cd $GOPATH/src/github.com/aforward/webl/weblui
    weblui

You can now crawl websites through the UI

TODO
------------

* Strengthen domain format inputs (e.g. http:// and trailing /)
* Resolve "//" as the same protocol as calling page
* Enable viewing of sitemap.xml
* Display graph link
* Read robots.txt and ensure following the directives
* Review code, extract additional tests
* Create docker deployment

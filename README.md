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

* Review code, extract additional tests
* Create docker deployment
* Adding the ability to manage multiple crawls over a domain and provide a diff of the results.
* Adding security to prevent abuse from crawling too often.
* Improve visualization based on how best to use the data (e.g. broken links, unused assets, etc). This will most likely involve an improved data store (like Postgres) to allow for reaching searching.
* Improved sitemap.xml generation to grab other fields like priority, last modified, etc.
* Improved resource meta-data like title, and keywords, as well as taking thumbnails of the webpage.
* Improved link identification by analyzing JS and CSS for urls.





var App = new Object();
App.host = "localhost:4005"

App.disable = function(me,id)
{
  var ele = me.parent().parent().parent().find(id)
  ele.attr('disabled', true); 
  ele.fadeTo( "fast", 0.33 ); 
}

App.enable = function(me,id)
{
  var ele = me.parent().parent().parent().find(id)
  ele.attr('disabled', false); 
  ele.fadeTo( "fast", 1 ); 
}

App.waiting = function(me)
{
  me.parent().parent().find(".loader").show();
}

App.done_waiting = function(me)
{
  me.parent().parent().find(".loader").hide();
}

App.crawl = function(me,url)
{
  $(".output").html("");
  if (window["WebSocket"]) {
    var socket = new WebSocket("ws://"+ App.host +"/ws/crawl");

    socket.onopen = function(event) {
      socket.send(url);
    };

    socket.onmessage = function(event) {
      var data = event.data;
      if (data == "exit") {
        App.enable(me,".url");
        App.enable(me,".crawl");
        App.done_waiting(me);
        var sitemap_url = "http://" + App.host + "/details/" + url;
        $(".output").append("View sitemap at <a href=\""+ sitemap_url +"\">" + sitemap_url  + "</a><br />");
        socket.close();
      } else {
        $(".output").append(data + "<br />");
      }
      $(".terminal").scrollTop($(".terminal").prop("scrollHeight") + 10);
    };
  } else {
    var params = { "url": url, };
    $.ajax({ type: "POST", url: "/crawl", data: params }).done(function( answer ) {
      App.enable(".url");
      App.enable(".crawl");
      App.done_waiting(me);
    });
  }
}

App.go = function(me,url) 
{
  $(".output").html("");

  App.waiting(me);
  App.disable(me,".url");
  App.disable(me,".crawl");

  $.ajax({ type: "GET", url: "/exists/" + url, data: {} }).done(function( answer ) {
    if (answer == "true") {
      App.done_waiting(me);
      window.location.assign("/details/" + url)
    } else {
      App.crawl(me,url);
    }
  });  
}
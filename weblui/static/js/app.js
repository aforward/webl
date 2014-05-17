
var App = new Object();

App.disable = function(id)
{
  var ele = $(id);
  ele.attr('disabled', true); 
  ele.fadeTo( "fast", 0.33 ); 
}

App.enable = function(id)
{
  var ele = $(id);
  ele.attr('disabled', false); 
  ele.fadeTo( "fast", 1 ); 
}


App.waiting = function()
{
  $(".loader").show(); 
}

App.done_waiting = function(did_work)
{
  $(".loader").hide();
}

App.create_crawl_socket = function() {
  var ws = new WebSocket("ws://{{$}}/crawl");
  return ws;
}

App.ws = App.create_crawl_socket();
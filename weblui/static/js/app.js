
var App = new Object();

App.disable = function(id)
{
  var ele = $(id);
  ele.prop('disabled', true); 
  ele.fadeTo( "fast", 0.33 ); 
}

App.waiting = function()
{
  $(".loader").show(); 
}

App.done_waiting = function(did_work)
{
  $(".loader").hide();
}

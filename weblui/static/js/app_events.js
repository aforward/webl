
//*****************
// CRAWL
//*****************

$(".crawl").on("click", function() {
  var params = {
    "url": $(".url").val(),
  };
  App.waiting();
  App.disable(".url");
  App.disable(".crawl");

  $.ajax({ type: "POST", url: "/crawl", data: params }).done(function( answer ) {
    App.done_waiting();
    location = "/aha";
  });
});
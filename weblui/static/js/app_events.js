
//*****************
// CRAWL
//*****************

$(".go").click(function() {
  var me = $(this);
  var url = $(".url").val();
  App.go(me,url);
});

$(".crawl").click(function() {
  var me = $(this);
  App.crawl(me,me.data("url"));
});

//*****************
// DELETE
//*****************

$(".delete").click(function() {
  var me = $(this)
  var url = me.data("url");

  App.waiting(me);
  App.disable(me,".crawl");
  App.disable(me,".delete");

  $.ajax({ type: "POST", url: "/delete/" + url, data: {} }).done(function( answer ) {
    App.done_waiting(me);
    var tr = me.parents(".domain-row");
    tr.fadeOut(400, function(){
      tr.remove();
    });
  });

});

//*****************
// ON LOAD
//*****************


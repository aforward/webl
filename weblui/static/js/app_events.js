
//*****************
// CRAWL
//*****************

$(".crawl").click(function() {
  var url = $(".url").val() 

  App.waiting();
  App.disable(".url");
  App.disable(".crawl");

  if (window["WebSocket"]) {
    console.log("websocket");
    var socket = new WebSocket("ws://localhost:4005/crawl");

    socket.onopen = function(event) {
      socket.send(url);
    };

    socket.onmessage = function(event) {
      var data = event.data;
      if (data == "exit") {
        App.enable(".url");
        App.enable(".crawl");
        App.done_waiting(true);
        socket.close();
      } else {
        $(".output").append(data);
      }
    };

  } else {
    var params = { "url": url, };
    $.ajax({ type: "POST", url: "/crawl", data: params }).done(function( answer ) {
      App.done_waiting();
      alert(answer)
    });
  }

});


      $(function() {
        var ws = new WebSocket("ws://localhost:8080/echo");
        ws.onmessage = function(e) {
          console.log("受信メッセージ:" + event.data);
        };

        var $ul = $('#msg-list');
        $('#sendBtn').click(function(){
          var data = $('#name').val();
          ws.send(data);
          console.log("送信メッセージ:" + data);
          $('<li>').text(data).appendTo($ul);
        });
      });
if (location.pathname === "/new") {
  var host = location.href.replace(/^http/, 'ws')
  var ws = new WebSocket(host);

  $(document).on('keyup', 'textarea', function() {
    var body = $('textarea').val();
    ws.send(body);
  });

  ws.onmessage = function(e) {
    var container = document.querySelector('.post');
    container.innerHTML = JSON.parse(e.data);
    $('pre code').each(function(i, e) {hljs.highlightBlock(e)});
 };
}

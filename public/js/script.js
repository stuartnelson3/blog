if (location.pathname === "/new_post") {
  var host = location.href.replace(/^http/, 'ws').replace(/new_post$/, 'markdown_preview')
  var ws = new WebSocket(host);

  $(document).on('keyup', 'textarea, .input', function() {
    var title = "# " + $('.input').val();
    var body = $('textarea').val();
    ws.send(title + "\n" + body);
  });

  ws.onmessage = function(e) {
    var container = document.querySelector('.post');
    container.innerHTML = e.data;
    $('pre code').each(function(i, e) {hljs.highlightBlock(e)});
 };
}

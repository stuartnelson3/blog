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

function addImageMarkdown(e, altText, imagePath) {
  var startPos = e.selectionStart;
  var endPos = e.selectionEnd;
  e.value = e.value.substring(0, startPos)
  + "![" + altText + "](" + imagePath + ")"
  + "\r\n"
  + e.value.substring(endPos, e.value.length);
}

$(document).on('drop', 'textarea', function(e) {
  e.preventDefault();
  var dt = e.originalEvent.dataTransfer;
  for (var i = 0; i < dt.files.length; i++) {
    var f = dt.files[i];

    t = e.currentTarget;
    addImageMarkdown(t, f.name, " ... ");

    var fd = new FormData();
    fd.append("file", f);

    var request = new XMLHttpRequest();
    request.open("POST", "/upload");

    request.onloadend = function(filename) {
      // Form a closure around the filename.
      return function(e) {
        var resp = e.currentTarget.response;
        // Update the text with the correct path for the image.
        var re = new RegExp("\\[" + filename + "\\]\\((\\s\\.\\.\\.\\s)\\)");
        t.value = t.value.replace(re, function($1, $2) {
          return $1.replace($2, resp.slice(1));
        });
        $(t).keyup();
      };
    }(f.name);
    request.send(fd);
  }

});

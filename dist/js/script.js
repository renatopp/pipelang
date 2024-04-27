$(function () {

  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("/js/pipe.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
  });

  $('#terminal').terminal(function (command, term) {

    if (command === '') {
      term.echo('')
      return
    }

    const res = pipe_eval(command)
    var div = ''
    if (res.startsWith('Error')) {
      div = $('<span class="response response-error">' + res + '</span>')
    } else if (res.startsWith('<')) {
      div = $('')
    } else {
      div = $('<span class="response response-ok">' + res + '</span>')
    }
    term.echo(div);
  }, {
    greetings: false,
    prompt: 'pipe> ',
    height: 400,
    keydown: function (e) {
      if (e.which == 82 && e.ctrlKey) {
        return true;
      }
    }
  }).focus();

  $.terminal.syntax('haskell')
});
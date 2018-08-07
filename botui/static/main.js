$(function() {
  // var DB_PATH = '/Users/cesar/coding/aerolith-infra/lexica/db/';
  var DB_PATH = '/home/ubuntu/word_db/';

  $.jsonRPC.setup({
    endPoint: '/rpc',
    namespace: ''
  });

  $('#connect-btn').click(connectHandler);
  $('#lexicon-change').click(lexiconChangeHandler);

  function connectHandler() {
    $.jsonRPC.request('AerobotService.Start', {
      params: {
        username: $('#username').val(),
        password: $('#password').val(),
        channel: $('#channel').val(),
        lexiconDb: $('#lexicon-db').val(),
      },
      success: function(result) {
        console.log(result);
      },
      error: function(result) {
        // Result is an RPC 2.0 compatible response object
      }
    });
  }

  function lexiconChangeHandler() {
    var lexName = $('#lexicon-select').val();
    if (lexName === 'TWL3.1') {
      lexName = 'America';
    }
    $.jsonRPC.request('AerobotService.LoadLexicon', {
      params: {
        newLexiconDb: DB_PATH + lexName + '.db'
      },
      success: function(result) {
        console.log(result);
      }
    });
  }

});

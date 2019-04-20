'use strict';

function updateTimeEls() {
  setInterval(() => {
    $('[data-ts]')
        .each((i, el) => {
          var $el = $(el);
          $el.text(formatTime($el.attr('data-ts')));
        });
  }, 1000);
}

function formatTime(ts) {
  var delta = new Date().getTime() - ts;

  if (delta < 1000) {
    return 'now'
  }

  if (delta < 1000 * 60) {
    return Math.round(delta / 1000) + 's ago';
  }

  if (delta < 1000 * 60 * 60) {
    var s = delta % (1000 * 60);
    var m = delta - s;
    return Math.round(m / (1000 * 60)) + 'm ' + Math.round(s / 1000) + 's ago';
  }

  var d = new Date();
  d.setTime(ts);

  if (delta < 1000 * 60 * 60 * 24) {
    return d.toLocaleTimeString();
  }

  return d.toLocaleDateString();
}
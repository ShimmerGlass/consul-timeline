'use strict';

Date.prototype.toDatetimeLocal = function toDatetimeLocal() {
  var date = this, ten = function(i) { return (i < 10 ? '0' : '') + i; },
      YYYY = date.getFullYear(), MM = ten(date.getMonth() + 1),
      DD = ten(date.getDate()), HH = ten(date.getHours()),
      II = ten(date.getMinutes()), SS = ten(date.getSeconds());
  return YYYY + '-' + MM + '-' + DD + 'T' + HH + ':' + II + ':' + SS;
};

Date.prototype.fromDatetimeLocal = (function(BST) {
  // BST should not be present as UTC time
  return new Date(BST).toISOString().slice(0, 16) === BST ?
      // if it is, it needs to be removed
      function() {
        return new Date(this.getTime() + (this.getTimezoneOffset() * 60000))
            .toISOString();
      } :
      // otherwise can just be equivalent of toISOString
      Date.prototype.toISOString;
}('2006-06-06T06:06'));

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
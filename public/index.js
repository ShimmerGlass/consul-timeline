'use strict';

function getEvents(start, service, limit, cb) {
  $.getJSON(
    "/events?limit=" + limit + "&start=" + Math.round(start / 1000) +
    "&service=" + (service || ''),
    cb)
}

Date.prototype.toDatetimeLocal =
  function toDatetimeLocal() {
    var
      date = this,
      ten = function (i) {
        return (i < 10 ? '0' : '') + i;
      },
      YYYY = date.getFullYear(),
      MM = ten(date.getMonth() + 1),
      DD = ten(date.getDate()),
      HH = ten(date.getHours()),
      II = ten(date.getMinutes()),
      SS = ten(date.getSeconds())
      ;
    return YYYY + '-' + MM + '-' + DD + 'T' +
      HH + ':' + II + ':' + SS;
  };

Date.prototype.fromDatetimeLocal = (function (BST) {
  // BST should not be present as UTC time
  return new Date(BST).toISOString().slice(0, 16) === BST ?
    // if it is, it needs to be removed
    function () {
      return new Date(
        this.getTime() +
        (this.getTimezoneOffset() * 60000)
      ).toISOString();
    } :
    // otherwise can just be equivalent of toISOString
    Date.prototype.toISOString;
}('2006-06-06T06:06'));

class Vue {
  constructor(container) {
    this.container = container;
    this.renderer = new Renderer(container, (dir) => { this.fetch(dir); });

    this.requestLimit = 100;

    this.filter = this.filterFromUrl();
    $('#status .service').val(this.filter.service || '');
    if (this.filter.start) {
      var d = new Date();
      d.setTime(this.filter.start);
      $('#time-selector-btn').text(d.toLocaleString());
      $('#time-selector .custom-in').val(d.toDatetimeLocal());
    } else {
      $('#time-selector .custom-in').val(new Date().toDatetimeLocal());
    }

    this.reset();

    var that = this;

    $('#status .service')
      .change(function () {
        var n = $(this).val();
        if (n == that.filter.service) {
          return;
        }

        that.filter.service = n;
        that.updateUrlFromFilter(that.filter);
        that.reset();
      });

    $("#back-to-top").click(() => { this.renderer.reset(); });

    $('#time-selector-btn').click(() => { $('#time-selector').toggle(); });
    $('#time-selector [data-val]')
      .click(function () {
        $('#time-selector-btn').text($(this).text());
        $('#time-selector [data-val]').removeClass('active');
        $(this).addClass('active');

        var offset = parseInt($(this).attr('data-val'));
        if (!offset) {
          that.filter.start = null;
        } else {
          that.filter.start =
            new Date().getTime() + parseInt($(this).attr('data-val'));
        }
        that.updateUrlFromFilter(that.filter);
        that.reset();
        $('#time-selector').hide();
      });

    $('#time-selector .custom-btn').click(() => {
      var time = $('#time-selector .custom-in').val();
      var d = new Date(time);
      $('#time-selector-btn').text(d.toLocaleString());

      that.filter.start = d.getTime();
      that.updateUrlFromFilter(that.filter);
      that.reset();
      $('#time-selector').hide();
    })
  }

  reset() {
    this.stopListenNew();
    this.renderer.empty();
    this.fetching = false;
    this.noMoreData = {};

    if (this.filter.start) {
      this.startTime = this.filter.start;
    } else {
      this.startTime = new Date().getTime();
      this.listenNew();
    }
    this.endTime = this.startTime;

    this.fetch(1);
  }

  fetch(dir) {
    if (this.fetching) {
      return;
    }

    if (this.noMoreData[dir]) {
      return;
    }

    var start;
    if (dir == 1) {
      start = this.endTime;
    } else {
      start = this.startTime;
    }

    this.fetching = true;
    getEvents(start, this.filter.service, this.requestLimit, data => {
      this.fetching = false;
      this.renderer[dir == 1 && 'appendLogs' || 'prependLogs'](data);
      if (data.length < this.requestLimit) {
        this.noMoreData[dir] = true;
        this.renderer.appendMessage('no_more_data', 'End of logs');
      }
      if (data.length && dir == 1) {
        this.endTime = new Date(data[data.length - 1].time).getTime();
      }
      if (data.length && dir == -1) {
        this.startTime = new Date(data[0].time).getTime();
      }
    });
  }

  listenNew() {
    var that = this;
    this.ws = new WebSocket(
      "ws://" + window.location.host + "/ws?" + this.filterToQs(this.filter));
    this.ws.onmessage = function (evt) {
      var e = JSON.parse(evt.data);
      that.startTime = e.time;
      that.renderer.prependLogs([e]);
    };

    this.ws.onopen = function () {
      that.wsShouldReconnect = true;
      $('#status .ws').html('<i class="fas fa-check passing"></i> Live');
    };

    this.ws.onclose = function () {
      if (!that.wsShouldReconnect) {
        return;
      }
      $('#status .ws')
        .html('<i class="fas fa-redo-alt fa-spin critical"></i> Disonnected');
      that.renderer.prependMessage(
        'ws_disconnect',
        'Live refresh got disconnected, events might be missing');
      that.wsReconnectTimeout =
        setTimeout(function () { that.listenNew(); }, 2000);
    }
  }

  stopListenNew() {
    this.wsShouldReconnect = false;
    clearTimeout(this.wsReconnectTimeout);
    $('#status .ws').html('<i class="fas fa-times warning"></i> Paused');
    if (!this.ws) {
      return;
    }

    this.ws.close();
  }

  filterFromUrl() {
    var filter = {};

    var searchParams = new URLSearchParams(window.location.search);
    filter.service = searchParams.get('service');
    filter.start = searchParams.get('start');

    return filter;
  }

  updateUrlFromFilter(filter) {
    var newRelativePathQuery =
      window.location.pathname + '?' + this.filterToQs(filter);
    history.pushState(null, '', newRelativePathQuery);
  }

  filterToQs(filter) {
    var searchParams = new URLSearchParams();
    if (filter.service) {
      searchParams.set("service", filter.service);
    } else {
      searchParams.delete("service");
    }
    if (filter.start) {
      searchParams.set("start", filter.start);
    } else {
      searchParams.delete("start");
    }
    return searchParams.toString();
  }
}

$(function () {
  new Vue($('#container'));
  new Select($('#status .service'), (cb) => { $.getJSON("/services", cb); });
});

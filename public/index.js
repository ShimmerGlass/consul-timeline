'use strict';

$(document).ajaxStart(() => { $("#loader").show() });
$(document).ajaxStop(() => { $("#loader").hide() });


function getEvents(start, filter, limit, cb) {
  $.getJSON(
    "/events?limit=" + limit + "&start=" + Math.round(start / 1000) +
    "&filter=" + (filter || ''),
    cb)
}

class Vue {
  constructor(container) {
    this.container = container;
    this.renderer = new Renderer(container, (dir) => { this.fetch(dir); });

    this.requestLimit = 100;

    this.filter = this.filterFromUrl();
    $('#status .filter').val(this.filter.filter || '');
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

    $('#status .filter')
      .change(function () {
        var n = $(this).val();
        if (n == that.filter.filter) {
          return;
        }

        that.filter.filter = n;
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

    $('#time-selector .custom-btn')
      .click(() => {
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
    getEvents(start, this.filter.filter, this.requestLimit, data => {
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
    var scheme = window.location.protocol == 'https:' ? 'wss' : 'ws';
    this.ws = new WebSocket(
      scheme + "://" + window.location.host + "/ws?" + this.filterToQs(this.filter));
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
    filter.filter = searchParams.get('filter');
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
    if (filter.filter) {
      searchParams.set("filter", filter.filter);
    } else {
      searchParams.delete("filter");
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
  new Select($('#status .filter'), (cb) => { $.getJSON("/filter-entries", cb); });
});

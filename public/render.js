'use strict';

var serviceStatusClasses = {
  0: "unknown",
  1: "missing",
  2: "critical",
  3: "warning",
  4: "passing",
  5: "critical",
};

class Renderer {
  constructor(container, requestLogs) {
    this.container = container;
    this.requestLogs = requestLogs;
    this.logs = [];
    this.maxLogs = 2000;  // max number of logs to hold in memory
    this.logsByKey = {};  // used for rows merging
    this.startIndex = 0;
    this.endIndex = 0;
    this.rowHeight = 40;
    this.topOddness = false;
    this.bottomOddness = true;

    this.visibleRows = Math.ceil(container.innerHeight() / this.rowHeight);
    this.rowOverflow = this.visibleRows * 3;

    this.container.scroll(() => { this.update(); });

    updateTimeEls();
  }

  update() {
    window.cancelAnimationFrame(this.annimationFrameId);
    this.annimationFrameId = window.requestAnimationFrame(() => {
      var follow = this.container.scrollTop() == 0;

      // if there are too many rows to draw (for example because we were paused
      // in background), just redraw the screen
      if (follow && this.startIndex > 100) {
        this.empty();
        return;
      }

      this.fillLeading();
      this.removeLeading();
      this.fillTrailing();
      this.removeTrailing();

      if (follow) {
        this.container.scrollTop(0);
        $("#back-to-top").hide();
      } else {
        $("#back-to-top .count").text(this.visibleStartIndex());
        $("#back-to-top").show();
      }
    });
  }

  appendLogs(logs) {
    if (!logs || !logs.length) {
      return;
    }
    this.logs = this.logs.concat(logs);
    this.update();
  }

  prependLogs(logs) {
    if (!logs.length) {
      return;
    }
    this.logs = logs.concat(this.logs);
    this.startIndex += logs.length;
    this.endIndex += logs.length;
    this.update();
  }

  prependMessage(kind, message) {
    if (this.logs.length && this.logs[0].kind == kind) {
      return;
    }

    this.prependLogs([{
      type: "message",
      kind: kind,
      message: message,
    }]);
  }

  appendMessage(kind, message) {
    if (this.logs.length && this.logs[this.logs.length - 1].kind == kind) {
      return;
    }

    this.appendLogs([{
      type: "message",
      kind: kind,
      message: message,
    }]);
  }

  reset() {
    this.topOddness = false;
    this.bottomOddness = true;
    this.startIndex = 0;
    this.endIndex = 0;
    this.container.find('.row').remove();
    this.update();
  }

  empty() {
    this.topOddness = false;
    this.bottomOddness = true;
    this.startIndex = 0;
    this.endIndex = 0;
    this.logs = [];
    this.logsByKey = {};
    this.container.find('.row').remove();
  }

  rowCount() { return this.logs.length; }

  removeLeading() {
    var scrollTop = this.container.scrollTop();
    var nToRemove = scrollTop / this.rowHeight - this.rowOverflow;
    if (nToRemove <= this.rowOverflow / 2) {
      return;
    }
    var rowEls = this.container.find('.row');
    for (var i = 0; i < nToRemove; i++) {
      rowEls[i].remove();
      this.startIndex++;
    }
    this.container.scrollTop(scrollTop - this.rowHeight * nToRemove);
  }

  fillLeading() {
    var scrollTop = this.container.scrollTop();
    var leading = scrollTop / this.rowHeight;
    var nToAdd = this.rowOverflow - leading;
    if (nToAdd <= this.rowOverflow / 2 &&
      this.startIndex > this.rowOverflow / 2) {
      return;
    }

    var added = 0;
    for (var i = 0; i < nToAdd; i++) {
      if (this.startIndex <= 0) {
        break;
      }

      this.startIndex--;
      this.drawRow(-1, this.logs[this.startIndex]);
      added++;
    }
    this.container.scrollTop(scrollTop + added * this.rowHeight);
  }

  removeTrailing() {
    var actual = (this.container.prop('scrollHeight') -
      this.container.scrollTop() - this.container.innerHeight()) /
      this.rowHeight;

    if (actual - this.rowOverflow <= this.rowOverflow / 2) {
      return;
    }

    var rowEls = this.container.find('.row');
    for (var i = 0; i < actual - this.rowOverflow; i++) {
      rowEls[rowEls.length - 1 - i].remove();
      this.endIndex--;
    }
  }

  fillTrailing() {
    var visibleStartIndex = this.visibleStartIndex();
    var wanted = visibleStartIndex + this.visibleRows + this.rowOverflow;

    if (wanted - this.endIndex <= this.rowOverflow / 2) {
      return;
    }

    for (var i = this.endIndex; i < wanted; i++) {
      if (i >= this.logs.length) {
        this.requestLogs(1);
        break;
      }

      this.drawRow(1, this.logs[i]);
      this.endIndex++;
    }
  }

  drawRow(direction, d) {
    if (d.type == 'message') {
      row = $('<div class="row message">' + d.message + '</div>');
      this.container[direction == 1 && 'append' || 'prepend'](row);
      return;
    }

    var key = d.time + '_' + d.node_name + '_' + d.service_id;
    var subKey = key + '_' + d.check_name;
    this.logsByKey[subKey] = d;

    var rows = this.container.find('.row[data-key="' + key + '"]').toArray();
    var row =
      $('<div class="row" data-key="' + key + '" data-subkey="' + subKey +
        '"></div>');

    if (rows.length) {
      $(rows[0]).hasClass('odd') && row.addClass('odd');
      $(rows[0]).hasClass('even') && row.addClass('even');
      if (direction == 1) {
        row.insertAfter(rows[rows.length - 1]);
        rows.push(row);
      } else {
        row.insertBefore(rows[0]);
        rows.unshift(row);
      }
    } else {
      if (direction == 1) {
        this.bottomOddness = !this.bottomOddness;
        row.addClass(this.bottomOddness && 'odd' || 'even');
        this.container.append(row);
      } else {
        this.topOddness = !this.topOddness;
        row.addClass(this.topOddness && 'odd' || 'even');
        this.container.prepend(row);
      }
      rows.push(row);
    }

    for (var i = 0; i < rows.length; i++) {
      var row = $(rows[i]);
      var d = this.logsByKey[row.attr('data-subkey')];

      row.empty();

      ((d, row) => { row.click(() => { this.drawDetails(d); }); })(d, row);

      if (i == 0) {
        var time = new Date(d.time).getTime();
        var html = '';

        html += '<div class="time" data-ts="' + time + '" title="' + d.time +
          '">' + formatTime(time) + '</div>';

        html += '<div class="node" title="' + d.node_ip + '">';
        if (d.old_node_status && d.new_node_status) {
          html += this.getStatusesMarkup(d.old_node_status, d.new_node_status);
        }
        html += '&nbsp;&nbsp;' + d.node_name + '</div>';

        if (d.service_id) {
          html += '<div class="service" title="' + d.service_id + '">';
          html += this.getStatusesMarkup(
            d.old_service_status, d.new_service_status);
          html += '&nbsp;&nbsp;' + d.service_name + '&nbsp;&nbsp;';
          html += '(' + d.old_instance_count +
            '&nbsp;<i class="fas fa-arrow-right missing" style="font-size: 0.8em"></i>&nbsp;' +
            d.new_instance_count + ')';
          html += '</div>';
        } else {
          html += '<div class="service"></div>';
        }

        row.append($(html));
      } else {
        var html = '';
        html += '<div class="time"></div>';
        html += '<div class="node"></div>';
        html += '<div class="service"></div>';
        row.append($(html));
      }

      if (d.check_name) {
        var html = '';
        html += '<div class="check">';
        html += this.getStatusesMarkup(d.old_check_status, d.new_check_status);
        html += '&nbsp;&nbsp;' + d.check_name;
        html += '</div>';

        html += '<div class="check-output">';
        html += $('<div />').text(d.check_output || '').html();
        html += '</div>';
        row.append($(html));
      }
    }
  }

  getStatusMarkup(statusCode) {
    var icon;
    var title;
    switch (statusCode) {
      case 1:
        icon = 'fa-question';
        title = 'Missing';
        break;
      case 2:
        icon = 'fa-times';
        title = 'Critical';
        break;
      case 3:
        icon = 'fa-exclamation-triangle';
        title = 'Warning';
        break;
      case 4:
        icon = 'fa-check';
        title = 'Passing';
        break;
      case 5:
        icon = 'fa-wrench';
        title = 'Maintenance';
        break;
    }

    return '<i class="fas ' + icon + ' ' + serviceStatusClasses[statusCode] +
      '" style="font-size: 1.2em" title="' + title + '"></i>';
  }

  getStatusesMarkup(oldStatus, newStatus) {
    var html = "";
    html += this.getStatusMarkup(oldStatus) + "";
    html +=
      '&nbsp;<i class="fas fa-arrow-right missing" style="font-size: 0.8em"></i>&nbsp;';
    html += this.getStatusMarkup(newStatus);
    return html;
  }

  visibleStartIndex() {
    return Math.floor(
      this.startIndex + (this.container.scrollTop() / this.rowHeight));
  }

  drawDetails(d) {
    var res = '<ul>';
    for (var i in d) {
      if (!d[i]) {
        continue;
      }
      var k = i.replace(/_/g, ' ');
      res += '<li><span class="key">' + k + ':</span> ';
      if (i.endsWith('_status')) {
        res += this.getStatusMarkup(d[i])
      } else if (i == 'check_output') {
        res += '<pre class="value">' + $('<div />').text(d[i] || '').html() +
          '</pre>';
      } else {
        res += d[i]
      }
      res += '</li>';
    }
    res += '</ul>';

    $('#details-pane .contents').html(res);
    $('#details-pane .close').one('click', () => { this.closeDetails(); });

    this.openDetails();
  }

  openDetails() {
    this.container.addClass('small');
    $('#details-pane').show();
  }

  closeDetails() {
    this.container.removeClass('small');
    $('#details-pane').hide();
  }
}
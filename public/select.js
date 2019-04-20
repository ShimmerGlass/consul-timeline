'use strict';

class Select {
  constructor(el, dataFn) {
    this.el = el;
    this.dataFn = dataFn;
    this.selectedIdx = 0;
    this.lastVal = '';
    this.entries = [];
    this.data = [];

    this.refresh();
    this.dropdown = $('<ul class="select-dd"></ul>')
                        .hide()
                        .appendTo($(window.document.body));

    el.focus(() => {
      var val = this.el.val();
      this.lastVal = val;
      this.entries = this.getEntries(val);
      this.open();
    });
    el.blur(() => { setTimeout(() => { this.dropdown.hide(); }, 100); });
    el.keyup((evt) => {
      this.dropdown.show();

      if (evt.code == "ArrowDown" && this.selectedIdx < this.entries.length) {
        this.selectedIdx++;
      } else if (evt.code == "ArrowUp" && this.selectedIdx > 0) {
        this.selectedIdx--;
      } else if (
          evt.code == "Enter" && this.selectedIdx < this.entries.length) {
        this.select(this.entries[this.selectedIdx]);
      }
      var val = this.el.val();
      if (val != this.lastVal) {
        this.selectedIdx = 0;
        this.lastVal = val;
      }
      this.entries = this.getEntries(val);
      this.update();
    });
  }

  refresh() {
    this.dataFn(d => { this.data = d; });
  }

  open() {
    this.update();
    this.dropdown.show();
  }

  update() {
    var that = this;
    this.dropdown.empty();
    for (var i in this.entries) {
      ((i, e) => {
        var el = $('<li>' + e + '</li>');
        if (i == this.selectedIdx) {
          el.addClass('selected');
        }
        el.click(() => {that.select(e)})

            this.dropdown.append(el);
      })(i, this.entries[i]);
    }

    var elLoc = this.el.offset();
    this.dropdown.css({
      width: this.el.outerWidth(),
      top: elLoc.top + this.el.outerHeight(),
      left: elLoc.left,
    });
  }

  select(val) {
    this.el.val(val);
    this.el.change();
    this.dropdown.hide();
  }

  getEntries(val) {
    var fuse = new Fuse(this.data, {});
    var idx = fuse.search(val);

    var res = [];
    for (var i in idx) {
      if (i >= 20) {
        break;
      }
      res.push(this.data[idx[i]]);
    }
    return res;
  }
}
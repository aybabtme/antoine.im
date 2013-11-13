// Convenience hacks
var puts = function (data) {
  "use strict";
  document.write(data);
};

var math = function (data) {
  "use strict";
  return "$" + data + "$";
};

var randomInt = function (from, to) {
  "use strict";
  var r = Math.random();
  return Math.round((to - from) * r) + from;
};

// To make HTML tables from JS
var VerticalTable = function (header) {
  "use strict";
  var rows = [];

  var rowAdder = function (row) {
    rows.push(row);
    return this;
  };

  var htmler = function () {
    var str = "<table>" +
      "\n\t<tr>";
    header.forEach(function (el) {
      str += "<th>" + el;
    });
    rows.forEach(function (row) {
      str += "\n\t<tr>";
      row.forEach(function (col) {
        str += "\t\t<td>" + col;
      });
    });
    return str += "\n</table>";
  };

  return {
    addEntry: rowAdder,
    toHTML: htmler
  };
};

var HorizontalTable = function (header) {
  "use strict";
  var cols = [];
  for (var i = 0; i < header.length; i++) {
    cols[i] = [header[i]];
  }

  var colAdder = function (col) {
    for (var i = 0; i < cols.length; i++) {
      cols[i].push(col[i]);
    }
    return this;
  };

  var htmler = function () {
    var str = "<table>";

    cols.forEach(function (row) {
      str += "\n\t<tr>";
      row.forEach(function (col) {
        str += "\t\t<td>" + col;
      });
    });

    return str += "\n</table>";
  };

  return {
    addEntry: colAdder,
    toHTML: htmler
  };
};

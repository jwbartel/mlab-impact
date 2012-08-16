// Copyright 2011 Google Inc. All Rights Reserved.

/**
 * @fileoverview This is the root module for the net-score project.
 *
 * @author gavaletz@google.com (Eric Gavaletz)
 * Additions made by bartel@cs.unc.edu (Jacob Bartel)
 */


goog.provide('results');

goog.require('goog.dom');

results.addListener = function(chart, eventName, callback) {
  google.visualization.events.addListener(chart, eventName, callback);
};


/**
 * Testing out the charting tools API.
 * http://code.google.com/apis/chart/interactive/docs/index.html
 *
 * http://code.google.com/apis/chart/interactive/docs/gallery/scatterchart.html
 */
results.scatterChart = function(dataObj, title, parentDiv) {
  var data = new google.visualization.DataTable(dataObj);
  var chart = new google.visualization.ScatterChart(parentDiv);
  var options = {'curveType': 'none', 'lineWidth': 1, 'pointSize': 0,
    'width': page.displayWidth, 'height': page.displayHeight,
    'title': title, 'vAxis': {'title': dataObj['cols'][1]['label']},
    'hAxis': {'title': dataObj['cols'][0]['label']},
    'legend': {'position': 'none'},
    'titlePosition': 'in'};
  chart.draw(data, options);
  return chart;
};


/**
 * http://code.google.com/apis/chart/interactive/docs/gallery/table.html
 */
results.table = function(dataObj, parentDiv) {
  var data = new google.visualization.DataTable(dataObj);
  var chart = new google.visualization.Table(parentDiv);
  var options = {'width': page.displayWidth, 'allowHtml': true,
    'page': 'enable', 'pageSize': 30};
  chart.draw(data, options);
  return chart;
};


/**
 * http://code.google.com/apis/chart/interactive/docs/gallery/piechart.html
 */
results.pieChart = function(dataObj, title, parentDiv) {
  var chart = new google.visualization.PieChart(parentDiv);
  var options = {'width': page.displayWidth, 'height': page.displayHeight,
    'title': title};
  chart.draw(dataObj, options);
  return chart;
};


/**
 * https://developers.google.com/chart/interactive/docs/gallery/columnchart
 */
results.columnChart = function(dataObj, title, showLegend, width,
    parentDiv, parentPos) {

  var chart = new google.visualization.ColumnChart(parentDiv);
  var legendPosition = (showLegend) ? 'bottom' : 'none';

  var options = {
                  'title': title,
                  'legend': {'position': legendPosition},
                  'colors': page.columnColors,
                  'chartArea': {'width': '80%'},
                  'width': width
                };
  chart.draw(dataObj, options);
  return chart;
};

// http://code.google.com/apis/chart/interactive/docs/reference.html#dataparam
results.dataToObj = function(x, y, xLabel, yLabel, xType, yType) {
  var dataObj = {'cols': [], 'rows': []};
  dataObj['cols'][0] = {'id': 'x', 'label': xLabel, 'type': xType};
  dataObj['cols'][1] = {'id': 'y', 'label': yLabel, 'type': yType};

  // if there is not an array for x then us a count.
  // otherwise use the smaller length to avoid index errors
  if (x == null) {
    for (i = 0; i < y.length; i++) {
      dataObj['rows'][i] = {'c': [{'v': i}, {'v': y[i]}]};
    }
  }
  else {
    for (i = 0; i < Math.min(x.length, y.length); i++) {
      dataObj['rows'][i] = {'c': [{'v': x[i]}, {'v': y[i]}]};
    }
  }
  return dataObj;
};


//Similar to the above method, but for data sorted into categories
results.categorizedDataToObject = function(xLabel, xGroups, yLabel,
    data, categories, categoryTypes, GroupNames) {
  var dataObj = new google.visualization.DataTable();
  dataObj.addColumn('string', xLabel);

  // Limit the categories by types because each category must have a type
  // but a label may be inferred
  for (var i = 0; i < categoryTypes.length; i++) {
    var label = (categories) ? categories[i] : i;
    dataObj.addColumn(categoryTypes[i], label);
  }

  //Groups are rows of categorized data
  numGroups = (xGroups) ? Math.min(xGroups.length, data.length) : data.length;
  for (var i = 0; i < numGroups; i++) {
    row = [];

    var groupName = (GroupNames) ? GroupNames[i] : (xGroups) ?
        xGroups[i] : '' + i;

    row[0] = groupName;

    for (var j = 0; j < data[i].length; j++) {
      row[j + 1] = data[i][j];
    }
    dataObj.addRow(row);
  }

  return dataObj;
};

results.recordsToObj = function(records) {
  var dataObj = {'cols': [], 'rows': []};
  var keys = Object.keys(records[0]);
  var max_i = keys.length;
  var max_j = records.length;
  for (i = 0; i < max_i; i++) {
    dataObj['cols'][i] = {
      'id': 'c' + i,
      'label': keys[i],
      'type': typeof records[0][keys[i]]
    };
  }
  for (j = 0; j < max_j; j++) {
    dataObj['rows'][j] = {'c': []};
    for (i = 0; i < max_i; i++) {
      dataObj['rows'][j]['c'].push({'v': records[j][keys[i]]});
    }
  }
  return dataObj;
};

// Copyright 2011 Google Inc. All Rights Reserved.

/**
 * @fileoverview This module provides a way to display our site.
 *
 * @author gavaletz@google.com (Eric Gavaletz)
 */

goog.provide('page');


goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.events.EventType');
goog.require('goog.object');
goog.require('goog.ui.Button');
goog.require('goog.ui.Component.EventType');
goog.require('goog.ui.FlatButtonRenderer');
goog.require('goog.ui.FlatMenuButtonRenderer');
goog.require('goog.ui.MenuRenderer');
goog.require('goog.ui.MenuSeparator');
goog.require('goog.ui.ProgressBar');
goog.require('goog.ui.RoundedTabRenderer');
goog.require('goog.ui.RoundedTabRenderer');
goog.require('goog.ui.Tab');
goog.require('goog.ui.TabBar');


goog.require('input');
goog.require('mapRenderer');
goog.require('results');


/**
 * Is the user logged in as an administrator?
 *
 * @type {boolean}
 */
page.loggedInAdmin = false;


/**
 * Sets the admin logged in flag.  Usually just used from the server.
 *
 * @param {boolean} b Is the user logged in as an administrator?
 */
page.setLoggedInAdmin = function(b) {
  page.loggedInAdmin = b;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetLoggedInAdmin', page.setLoggedInAdmin);


/**
 * Is the user logged in?
 *
 * @type {boolean}
 */
page.loggedIn = false;


/**
 * Sets the logged in flag.  Usually just used from the server.
 *
 * @param {boolean} b Is the user logged in?
 */
page.setLoggedIn = function(b) {
  page.loggedIn = b;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetLoggedIn', page.setLoggedIn);


/**
 * The URL used to log in a user.
 *
 * @type {string}
 */
page.loginUrl = '';


/**
 * Sets the URL used to log in a user.  Usually just used from the server.
 *
 * @param {string} s the URL to be used.
 */
page.setLoginUrl = function(s) {
  page.loginUrl = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetLoginUrl', page.setLoginUrl);


/**
 * The URL used to log out a user.
 *
 * @type {string}
 */
page.logoutUrl = '';


/**
 * Sets the URL used to log out a user.  Usually just used from the server.
 *
 * @param {string} s the URL to be used.
 */
page.setLogoutUrl = function(s) {
  page.logoutUrl = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetLogoutUrl', page.setLogoutUrl);


/**
 * The user's name to be displayed.
 *
 * @type {string}
 */
page.userName = '';


/**
 * Sets the user's name.  Usually just used from the server.
 *
 * @param {string} s the user's name to be used.
 */
page.setUserName = function(s) {
  page.userName = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetUserName', page.setUserName);


/**
 * The width of the page in pixles.
 *
 * @type {number}
 */
page.displayWidth = 640;


/**
 * Sets the width of the page in pixles.  Usually just used from the server.
 *
 * @param {number} n the width to be used.
 */
page.setDisplayWidth = function(n) {
  page.displayWidth = n;
  page.setAttributes();
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'pageSetDisplayWidth', page.setDisplayWidth);


/**
 * The height of the page in pixles.
 *
 * @type {number}
 */
page.displayHeight = 360;


/**
 * Sets the height of the page in pixles.  Usually just used from the server.
 *
 * @param {number} n the height to be used.
 */
page.setDisplayHeight = function(n) {
  page.displayHeight = n;
  page.setAttributes();
};
//used by the closure compiler to export this function call
goog.exportProperty(window, 'pageSetDisplayHeight', page.setDisplayHeight);


/**
 * The colors of the columns in a chart
 *
 * @type {array}
 */
page.columnColors = ['#6082C4', '#BFCDE7', '#90A8D6'];


/**
 * Sets the attributes of the page when the height or width change.
 */
page.setAttributes = function() {
  page.attributes = {};
  page.attributes.display = {'style': 'height: ' + page.displayHeight +
        'px; width: ' + page.displayWidth + 'px;',
    'valign': 'middle', 'align': 'center'};

  page.attributes.table = {'style': 'width: ' + page.displayWidth + 'px;',
    'align': 'center'};

  page.attributes.message = {'style':
        'padding-top: 10px; padding-bottom: 10px; width: ' + page.displayWidth +
        'px;', 'valign': 'top', 'align': 'left'};
};


/**
 * Name space for progress components.
 *
 * @type {Object.{*}}
 */
page.progress = {};


/******************************************************************************
 * START MAIN PAGE
 *****************************************************************************/


/**
 * Setup and display the page header.
 */
page.headerSetup = function() {

  //Not necessary because we do not use logins to oauth for bigquery

  /*if (typeof page.header == 'undefined') {
    //header
    if (page.loggedIn) {
      var home = goog.dom.createDom('a', {'href': '/'}, 'home');
      var logout = goog.dom.createDom('a', {'href': page.logoutUrl},
          'logout ' + page.userName);
      if (page.loggedInAdmin) {
        //TODO(user) links for adminy things.
        page.header = goog.dom.createDom('div', {'id': 'header',
          'class': 'header'}, logout,
        ' | ', home);
      }
      else {
        page.header = goog.dom.createDom('div', {'id': 'header',
          'class': 'header'}, logout,
        ' | ', home);
      }
    }
    else {
      var login = goog.dom.createDom('a', {'href': page.loginUrl},
          'Google login');
      var home = goog.dom.createDom('a', {'href': '/'}, 'home');
      page.header = goog.dom.createDom('div', {'id': 'header',
        'class': 'header'}, login, ' | ', home);
    }
    goog.dom.appendChild(document.body, page.header);
  }*/
};


/**
 * Setup and display the page content section.
 */
page.contentSetup = function() {
  if (typeof page.content == 'undefined') {
    //content
    //EVERYTHING DISPLAYED AFTER THE WELCOME SHOULD GO IN HERE!
    var pushFooter = goog.dom.createDom('div', {'class': 'pushFooter'});
    //var pushHeader = goog.dom.createDom('div', {'class': 'pushHeader'});
    page.content = goog.dom.createDom('div', {'id': 'content'});
    var wrapper = goog.dom.createDom('div', {'class': 'wrapper'}, //pushHeader,
        page.content, pushFooter);
    goog.dom.appendChild(document.body, wrapper);
  }
};


/**
 * Setup and display the page footer.
 */
page.footerSetup = function() {
  if (typeof page.footer == 'undefined') {
    //footer
    var source = goog.dom.createDom('a',
        {'href': 'http://code.google.com/p/mlab-impact/', 'target': '_blank'},
        'source code');
    var mlab = goog.dom.createDom('a', {'href': 'http://www.measurementlab.net/',
      'target': '_blank'}, 'measurement lab');
    var netscore = goog.dom.createDom('a',
        {'href': 'http://www.net-score.org', 'target': '_blank'},
        'net-score');
    var usCensusAPI = goog.dom.createDom('a',
        {'href': 'http://www.census.gov/developers', 'target': '_blank'},
        'united states census api');
    var br = goog.dom.createDom('br');
    var legalese = 'This product uses the Census Bureau Data API but is not endorsed or certified by the Census Bureau.';
    page.footer = goog.dom.createDom('div', {'id': 'footer', 'class': 'footer'},
        source, ' | ', mlab, ' | ', netscore, ' | ', usCensusAPI, br, legalese);
    goog.dom.appendChild(document.body, page.footer);
  }
};


/**
 * Setup  the logo.
 */
page.logoSetup = function() {
  if (typeof page.logo == 'undefined') {
    page.logo = goog.dom.createDom('img',
        {'src': '/images/impact-logo.png',
          'alt': 'M-Lab Impact', 'id': 'logo'});
  }
};


/**
 * Setup the spinner to show when results are being loaded
 */
page.spinnerSetup = function() {
  if (typeof page.spinner == 'undefined') {
    page.spinner = goog.dom.createDom('img',
        {'src': '/images/spinner.gif',
          'alt': 'Loading...', 'id': 'spinner',
          'style': 'display:none;'});
  }
};


/**
 * Resets the display to remove any old query results
 */
page.resetAnswerSpace = function() {
  if (page.answer != null) goog.dom.removeChildren(page.answer);
  //if (page.pureData != null) page.pureData.style.display = 'none';
  if (page.map != null) page.map.style.opacity = 0.0;
};


/**
 * Initialize all the elements for building query forms
 */
page.querySetup = function() {
  if (typeof page.query == 'undefined') {

    var question = goog.dom.createDom('div',
        {'id': 'question'});
    input.queryFormSetup(question);
    page.answer = goog.dom.createDom('div',
        {'id': 'answer'});

    page.logoSetup();
    page.spinnerSetup();
    page.leftColumn = goog.dom.createDom('td', {'id': 'leftColumn'},
        page.spinner, page.logo, question);
    page.rightColumn = goog.dom.createDom('td',
        {'id': 'rightColum'}, page.answer);

    contentRow = goog.dom.createDom('tr', {},
        page.leftColumn, page.rightColumn);
    contentTable = goog.dom.createDom('table', {}, contentRow);

    page.query = goog.dom.createDom('div', {'id': 'query'}, contentTable);
    goog.dom.appendChild(page.content, page.query);

    page.tabsSetup();

  } else {
    goog.dom.setProperties(page.logo, {'id': 'logo'});
  }
};


/**
 * Insert a new row into a table containing the specified contents
 *
 * @param {DOM} table The table to add the row to.
 * @param {DOM} rowContents The row to be added to the table.
 */
page.insertRow = function(table, rowContents) {
  td = goog.dom.createDom('td', {}, rowContents);
  row = goog.dom.createDom('tr', {}, td);
  goog.dom.appendChild(table, row);
};


/**
 * Set up the page for impact queries.
 */
page.welcomeSetup = function() {
  page.headerSetup();
  page.contentSetup();
  page.footerSetup();
  //page.logoSetup();
  page.querySetup();
};


/**
 * Clear the page content.
 */
page.clearContent = function() {
  goog.dom.removeChildren(page.content);
};


/**
 * Set up the tabs for showing visualizations
 */
page.tabsSetup = function() {

  if (page.map == null) {
    page.map = goog.dom.createDom('div', {'id': 'map'});
    goog.dom.appendChild(page.leftColumn, page.map);
    page.map.style.opacity = 0.0;
  }

  tabsDiv = document.getElementById('visualizationTabs');
  if (tabsDiv == null) {
    tabsDiv = goog.dom.createDom('div', {'id' : 'visualizationTabs'});
    goog.dom.appendChild(page.answer, tabsDiv);

    var tabBar = new goog.ui.TabBar();
    tabBar.render(tabsDiv);
    var tab1 = new goog.ui.Tab('Project Info');
    var tab2 = new goog.ui.Tab('Network Data');
    var tab3 = new goog.ui.Tab('Population');
    var tab4 = new goog.ui.Tab('Income Levels');
    var tab5 = new goog.ui.Tab('JSON Result');
    tabBar.addChild(tab1, true);
    tabBar.addChild(tab2, true);
    tabBar.addChild(tab3, true);
    tabBar.addChild(tab4, true);
    tabBar.addChild(tab5, true);
    tabBar.setSelectedTab(tab1);
    page.setupTabSelections(tabBar);
    goog.dom.appendChild(tabsDiv,
        goog.dom.createDom('div', {'class': 'goog-tab-bar-clear'}));
    page.tabContent = goog.dom.createDom('div',
        {'class': 'goog-tab-content', 'id': 'tabContent'});
    goog.dom.appendChild(tabsDiv, page.tabContent);
    page.setupCurrentTab();
  }

};


/**
 * The label for the currently selected Tab
 */
page.currentTabLabel = 'Project Info';


/**
 * Setup the tabs to load
 *
 * @param {TabBar} tabBar The bar containing the tabs.
 */
page.setupTabSelections = function(tabBar) {

  goog.events.listen(tabBar, goog.ui.Component.EventType.SELECT, function(e) {
    page.currentTabLabel = e.target.getCaption();

    page.setupCurrentTab();
  });

};


/**
 * Calls the setup for the current tab based on the current label
 */
page.setupCurrentTab = function() {
  if (page.currentTabLabel == 'Project Info') {
    page.setupProjectInfoTab();
  }else if (page.currentTabLabel == 'Network Data') {
    page.setupNetworkTab(impact.queryResult);
  }else if (page.currentTabLabel == 'Population') {
    page.setupPopulationTab(impact.queryResult);
  }else if (page.currentTabLabel == 'Income Levels') {
    page.setupIncomeTab(impact.queryResult);
  }else if (page.currentTabLabel == 'JSON Result') {
    page.setupJSONTab(impact.queryResult);
  }
};


/**
 * Setup a description of the project.
 */
page.setupProjectInfoTab = function() {
  goog.dom.removeChildren(page.tabContent);

  header1 = goog.dom.createDom('h2', {'class': 'tabContentText'}, 'Project Overview');
  content1 = goog.dom.createDom('p', {'class': 'tabContentText'}, 'The effects of network properties may have long reaching implications beyond effective file and data transfer.  With the use of the Internet and networked systems in domains such as commerce and social communication, inefficiencies or variations in connections can correspond to major differences in these outside domains. The goal of this project is to unite these network properties with outside domains to better understand the effects they have on each other.');

  header2 = goog.dom.createDom('h2', {'class': 'tabContentText'}, 'Life of a query');
  content2 = goog.dom.createDom('p', {'class': 'tabContentText'}, 'Each query begins using the forms shown on the left.  In this form, you can restrict the client or the server to a specific locality based on a city, region, and/or country, or based on a latitude and longitude.  The client is the location that initiated the tests that collected the network data displayed here.  The server is the location that received and responded to the test.'
  
   content3 = goog.dom.createDom('p', {'class': 'tabContentText'}, 'When you hit submit, we first pass your specifications to Google Geolocator to map your locations and fix any errors in your entries. This means you can use any standard abberviations you would normally use to find a location on Google Maps.  From here, we pass this information to our backend which searches our data sources to find matching tests and census information pertaining to your query.  This  query is sent as GET request, such as the following: "http://mlab-impact.appspot.com/query?cType=cityregioncountry&cRegion=Michigan&cCountry=United+States&sType=cityregioncountry&sCity=New+York&sCounty=New+York&sRegion=New+York&sCountry=United+States .  Because this data is retrieved through a GET request, you also have the option of using only GET requests to retrieve your data.  Of course there will not be any correction of locality information or visualizations, so we recommend this for experienced users only.');
  
   content4 = goog.dom.createDom('p', {'class': 'tabContentText'}, 'When our backend receives the query, we retrieve the appropriate census information if it is available.  We then ask out collection of network tests for the appropriate statistics based on your query.  Because of the large number of tests, this can take a significant amount of time, so we do not always wait for it to complete.  If it completes quickly, we return the result when you first ask us.  However, if there is a longer wait, we send back the census information and tell your end to periodically check with us if the job has completed.  We will generate visualizations for the census information as soon as possible and generate visualizations  for the network data as they become available.');

  header5 = goog.dom.createDom('h2', {'class': 'tabContentText'}, 'Data Sources');
  content5 = goog.dom.createDom('p', {'class': 'tabContentText'}, 'We use data collected from multiple sources, which currently include data collected through measurement lab and the United States Census API, both of which are linked in the footer below.  The census data is comprised of data from the American Community Survey (ACS) and the 2010 Census Summary File 1 (SF1).');

  goog.dom.appendChild(page.tabContent, goog.dom.createDom('br'));
  goog.dom.appendChild(page.tabContent, header1);
  goog.dom.appendChild(page.tabContent, content1);
  goog.dom.appendChild(page.tabContent, goog.dom.createDom('br'));
  goog.dom.appendChild(page.tabContent, header2);
  goog.dom.appendChild(page.tabContent, content2);
  goog.dom.appendChild(page.tabContent, content3);
  goog.dom.appendChild(page.tabContent, content4);
  goog.dom.appendChild(page.tabContent, goog.dom.createDom('br'));
  goog.dom.appendChild(page.tabContent, header5);
  goog.dom.appendChild(page.tabContent, content5);
};


/**
 * Setup the network visualizations.
 *
 * @param {Object} result The result of a query.
 */
page.setupNetworkTab = function(result) {
  goog.dom.removeChildren(page.tabContent);
  if (result == null || result['network data'] == null ||
      result['network data'] == 'loading...') {
    return;
  }

  var options = {};


  var graph1 = goog.dom.createDom('div', {'id': 'graph1', 'class': 'chart'});
  var graph2 = goog.dom.createDom('div', {'id': 'graph2', 'class': 'chart'});
  var graph3 = goog.dom.createDom('div', {'id': 'graph3', 'class': 'chart'});
  var graph4 = goog.dom.createDom('div', {'id': 'graph4', 'class': 'chart'});
  var graph5 = goog.dom.createDom('div', {'id': 'graph5', 'class': 'chart'});
  var graph6 = goog.dom.createDom('div', {'id': 'graph6', 'class': 'chart'});

  table = page.putIn1X3Table(graph1, graph2, graph3);
  table2 = page.putIn1X3Table(graph4, graph5, graph6);
  goog.dom.appendChild(page.tabContent, table);
  goog.dom.appendChild(page.tabContent, table2);

  if (result['network data']['SampleRTT'] != null &&
      result['network data']['SmoothedRTT'] != null &&
      result['network data']['MaxRTT'] != null) {
    var RTTData = [
      [
        parseFloat(result['network data']['SampleRTT']['average']),
        parseFloat(result['network data']['SmoothedRTT']['average']),
        parseFloat(result['network data']['MaxRTT']['average'])
      ]
    ];

    var RTTCategories = ['Sample', 'Smoothed', 'Max'];
    var RTTTypes = ['number', 'number', 'number'];
    var RTTdataObj = results.categorizedDataToObject('',
        ['RTT (ms)'], '', RTTData, RTTCategories, RTTTypes);
    results.columnChart(RTTdataObj, '', true, 225, graph1);
  }

  if (result['network data']['SegsIn'] != null &&
      result['network data']['SegsIn'] != null) {
    var SegData = [
      [
        parseFloat(result['network data']['SegsIn']['average']),
        parseFloat(result['network data']['SegsOut']['average'])
      ]
    ];
    var SegCategories = ['Incoming', 'Outgoing'];
    var SegTypes = ['number', 'number'];
    var SegdataObj = results.categorizedDataToObject('',
        ['Segments'], '', SegData, SegCategories, SegTypes);

    results.columnChart(SegdataObj, '', true, 225, graph2);
  }

  if (result['network data']['CongSignals'] != null &&
          result['network data']['CongAvoid'] != null) {
    var CongData = [
      [
        parseFloat(result['network data']['CongSignals']['average']),
        parseFloat(result['network data']['CongAvoid']['average'])
      ]
    ];
    var CongCategories = ['Congestion Signals', 'Congestion Avoidance'];
    var CongTypes = ['number', 'number'];
    var CongdataObj = results.categorizedDataToObject('',
        ['Congestion'], '', CongData, CongCategories, CongTypes);

    results.columnChart(CongdataObj, '', true, 225, graph3);
  }

  if (result['network data']['CurCwnd']['average'] != null && result['network data']['SmoothedRTT'] != null) {
    var CongData = [
      [
        (parseFloat(result['network data']['CurCwnd']['average'])/parseFloat(result['network data']['SmoothedRTT']['average']))/
      ]
    ];
    var CongCategories = ['Current'];
    var CongTypes = ['number'];
    var CongdataObj = results.categorizedDataToObject('',
        ['Throughput (Mb/s)'], '', CongData, CongCategories, CongTypes);1048.576

    results.columnChart(CongdataObj, '', true, 225, graph5);
  }

  if (result['network data']['DataSegsIn'] != null &&
          result['network data']['DataSegsOut'] != null) {
    var CongData = [
      [
        parseFloat(result['network data']['DataSegsIn']['average']),
        parseFloat(result['network data']['DataSegsOut']['average'])
      ]
    ];
    var CongCategories = ['In', 'Out'];
    var CongTypes = ['number', 'number'];
    var CongdataObj = results.categorizedDataToObject('',
        ['Data Segments'], '', CongData, CongCategories, CongTypes);

    results.columnChart(CongdataObj, '', true, 225, graph4);
  }
};


/**
 * Setup the tab content to display data on population demographics
 *
 * @param {Object} result The result of the query.
 */
page.setupPopulationTab = function(result) {
  goog.dom.removeChildren(page.tabContent);

  if (result == null || result['SF1'] == null) {
    return;
  }

  var options = {};


  var graph1 = goog.dom.createDom('div', {'id': 'graph1', 'class': 'chart'});
  goog.dom.appendChild(page.tabContent, graph1);
  var graph2 = goog.dom.createDom('div', {'id': 'graph2', 'class': 'chart'});
  goog.dom.appendChild(page.tabContent, graph2);

  var demographicCategories = ['Population Count'];
  var demographicTitles = ['American Indian or Alaska Native', 'Asian',
    'Black or African American', 'White', 'Some other Race',
    'Two or more Races'];
  var demographicTypes = ['number'];



  if (result['SF1']['client'] != null &&
      result['SF1']['client']['Asian alone'] != null) {

    data = result['SF1']['client'];
    var demographicData = [
      [parseInt(data['American Indian and Alaska Native alone']['total'])],
      [parseInt(data['Asian alone']['total'])],
      [parseInt(data['Black or African American alone']['total'])],
      [parseInt(data['White alone']['total'])],
      [parseInt(data['Some Other Race alone']['total'])],
      [parseInt(data['Two or More Races']['total'])]
    ];

    var demographicDataObj = results.categorizedDataToObject('Race',
        demographicTitles, '', demographicData,
        demographicCategories, demographicTypes);

    results.pieChart(demographicDataObj,
        'Client Race Demographics', graph1);
  }

  if (result['SF1']['server'] != null &&
      result['SF1']['server']['Asian alone'] != null) {

    data = result['SF1']['server'];
    var demographicData = [
      [parseInt(data['American Indian and Alaska Native alone']['total'])],
      [parseInt(data['Asian alone']['total'])],
      [parseInt(data['Black or African American alone']['total'])],
      [parseInt(data['White alone']['total'])],
      [parseInt(data['Some Other Race alone']['total'])],
      [parseInt(data['Two or More Races']['total'])]
    ];

    var demographicDataObj = results.categorizedDataToObject('Race',
        demographicTitles, '', demographicData, demographicCategories,
        demographicTypes);

    results.pieChart(demographicDataObj, 'Server Race Demographics', graph2);
  }

};


/**
 * Sets up the tab content to display visualizations of income data
 *
 * @param {Object} result The result of the query.
 */
page.setupIncomeTab = function(result) {
  goog.dom.removeChildren(page.tabContent);

  goog.dom.removeChildren(page.tabContent);

  if (result == null || result['ACS'] == null) {
    return;
  }

  var options = {};

  var graph1 = goog.dom.createDom('div', {'id': 'graph1', 'class': 'chart'});
  goog.dom.appendChild(page.tabContent, graph1);
  var graph2 = goog.dom.createDom('div', {'id': 'graph2', 'class': 'chart'});
  goog.dom.appendChild(page.tabContent, graph2);

  var incomeCategories = ['Population Count'];
  var incomeTitles = [
    'Less than $10,000',
    '$10,000 to $14,999',
    '$15,000 to $19,999',
    '$20,000 to $24,999',
    '$25,000 to $29,999',
    '$30,000 to $34,999',
    '$35,000 to $39,999',
    '$40,000 to $44,999',
    '$45,000 to $49,999',
    '$50,000 to $59,999',
    '$60,000 to $74,999',
    '$75,000 to $99,999',
    '$100,000 to $124,999',
    '$125,000 to $149,999',
    '$150,000 to $199,999',
    '$200,000 or more'
  ];
  var incomeTypes = ['number'];

  if (result['ACS']['client'] != null &&
      result['ACS']['client']['Income'] != null) {

    var data = result['ACS']['client']['Income'];
    var incomeData = [
      [parseInt(data['Less than $10,000']['total'])],
      [parseInt(data['$10,000 to $14,999']['total'])],
      [parseInt(data['$15,000 to $19,999']['total'])],
      [parseInt(data['$20,000 to $24,999']['total'])],
      [parseInt(data['$25,000 to $29,999']['total'])],
      [parseInt(data['$30,000 to $34,999']['total'])],
      [parseInt(data['$35,000 to $39,999']['total'])],
      [parseInt(data['$40,000 to $44,999']['total'])],
      [parseInt(data['$45,000 to $49,999']['total'])],
      [parseInt(data['$50,000 to $59,999']['total'])],
      [parseInt(data['$60,000 to $74,999']['total'])],
      [parseInt(data['$75,000 to $99,999']['total'])],
      [parseInt(data['$100,000 to $124,999']['total'])],
      [parseInt(data['$125,000 to $149,999']['total'])],
      [parseInt(data['$150,000 to $199,999']['total'])],
      [parseInt(data['$200,000 or more']['total'])]
    ];

    var incomeDataObj = results.categorizedDataToObject('Income',
        incomeTitles, '', incomeData, incomeCategories, incomeTypes);
    results.columnChart(incomeDataObj, 'Client Incomes', false, 723, graph1);
  }

  if (result['ACS']['server'] != null &&
      result['ACS']['server']['Income'] != null) {

    var data = result['ACS']['server']['Income'];
    var incomeData = [
      [parseInt(data['Less than $10,000']['total'])],
      [parseInt(data['$10,000 to $14,999']['total'])],
      [parseInt(data['$15,000 to $19,999']['total'])],
      [parseInt(data['$20,000 to $24,999']['total'])],
      [parseInt(data['$25,000 to $29,999']['total'])],
      [parseInt(data['$30,000 to $34,999']['total'])],
      [parseInt(data['$35,000 to $39,999']['total'])],
      [parseInt(data['$40,000 to $44,999']['total'])],
      [parseInt(data['$45,000 to $49,999']['total'])],
      [parseInt(data['$50,000 to $59,999']['total'])],
      [parseInt(data['$60,000 to $74,999']['total'])],
      [parseInt(data['$75,000 to $99,999']['total'])],
      [parseInt(data['$100,000 to $124,999']['total'])],
      [parseInt(data['$125,000 to $149,999']['total'])],
      [parseInt(data['$150,000 to $199,999']['total'])],
      [parseInt(data['$200,000 or more']['total'])]
    ];

    var incomeDataObj = results.categorizedDataToObject('Income',
        incomeTitles, '', incomeData, incomeCategories, incomeTypes);
    //results.pieChart(incomeDataObj, 'Server Incomes', graph2);
    results.columnChart(incomeDataObj, 'Server Incomes', false, 723, graph2);
  }
};


/**
 * Setups the content of the tabs to display the JSON representation
 * of the result
 *
 * @param {Object} result The result of the query.
 */
page.setupJSONTab = function(result) {
  goog.dom.removeChildren(page.tabContent);

  if (result == null) {
    return;
  }

  var pre = goog.dom.createDom('pre', {},
            JSON.stringify(result, null, '\t'));
  pre.style.textAlign = 'left';
  goog.dom.appendChild(page.tabContent, pre);
};


/**
 *  Returns a 3 x 1 table containing the three elements
 *
 *  @param {Object} e1 DOM object to be added to the table.
 *  @param {Object} e2 DOM object to be added to the table.
 *  @param {Object} e3 DOM object to be added to the table.
 *
 *  @return {Object} The table containing the 3 values.
 */
page.putIn1X3Table = function(e1, e2, e3) {
  var td1 = goog.dom.createDom('td', {}, e1);
  var td2 = goog.dom.createDom('td', {}, e2);
  var td3 = goog.dom.createDom('td', {}, e3);
  var row = goog.dom.createDom('tr', {}, td1, td2, td3);
  var table = goog.dom.createDom('table', {}, row);
  return table;
};


/**
 * Setup the map for displaying the locations specified in the query
 *
 * @param {Object} result The result of a query made using mlab-impact.
 */
page.setupMap = function(result) {
  if (page.mapCanvas == null) {
    page.mapCanvas = goog.dom.createDom('div', {'id': 'mapCanvas'});
    goog.dom.appendChild(page.map, page.mapCanvas);
  }else if (page.map == null) {
    page.map = goog.dom.createDom('div', {'id': 'map'}, page.mapCanvas);
  }
  page.map.style.opacity = 1.0;

  var clientPos = null;
  var serverPos = null;

  if (result['client']) {
    var clientLoc = result['client'];
    if (clientLoc['lat'] && clientLoc['lng']) {
      var lat = clientLoc['lat'];
      var lng = clientLoc['lng'];
      clientPos = new google.maps.LatLng(lat, lng);
    }
  }

  if (result['server']) {
    var serverLoc = result['server'];
    if (serverLoc['lat'] && serverLoc['lng']) {
      var lat = serverLoc['lat'];
      var lng = serverLoc['lng'];
      serverPos = new google.maps.LatLng(lat, lng);
    }
    serverPos = new google.maps.LatLng(lat, lng);
  }

  var map = mapRenderer.drawMap(page.mapCanvas, clientPos,
      'Client', serverPos, 'Server');
  mapRenderer.adjustMapBoundaries(map, clientPos, serverPos);
};


//Needed for something, but I cannot remember what...oops!
goog.events.listen(window, goog.events.EventType.UNLOAD,
    goog.events.removeAll);

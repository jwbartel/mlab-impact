// Copyright 2011 Google Inc. All Rights Reserved.

/**
 * @fileoverview This is the root module for the net-score project.
 * @author gavaletz@google.com (Eric Gavaletz)
 */


goog.provide('impact');


goog.require('goog.dom');
goog.require('input');
goog.require('page');
goog.require('requester');
goog.require('track');


/**
 * A place to keep track of errors so we can report them.
 *
 * @type {Array.<string>}
 */
impact.error = [];


/**
 * Indicates that the visualization libraries have been loaded and that their
 * dependencies can start.
 *
 * @type {boolean}
 */
//impact.visReady = false;


/**
 * Indicates that the maps API javascript has been successfuly loaded.
 * @type {boolean}
 */
//impact.mapsReady = false;


/**
 * http://code.google.com/apis/loader/
 *
 * This is some tricky stuff and the google.load is not well documented and does
 * some strange things if you are not careful.  Always list the callback
 * function in optional parameters.  See this message for an explination:
 *
 * http://groups.google.com/group/google-maps-api/msg/936506d83fea4458
 */
impact.loadLibs = function() {
  google.load('maps', '3', {'callback': impact.mapsCallback,
    'other_params': 'sensor=' + track.isMobile()});
  google.load('visualization', '1', {'callback': impact.visualizationCallback,
    'packages': ['corechart', 'table']});
};


/**
 * Tracks whther maps were loaded through Google Loader
 *
 * @type {boolean}
 */
impact.mapsLoaded = false;


/**
 * Tracks whther visualizations were loaded through Google Loader
 *
 * @type {boolean}
 */
impact.visualizationsLoaded = false;


/**
 * The callback function when maps have been loaded by Google Loader
 */
impact.mapsCallback = function() {
  impact.mapsLoaded = true;
  if (impact.visualizationsLoaded) {
    page.welcomeSetup();
  }
};


/**
 * The callback function when visualizations have been loaded by Google Loader
 */
impact.visualizationCallback = function() {
  impact.visualizationsLoaded = true;
  if (impact.mapsLoaded = true) {
    page.welcomeSetup();
  }
};


/**
 * The primary entry point for impact.
 */
impact.start = function() {
  //used for all tests to enable tracking via Google analytics.
  track.trackMe(track.trackingCode);

  //used only when we plan on running tests and want to record the user agent
  track.userAgent();

  //dynamically loaded JavaScript that can not be included at compile time
  impact.loadLibs();

};

//used by the closure compiler to eport this function call
goog.exportProperty(window, 'impactStart', impact.start);


/**
 * Send queries through the impact system.
 */
impact.doQuery = function() {

  input.setLoading();
  impact.getGeolocations();

};


/**
 * Ask impact for various data values after the maps have been loaded
 */
impact.askImpact = function() {
  clientStr = impact.parseLocationVals('client', 'c');
  serverStr = impact.parseLocationVals('server', 's');

  url = '/query?' + clientStr + '&' + serverStr;

  listener = function(r) {
    impact.onQuerySuccess(r);
  }
  params = input.getQueryParameters();
  requester.get(url, listener);
};


/**
 *  Parse the location geolocation to pass as a GetQuery based on values
 *  in the query result
 *
 *  @param {string} label The subfield in queryResult that contains the
 *    necessary location.
 *
 *  @param {string} prefix The value 'c' or 's' to specifiy client or
 *    server respectively.
 *
 *  @return {string} The portion of the GET param string representing the
 *    geolocation.
 */
impact.parseLocationVals = function(label, prefix) {
  vals = impact.queryResult[label];
  type = vals['type'];
  str = '';
  if (type != null && type != 'global') {
    str = prefix + 'Type=cityregioncountry';
    if (vals['City'] != null) {
      str += '&' + prefix + 'City=' + vals['City'].replace(' ', '+');
    }
    if (vals['County'] != null) {
      str += '&' + prefix + 'County=' + vals['County'].replace(' ', '+');
    }
    if (vals['State/Region'] != null) {
      str += '&' + prefix + 'Region=' + vals['State/Region'].replace(' ', '+');
    }
    if (vals['Country'] != null) {
      str += '&' + prefix + 'Country=' + vals['Country'].replace(' ', '+');
    }
  }else {
    str = prefix + 'Type=global';
  }
  return str;
};


/**
 * Get the geolocations specified in the form through the Google Maps API
 */
impact.getGeolocations = function() {
  params = input.getQueryParameters();
  impact.queryResult = {};

  listener = function() {
    if (impact.queryResult['client'] != null &&
        impact.queryResult['server'] != null) {

      page.setupMap(impact.queryResult);
      page.setupCurrentTab();
      impact.askImpact();
    }else if (impact.queryResult['err'] != null) {
      input.setAvailable();
    }
  }

  impact.getSingleGeolocation('client', params['client'], listener);
  impact.getSingleGeolocation('server', params['server'], listener);
};


/**
 * Get the single geolocation of a client or user
 *
 * @param {string} src The value client or server.
 * @param {Object} params The object representing the location.
 * @param {function} listener The callback on return.
 */
impact.getSingleGeolocation = function(src, params, listener) {
  if (params['type'] == 'global') {
    impact.queryResult[src] = {'type' : params['type']};
    listener();
  }else {

    geocoder = new google.maps.Geocoder();
    request = {};
    handler = function(results, geoStatus) {
      impact.parseGeoResults(results, geoStatus, src, params['type']);
      listener();
    }

    if (params['type'] == 'latlng') {

      lat = parseFloat(params['latitude']);
      lng = parseFloat(params['longitude']);
      latlng = new google.maps.LatLng(lat, lng);
      request['location'] = latlng;

    }else if (params['type'] == 'cityregioncountry') {

      address = impact.parseAddress(params);
      request['address'] = address;
    }

    geocoder.geocode(request, handler);
  }
};


/**
 * Parse the results returned by the Google Geolocater
 *
 * @param {Object} results The results of the geolocater.
 * @param {string} geoStatus The status code returned by geolocater.
 * @param {string} src Either client or server.
 * @param {string} type The type specified by the user in the form.
 */
impact.parseGeoResults = function(results, geoStatus, src, type) {
  if (geoStatus == 'OK') {
    loc = {};
    if (results != null && results.length > 0) {
      for (var r = 0; r < results.length; r++) {
        result = results[r];
        address = result['address_components'];

        for (var i = 0; i < address.length; i++) {
          label = null;
          type = address[i]['types'][0];
          name = address[i]['long_name'];

          if (type == 'country') {
            label = 'Country';
          } else if (type == 'administrative_area_level_1') {
            label = 'State/Region';
          } else if (type == 'administrative_area_level_2') {
            label = 'County';
          } else if (type == 'locality') {
            label = 'City';
          } else if (type == 'postal_code') {
            label = 'Zip';
          }

          if (label != null && loc[label] == null && name != null) {
            loc[label] = name;
          }

        }

        geometry = result['geometry']['location'];
        loc['lat'] = geometry.lat();
        loc['lng'] = geometry.lng();

      }
    }
    loc['type'] = type;
    impact.queryResult[src] = loc;
    impact.queryResult['network data'] = 'loading...';

  }else {
    impact.queryResult['err'] = type + ' geo: ' + geoStatus;
  }
};


/**
 * Parse the address to send to Google geolocater
 *
 * @param {Object} params The representation of the geolocation.
 *
 * @return {string} An address to pass to Google Geolocater.
 */
impact.parseAddress = function(params) {
  address = '';
  if (params['city'] != '') {
    address += params['city'];
    address += ', ';
  }
  if (params['region'] != '') {
    address += params['region'];
    address += ', ';
  }
  if (params['country'] != '') {
    address += params['country'];
  }

  return address;
};


/**
 * Handle a the result of a successful query
 *
 * @param {Object} response The result of the query.
 */
impact.onQuerySuccess = function(response) {

  response['client'] = impact.queryResult['client'];
  response['server'] = impact.queryResult['server'];
  impact.queryResult = response;

  var url = '';

  var isIncomplete = false;
  if (impact.queryResult['network data'] != null) {
    var isIncomplete = !impact.queryResult['network data']['complete']  &&
        typeof impact.queryResult['err'] == 'undefined';
    if (isIncomplete) {
      url = '/bq_job?jobID=' + impact.queryResult['network data']['jobID'];
    }else {
      input.setAvailable();
    }
    impact.queryResult['network data'] = 'loading...';
  }

  page.setupCurrentTab();
  page.setupMap(impact.queryResult);
  //page.pureDataSetup(impact.queryResult);

  if (isIncomplete) {
    impact.doRequery(url);
  }else{
    input.setAvailable();
  }
};


/**
 * Ask for the result of the BigQuery if the job did not complete.
 *
 * @param {string} url The url to send the request to.
 */
impact.doRequery = function(url) {
  impact.waitTime = 2 * impact.waitTime;
  setTimeout(function() {
    requester.get(url, impact.onRequerySuccess);
  }, impact.waitTime / 2);
};


/**
 * Handle a response from the doRequery()
 *
 * @param {Object} response The result of the doRequery().
 */
impact.onRequerySuccess = function(response) {
  if (response['complete']) {
    impact.queryResult['network data'] = response;
    page.setupCurrentTab();
    input.setAvailable();
  }else {
    url = '/bq_job?jobID=' + response['jobID'];
    impact.doRequery(url);
  }
};


/**
 * The current time to wait between retrying a request for results in seconds
 */
impact.waitTime = 3000;


//used by the closure compiler to eport this function call
goog.exportProperty(window, 'impactSendQuery', impact.sendQuery);


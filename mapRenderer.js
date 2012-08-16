// Copyright 2011 Google Inc. All Rights Reserved.

/**
 * @fileoverview This module provides a way to render maps.
 *
 * @author bartel@cs.unc.edu (Jacob Bartel)
 */


goog.provide('mapRenderer');

/**
 * Setup the map for displaying the locations specified in the query
 *
 * @param {Object} result The result of a query made using mlab-impact.
 */
mapRenderer.drawMap = function(parentDiv, pos1, label1, pos2, label2) {
  var defaultCenter = new google.maps.LatLng(35.9131, -79.0561);

  var mapOptions = {
    center: defaultCenter,
    zoom: 1,
    mapTypeId: google.maps.MapTypeId.ROADMAP
  };
  var map = new google.maps.Map(parentDiv, mapOptions);

  var clientPos = null;
  var serverPos = null;

  if (pos1 != null) {
    var marker = new google.maps.Marker({
      position: pos1,
      title: label1,
      map: map
    });
  }

  if (pos2 != null) {
    var marker = new google.maps.Marker({
      position: pos2,
      title: label2,
      map: map
    });
  }

  return map;
};


/**
 * Adjust the boundaries for a map for a client and server
 *
 * @param {google.maps.Map} map The map which needs to be adjusted.
 * @param {google.maps.LatLng} clientPos Position of the client.
 * @param {google.maps.LatLng} serverPos Position of the server.
 */
mapRenderer.adjustMapBoundaries = function(map, clientPos, serverPos) {
  if (clientPos == null || serverPos == null) {
    return;
  }

  if (clientPos == null) {
    map.setZoom(8);
    map.setCenter(serverPos);
    return;
  }

  if (serverPos == null) {
    map.setZoom(8);
    map.setCenter(clientPos);
    return;
  }

  var border = 0;

  var southLat = Math.min(clientPos.lat(), serverPos.lat()) - border;
  var northLat = Math.max(clientPos.lat(), serverPos.lat()) + border;
  var westLng = 0;
  var eastLng = 0;

  serverToClient = mapRenderer.lngDistance(clientPos.lng(), serverPos.lng());
  clientToServer = mapRenderer.lngDistance(serverPos.lng(), clientPos.lng());

  if (clientToServer <= serverToClient) {
    eastLng = clientPos.lng() + border;
    westLng = serverPos.lng() - border;
  }else {
    eastLng = serverPos.lng() + border;
    westLng = clientPos.lng() - border;
  }

  sw = new google.maps.LatLng(southLat, westLng);
  ne = new google.maps.LatLng(northLat, eastLng);
  map.setZoom(8);
  map.fitBounds(new google.maps.LatLngBounds(sw, ne));
};


/**
 * Find distance between two longitude values in degrees going west to east
 *
 * @param {float} start The starting point of the distance.
 * @param {float} end The end point for the distance.
 *
 * @return {float} The distance between start and end in degrees.
 */
mapRenderer.lngDistance = function(start, end) {
  if (end >= start) {
    return Math.abs(end - start);
  }else {
    return 360 - mapRenderer.lngDistance(end, start);
  }
};





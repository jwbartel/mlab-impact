// Copyright 2011 Google Inc. All Rights Reserved.

/**
 * @fileoverview Description of this file.
 * @author gavaletz@google.com (Eric Gavaletz)
 */


goog.provide('track');


goog.require('goog.Uri');
goog.require('goog.userAgent');
goog.require('goog.userAgent.flash');
goog.require('goog.userAgent.product');
goog.require('goog.userAgent.product.isVersion');


/**
 * The tracking code to be used with Google Analytics.
 *
 * @type {string}
 */
track.trackingCode = 'UA-16051354-8';


/**
 * Sets the tracking code.  Usually just used from the server.
 *
 * @param {string} s the tracking code to be used.
 */
track.setTrackingCode = function(s) {
  track.trackingCode = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'trackSetTrackingCode', track.setTrackingCode);


/**
 * The remote address that made the request.
 *
 * @type {string}
 */
track.remoteAddr = '';


/**
 * Sets the remote address.  Usually just used from the server.
 *
 * @param {string} s the remote address to be used.
 */
track.setRemoteAddr = function(s) {
  track.remoteAddr = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'trackSetRemoteAddr', track.setRemoteAddr);


/**
 * The refereing site for the request.
 *
 * @type {string}
 */
track.referer = '';


/**
 * Sets the refereing site.  Usually just used from the server.
 *
 * @param {string} s the refereing site to be used.
 */
track.setReferer = function(s) {
  track.referer = s;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'trackSetReferer', track.setReferer);


/**
 * The server-side time of the request (UTC seconds).
 *
 * @type {number}
 */
track.requestTime = 0;


/**
 * Sets the time of the request.  Usually just used from the server.
 *
 * @param {number} n the time of the request to be used.
 */
track.setRequestTime = function(n) {
  track.requestTime = n;
};
//used by the closure compiler to eport this function call
goog.exportProperty(window, 'trackSetRequestTime', track.setRequestTime);


/**
 * Tracking code from the analytics site.
 *
 * @type {Array.{*}}
 */
var _gaq = _gaq || [];


/**
 * Do the Google analytics stuff.
 *
 * @param {string} trackingCode from the analytics site.
 */
track.trackMe = function(trackingCode) {
  _gaq.push(['_setAccount', trackingCode]);
  _gaq.push(['_setDomainName', 'net-score.org']);
  _gaq.push(['_setAllowLinker', true]);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script');
    ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ?
        'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(ga, s);
  })();
};


/**
 * Track the user agent stuff.
 *
 * This is a well kept secret...
 * http://closure-library.googlecode.com/svn/docs/closure_goog_useragent_useragent.js.html
 */
track.userAgent = function() {
  //product
  track.productVersion = goog.userAgent.product.VERSION;
  if (goog.userAgent.product.CHROME) {
    track.product = 'CHROME';
  } else if (goog.userAgent.product.FIREFOX) {
    track.product = 'FIREFOX';
  } else if (goog.userAgent.product.IE) {
    track.product = 'IE';
  } else if (goog.userAgent.product.ANDROID) {
    track.product = 'ANDROID';
  } else if (goog.userAgent.product.IPAD) {
    track.product = 'IPAD';
  } else if (goog.userAgent.product.IPHONE) {
    track.product = 'IPHONE';
  } else if (goog.userAgent.product.SAFARI) {
    track.product = 'SAFARI';
  } else if (goog.userAgent.product.OPERA) {
    track.product = 'OPERA';
  } else if (goog.userAgent.product.CAMINO) {
    track.product = 'CAMINO';
  } else {
    track.product = 'OTHER';
  }

  //platform
  track.platformVersion = goog.userAgent.PLATFORM;
  if (goog.userAgent.LINUX) {
    track.platform = 'LINUX';
  } else if (goog.userAgent.MAC) {
    if (track.product == 'IPAD' || track.product == 'IPHONE') {
      track.platform = 'IOS';
    }
    else {
      track.platform = 'MAC';
    }
  } else if (goog.userAgent.WINDOWS) {
    track.platform = 'WINDOWS';
  } else {
    track.platform = track.platformVersion;
  }

  //flash
  track.flash = goog.userAgent.flash.HAS_FLASH;
  track.flashVersion = goog.userAgent.flash.VERSION;

  //renderer
  var renderer;
  track.rendererVersion = goog.userAgent.VERSION;
  if (goog.userAgent.WEBKIT) {
    track.renderer = 'WEBKIT';
  } else if (goog.userAgent.GECKO) {
    track.renderer = 'GECKO';
  } else if (goog.userAgent.IE) {
    track.renderer = 'IE';
  } else if (goog.userAgent.OPERA) {
    track.renderer = 'OPERA';
  } else {
    track.renderer = 'OTHER';
  }

  //mobile
  track.mobile = goog.userAgent.MOBILE;
};


/**
 * Is the remote device iOS?
 *
 * @return {boolean} Is the aganet an iOS device?
 */
track.isIos = function() {
  return (goog.userAgent.product.IPAD || goog.userAgent.product.IPHONE);
};


/**
 * Is the remote device mobile?
 *
 * @return {boolean} Is the aganet a mobile device?
 */
track.isMobile = function() {
  return goog.userAgent.MOBILE;
};


/**
 * Does the remote device have flash?
 *
 * @return {boolean} Does the remote device have flash?
 */
track.hasFlash = function() {
  return goog.userAgent.flash.HAS_FLASH;
};

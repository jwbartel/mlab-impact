
/**
 * @fileoverview This is the module that gets results for the visualizer.
 *
 * @author gavaletz@google.com (Eric Gavaletz)
 */


goog.provide('requester');


goog.require('goog.events');
goog.require('goog.json');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');


/**
 * Get the data to the given uri
 *
 * @param {string} url The url to send the request to.
 * @param {function=} opt_response_handler Function that handles the response.
 */
requester.get = function(url, opt_response_handler) {
  xhr = new goog.net.XhrIo();
  if (goog.isDefAndNotNull(opt_response_handler)) {
    responder = function(e) {
      response = xhr.getResponseJson();
      opt_response_handler(response);
    }
    goog.events.listen(xhr, goog.net.EventType.COMPLETE,
        responder, false);
  }
  xhr.send(url);
};

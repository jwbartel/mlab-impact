
/**
 * @fileoverview This module provicdes a way to display our site.
 *
 * @author bartel@cs.unc.edu (Jacob Bartel)
 */

goog.provide('input');

goog.require('goog.dom');
goog.require('goog.ui.FlatButtonRenderer');
goog.require('goog.ui.Select');


/**
 * marks whether the elements of the form are disabled
 */
input.formDisabled = false;


/**
 * Set the form as loading data, and thus disabled
 */
input.setLoading = function() {
  input.formDisabled = true;
  page.spinner.style.display = 'inline';
  if (typeof input.queryForm != 'undefined') {
    input.queryForm.style.opacity = 0.5;
    input.queryForm.style.filter = 'alpha(opacity=0.5)';
    clientType = input.clientQueryType.value;
    if (clientType == 'latlng') {
      document.getElementById('clientLat').disabled = true;
      document.getElementById('clientLng').disabled = true;
    }else if (clientType == 'cityregioncountry') {
      document.getElementById('clientCountry').disabled = true;
      document.getElementById('clientRegion').disabled = true;
      document.getElementById('clientCity').disabled = true;
    }

    serverType = input.serverQueryType.value;
    if (serverType == 'latlng') {
      document.getElementById('serverLat').disabled = true;
      document.getElementById('serverLng').disabled = true;
    }else if (serverType == 'cityregioncountry') {
      document.getElementById('serverCountry').disabled = true;
      document.getElementById('serverRegion').disabled = true;
      document.getElementById('serverCity').disabled = true;
    }
  }
};


/**
 * Set the form as available to make a query
 */
input.setAvailable = function() {
  input.formDisabled = false;
  page.spinner.style.display = 'none';
  if (typeof input.queryForm != 'undefined') {
    input.queryForm.style.opacity = 1.0;
    input.queryForm.style.filter = 'alpha(opacity=1.0)';
    clientType = input.clientQueryType.value;
    if (clientType == 'latlng') {
      document.getElementById('clientLat').disabled = false;
      document.getElementById('clientLng').disabled = false;
    }else if (clientType == 'cityregioncountry') {
      document.getElementById('clientCountry').disabled = false;
      document.getElementById('clientRegion').disabled = false;
      document.getElementById('clientCity').disabled = false;
    }

    serverType = input.serverQueryType.value;
    if (serverType == 'latlng') {
      document.getElementById('serverLat').disabled = false;
      document.getElementById('serverLng').disabled = false;
    }else if (serverType == 'cityregioncountry') {
      document.getElementById('serverCountry').disabled = false;
      document.getElementById('serverRegion').disabled = false;
      document.getElementById('serverCity').disabled = false;
    }
  }
};


/**
 * Get the parameters to be passed to the querying service
 *
 * @return {Object} The collected specifications of the geolocations as
 *   given by the user in the form.
 */
input.getQueryParameters = function() {
  params = {};
  if (typeof input.queryForm != 'undefined') {
    client = {};
    client['type'] = input.clientQueryType.value;
    if (client['type'] == 'latlng') {
      client['latitude'] = document.getElementById('clientLat').value;
      client['longitude'] = document.getElementById('clientLng').value;
    }else if (client['type'] == 'cityregioncountry') {
      client['country'] = document.getElementById('clientCountry').value;
      client['region'] = document.getElementById('clientRegion').value;
      client['city'] = document.getElementById('clientCity').value;
    }
    params['client'] = client;

    server = {};
    server['type'] = input.serverQueryType.value;
    if (server['type'] == 'latlng') {
      server['latitude'] = document.getElementById('serverLat').value;
      server['longitude'] = document.getElementById('serverLng').value;
    }else if (server['type'] == 'cityregioncountry') {
      server['country'] = document.getElementById('serverCountry').value;
      server['region'] = document.getElementById('serverRegion').value;
      server['city'] = document.getElementById('serverCity').value;
    }
    params['server'] = server;
  }

  return params;
};


/**
 * Setup the form for inputting a query
 *
 * @param {DOM} parent The element holding the form.
 */
input.queryFormSetup = function(parent) {
  if (input.queryForm == null || typeof input.queryForm == 'undefined') {
    input.queryForm = goog.dom.createDom('form',
        {'id': 'questionForm'});

    queryElements = goog.dom.createDom('table', {'id': 'queryElements'});

    input.clientQueryType = goog.dom.createDom('input',
        {'type': 'hidden', 'id': 'clientType',
          'name': 'clientType', 'value': 'global'});
    clientSearchFields = goog.dom.createDom('table',
        {'class' : 'searchFields'});
    clientQueryTypeSelector =
        input.searchSelectorSetup('Filter client by',
        clientSearchFields, input.clientQueryType, 'client');

    goog.dom.appendChild(input.queryForm, clientQueryTypeSelector);
    goog.dom.appendChild(input.queryForm, input.clientQueryType);
    goog.dom.appendChild(input.queryForm, clientSearchFields);

    input.serverQueryType = goog.dom.createDom('input',
        {'type': 'hidden', 'id': 'serverType',
          'name': 'serverType', 'value': 'global'});
    serverSearchFields =
        goog.dom.createDom('table', {'class': 'searchFields'});
    serverQueryTypeSelector =
        input.searchSelectorSetup('Filter server by',
        serverSearchFields, input.serverQueryType, 'server');

    goog.dom.appendChild(input.queryForm, serverQueryTypeSelector);
    goog.dom.appendChild(input.queryForm, input.serverQueryType);
    goog.dom.appendChild(input.queryForm, serverSearchFields);

    submitButton = new goog.ui.Button('Submit',
        goog.ui.FlatButtonRenderer.getInstance());

    goog.dom.appendChild(input.queryForm, queryElements);
    submitButton.render(input.queryForm);
    submitButton.getElement().id = 'submitButton';

    goog.dom.appendChild(parent, input.queryForm);


    goog.events.listen(submitButton.getElement(),
        goog.events.EventType.CLICK, function(e) {
          if (!input.formDisabled) {
            impact.doQuery();
          }
        });

  }
};


/**
 * Setup the selector for how to search
 *
 * @param {string} title The initial label of the selector.
 * @param {DOM} searchFieldDiv Where to put search fields.
 * @param {DOM} queryType The hidden input that contains the type of
 *   geolocation used.
 * @param {string} idPrefix The prefix used for all ids.
 *
 * @return {DOM} The wrapping div that contains the selector.
 */
input.searchSelectorSetup = function(title, searchFieldDiv,
    queryType, idPrefix) {
  var selector = new goog.ui.Select('All Locations', null,
      goog.ui.FlatButtonRenderer.getInstance(), null);

  input.buildOption(selector, 'All Locations', input.buildGlobalForm,
      searchFieldDiv, queryType, idPrefix);
  selector.addItem(new goog.ui.MenuSeparator());
  input.buildOption(selector, 'Latitude/Longitude', input.buildLatLongForm,
      searchFieldDiv, queryType, idPrefix, idPrefix);
  selector.addItem(new goog.ui.MenuSeparator());
  input.buildOption(selector, 'City/Region/Country',
      input.buildCityRegionCountryForm, searchFieldDiv, queryType, idPrefix);

  var row = goog.dom.createDom('tr');
  var title = goog.dom.createDom('td',
      {'class': 'selectorTitle'}, title);
  var selector_spot = goog.dom.createDom('td');
  var wrapper = goog.dom.createDom('table', {}, row);
  goog.dom.appendChild(row, title);
  goog.dom.appendChild(row, selector_spot);
  selector.render(selector_spot);


  return wrapper;
};


/**
 * Build option an option for the selector
 *
 * @param {DOM} selectorParent The selector that will contain the option.
 * @param {string} label  What to be displayed for the option.
 * @param {function} action What to do when the option is selected.
 * @param {DOM} searchFieldDiv Where search fields should be created.
 * @param {DOM} queryType The hidden input that should specify the type
 *   of geolocation used.
 * @param {string} idPrefix The prefix for each id.
 */
input.buildOption = function(selectorParent, label, action,
    searchFieldDiv, queryType, idPrefix) {

  var option = new goog.ui.Button(label,
      goog.ui.FlatMenuButtonRenderer.getInstance());

  selectorParent.addItem(option);
  goog.events.listen(option.getElement(),
      goog.events.EventType.CLICK, function(e) {
        if (!input.formDisabled) {
          action(searchFieldDiv, queryType, idPrefix);
        }
      });
};


/**
 * Build the form for querying by latitude longitude
 *
 * @param {DOM} searchFieldDiv Where search fields should be created.
 * @param {DOM} queryType The hidden input that should specify the type
 *   of geolocation used.
 */
input.buildGlobalForm = function(searchFieldDiv, queryType) {
  if (queryType.value == 'global') {
    return;
  }

  goog.dom.removeChildren(searchFieldDiv);
  queryType.value = 'global';
};


/**
 * Build the form for querying by latitude longitude
 *
 * @param {DOM} searchFieldDiv Where search fields should be created.
 * @param {DOM} queryType The hidden input that should specify the type
 *   of geolocation used.
 * @param {string} idPrefix The prefix for each id.
 */
input.buildLatLongForm = function(searchFieldDiv, queryType, idPrefix) {
  if (queryType.value == 'latlng') {
    return;
  }

  goog.dom.removeChildren(searchFieldDiv);

  queryType.value = 'latlng';
  input.queryElementSetup(searchFieldDiv, 'Latitude', idPrefix + 'Lat');
  input.queryElementSetup(searchFieldDiv, 'Longitude', idPrefix + 'Lng');
};


/**
 * Build the form for querying by City, State, and Country
 *
 * @param {DOM} searchFieldDiv Where search fields should be created.
 * @param {DOM} queryType The hidden input that should specify the type
 *   of geolocation used.
 * @param {string} idPrefix The prefix for each id.
 */
input.buildCityRegionCountryForm = function(searchFieldDiv,
    queryType, idPrefix) {

  if (queryType.value == 'cityregioncountry') {
    return;
  }

  goog.dom.removeChildren(searchFieldDiv);

  queryType.value = 'cityregioncountry';
  input.queryElementSetup(searchFieldDiv, 'City', idPrefix + 'City');
  input.queryElementSetup(searchFieldDiv, 'State/Region', idPrefix + 'Region');
  input.queryElementSetup(searchFieldDiv, 'Country', idPrefix + 'Country');
};


/**
 * Setup an element of the query form in a table of query elements
 *
 * @param {dom element} parentTable The table containing all visible
 *                                  query elements.
 * @param {string} name The name shown when the table is displayed.
 * @param {string} id The prefix used for ids in the label and input in
 *                    the form.
 */
input.queryElementSetup = function(parentTable, name, id) {
  var title = goog.dom.createDom('td', {'class': 'entryTitle'}, name);
  var entry = goog.dom.createDom('input', {'type': 'text', 'id': id});
  entry.disabled = input.formDisabled;
  var entryWrapper = goog.dom.createDom('td', {'class': 'entryInput'}, entry);

  var row = goog.dom.createDom('tr', {});
  goog.dom.appendChild(row, title);
  goog.dom.appendChild(row, entryWrapper);

  goog.dom.appendChild(parentTable, row);
};

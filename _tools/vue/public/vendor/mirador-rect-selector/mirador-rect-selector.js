// mirador-rect-selector.js
// highlight rectangle region via url parameters

(function($, jQuery, OpenSeadragon) {
  /* ViewFromUrl (UpdateUrlFromView) plugin or like should be installed */
  if(typeof UpdateUrlFromView === 'undefined') {
    console.info('ViewFromUrl (UpdateUrlFromView) plugin or like '
      + 'should be installed.');
  }

  /* local variables */
  const _data = {
    initialized: false,
  };

  /* init options */
  const _options = {
    annoPosition: 'center', // see _availableAnnoPositions
    annoVertical: false, // makes annotation vertical if true
    zoomMargin: 300, // offset in zooming
    zoom: 'on', // regards anything but 'off' as 'on'
  };

  const _availableAnnoPositions = [
    'center', // default
    'left',
    'right',
    'upper',
    'lower',
  ];

  /* url query parameters */
  const _availableParams = [
    'canvas', // canvas=<canvas id>
    'xywh', // xywh=<x>,<y>,<width>,<height>
    'chars', // chars=<annotation>
    'zoom', // regards anythin but 'off' as 'on'
  ];

  // utilily
  function parseQueryString() {
    window.location.search.replace(/^\?/, '').split('&').forEach((e) => {
      const kv = e.split('=');
      if(_availableParams.includes(kv[0])) {
        _data[decodeURIComponent(kv[0])] = decodeURIComponent(kv[1] || '');
      }
    });
    _data.initialized = true;
  }

  /* originally defined in mirador/js/src/utils/saveController.js */
  const origFuncSaveControllerInit = $.SaveController.prototype.init;
  $.SaveController.prototype.init = function(config) {
    origFuncSaveControllerInit.apply(this, [config]);

    if(config && config.rectSelector) {
      jQuery.extend(_options, config.rectSelector);
    }
  };

  const origImageViewSyncAllImageResourceProperties = $.ImageView.prototype
    .syncAllImageResourceProperties;
  $.ImageView.prototype.syncAllImageResourceProperties = function(...args) {
    origImageViewSyncAllImageResourceProperties.apply(this, args);

    // only for the first time
    if(_data.initialized) {
      return;
    }

    parseQueryString();
    if(_data.canvas === this.canvasID
      && _data.xywh && _data.xywh.match(/^(\d+),(\d+),(\d+),(\d+)$/)) {
      const xywh = _data.xywh.split(',').map(v => parseInt(v, 10));
      _data.bounds = new (Function.prototype.bind
        .apply(OpenSeadragon.Rect, [null].concat(xywh)));

      _data.elem = document.createElement('div');
      _data.elem.id = 'rect-selector-overlay';
      this.osd.addOverlay({
        element: _data.elem,
        location: _data.bounds,
      });

      // zoom
      if((typeof _data.zoom !== 'undefined' && _data.zoom !== 'off')
        || (_options.zoom !== 'off' && typeof _data.zoom === 'undefined')) {
        xywh[0] -= _options.zoomMargin;
        xywh[1] -= _options.zoomMargin;
        xywh[2] += _options.zoomMargin * 2;
        xywh[3] += _options.zoomMargin * 2;
        const margined = new (Function.prototype.bind
          .apply(OpenSeadragon.Rect, [null].concat(xywh)));
        this.osd.viewport.fitBounds(margined);
      }

      // annotation (chars)
      if(_data.chars) {
        // container
        _data.anno = document.createElement('div');
        _data.anno.id = 'rect-selector-anno-overlay';
        _data.anno.classList.add(`rect-selector-anno-overlay-${
          _availableAnnoPositions.includes(_options.annoPosition)
            ? _options.annoPosition : 'center'}`);
        // chars
        const div = document.createElement('div');
        div.textContent = _data.chars;
        if(_options.annoVertical) {
          div.classList.add('vertical-rl');
        }
        _data.anno.appendChild(div);
        this.osd.addOverlay({
          element: _data.anno,
          location: _data.bounds,
        });
      }
    }
  };

  const origImageViewUpdateImage = $.ImageView.prototype.updateImage;
  $.ImageView.prototype.updateImage = function(...args) {
    origImageViewUpdateImage.apply(this, args);

    // the second time and following
    if(this.canvasID === _data.canvas) {
      this.osd.addOverlay({
        element: _data.elem,
        location: _data.bounds,
      });
      this.osd.addOverlay({
        element: _data.anno,
        location: _data.bounds,
      });
    } else {
      this.osd.removeOverlay(_data.elem);
      this.osd.removeOverlay(_data.anno);
    }
  };
})(Mirador, jQuery, OpenSeadragon);

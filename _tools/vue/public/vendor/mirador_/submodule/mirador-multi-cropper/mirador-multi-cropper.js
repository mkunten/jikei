// mirador-multi-cropper.js

(function($, $jq) {
  // options
  const _options = {
    activeOnStart: false,
    multiCropperCls: 'multi-cropper-canvas',
    multiCropperState: 'off',
  };

  // local variables
  let _active = false; // dummy
  let _rects = [];

  // templates
  const _templates = {
    button: $.Handlebars.compile([
      '<div class="multi-cropper-controls{{#if active}} active{{/if}}">',
      '<a class="multi-cropper-layer hud-control" role="button" title="{{t "toggle-multi-cropper"}}" aria-label="{{t "toggle-multi-cropper"}}">',
      '<i class="fa fa-lg fa-clone"></i>',
      '</a>',
      '</div>',
    ].join('')),
    controls: $.Handlebars.compile([
      '<span class="multi-cropper-menu">',
      '<a class="multi-cropper-pointer-mode hud-control selected" title="{{t "pointerTooltip"}}">',
      '<i class="fa fa-mouse-pointer"></i>',
      '</a>',
      '<a class="multi-cropper-check_box_outline_blank-mode hud-control multi-cropper-edit-mode" title="{{t "rectangleTooltip"}}">',
      '<i class="material-icons">check_box_outline_blank</i>',
      '</a>',
      '</span>',
    ].join('')),
  };

  // custom: settings
  $jq.extend($.DEFAULT_SETTINGS, {
    multiCropper: _options,
  });

  // /* custom: widget */
  // custom: contextControls
  const origContextControlsPrototypeInit = $.ContextControls.prototype.init;
  $.ContextControls.prototype.init = function() {
    origContextControlsPrototypeInit.call(this);

    if(_options.activeOnStart) {
      _active = _options.activeOnStart;
    }
    const controls = $jq(_templates.button({
      active: _active,
    }));
    this.multiCropperElement = $jq(_templates.controls()).hide()
      .appendTo(controls);
    controls.appendTo(this.container);
  };

  $.ContextControls.prototype.multiCropperShow = function() {
    this.multiCropperElement.fadeIn('150');
  };

  $.ContextControls.prototype.multiCropperHide = function() {
    this.multiCropperElement.fadeOut('150');
  };

  // custom: hud
  const origHudPrototypeCreateStateMachines = $.Hud.prototype.createStateMachines;
  $.Hud.prototype.createStateMachines = function() {
    origHudPrototypeCreateStateMachines.call(this);

    const _this = this;
    const duration = 200;

    this.multiCropperStateMachine = StateMachine.create({
      events: [
        { name: 'startup', from: 'none', to: 'off' },
        { name: 'displayOn', from: 'off', to: 'pointer'},
        { name: 'displayOff', from: ['pointer', 'rect'], to: 'off'},
        { name: 'choosePointer', from: ['pointer', 'rect'], to: 'pointer'},
        { name: 'chooseRect', from: 'pointer', to: 'rect'},
        { name: 'refresh', from: 'pointer', to: 'pointer'},
        { name: 'refresh', from: 'rect', to: 'rect'}
      ],
      callbacks: {
        onstartup: function(event, from, to) {
          console.log(1, from, to);
          _this.eventEmitter.publish(('windowUpdated'), {
            id: _this.windowId,
            multiCropperStateMachine: to
          });
        },
        ondisplayOn: function(event, from, to) {
          console.log(2, from, to);
          _this.eventEmitter.publish('HUD_ADD_CLASS.'+_this.windowId, ['.multi-cropper-layer', 'selected']);
          _this.contextControls.multiCropperShow();
          _this.eventEmitter.publish('multiCropperModeChange.' + _this.windowId, 'displayMultiCropper');
          _this.eventEmitter.publish('HUD_ADD_CLASS.'+_this.windowId, ['.multi-cropper-pointer-mode', 'selected']);
          _this.eventEmitter.publish('DEFAULT_CURSOR.' + _this.windowId);
          _this.eventEmitter.publish(('windowUpdated'), {
            id: _this.windowId,
            multiCropperStateMachine: to
          });
        },
        ondisplayOff: function(event, from, to) {
          console.log(3, from, to);
          // if (_this.annoEndpointAvailable) {
          //   _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-edit-mode', 'selected']);
          //   _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-pointer-mode', 'selected']);
          //   _this.eventEmitter.publish('CANCEL_ACTIVE_MULTI_CROPPER.'+_this.windowId);
          //   _this.contextControls.multiCropperHide();
          // }
          _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-layer', 'selected']);
          _this.eventEmitter.publish('multiCropperModeChange.' + _this.windowId, 'default');
          _this.eventEmitter.publish(('windowUpdated'), {
            id: _this.windowId,
            multiCropperStateMachine: to
          });
        },
        onchoosePointer: function(event, from, to) {
          console.log(4, from, to);
          _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-edit-mode', 'selected']);
          _this.eventEmitter.publish('HUD_ADD_CLASS.'+_this.windowId, ['.multi-cropper-pointer-mode', 'selected']);
          _this.eventEmitter.publish('multiCropperModeChange.' + _this.windowId, 'displayMultiCroppers');
          _this.eventEmitter.publish('DEFAULT_CURSOR.' + _this.windowId);
          _this.eventEmitter.publish(('windowUpdated'), {
            id: _this.windowId,
            multiCropperStateMachine: to
          });
        },
        onchooseRect: function(event, from, to, shape) {
          console.log(5, from, to);
          _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-pointer-mode', 'selected']);
          _this.eventEmitter.publish('HUD_REMOVE_CLASS.'+_this.windowId, ['.multi-cropper-edit-mode', 'selected']);
          _this.eventEmitter.publish('HUD_ADD_CLASS.'+_this.windowId, ['.multi-cropper-check_box_outline_blank-mode', 'selected']);
          _this.eventEmitter.publish('multiCropperModeChange.' + _this.windowId, 'creatingMultiCropper');
          _this.eventEmitter.publish('CROSSHAIR_CURSOR.' + _this.windowId);
          _this.eventEmitter.publish('toggleDrawingTool.'+_this.windowId, shape);

          _this.eventEmitter.publish(('windowUpdated'), {
            id: _this.windowId,
            multiCropperStateMachine: to
          });
        },
        onrefresh: function(event, from, to) {
          console.log(6, from, to);
        }
      }
    });
  };

  const origHudPrototypeListenForActions = $.Hud.listenForActions;
  $.Hud.listenForActions = function() {
    origHudPrototypeListenForActions.call(this);

    const _this = this;
    this.eventEmitter.subscribe('SET_MULTI_CROPPER_STATE_MACHINE_POINTER.' + this.windowId, function(event) {
      if (_this.multiCropperStateMachine.current === 'none') {
        _this.multiCropperStateMachine.startup();
      } else if (_this.multiCropperStateMachine.current === 'off') {
        _this.multiCropperStateMachine.displayOn();
      } else {
        _this.multiCropperStateMachine.choosePointer();
      }
    });
  };

  // custom: imageView
  const origImageViewPrototypeInit = $.ImageView.prototype.init;
  $.ImageView.prototype.init = function() {
    origImageViewPrototypeInit.call(this);

    this.elemMultiCropper = $jq('<div />')
      .addClass(_options.multiCropperCls)
      .appendTo(this.element);
  };

  const origImageViewPrototypeBindEvents = $.ImageView.prototype.bindEvents;
  $.ImageView.prototype.bindEvents = function() {
    origImageViewPrototypeBindEvents.call(this);

    const _this = this;

    this.element.find('.multi-cropper-layer').on('click', function() {
      console.log(_this.hud.multiCropperStateMachine.current);
      if (_this.hud.multiCropperStateMachine.current === 'none') {
        _this.hud.multiCropperStateMachine.startup(this);
      }
      if (_this.hud.multiCropperStateMachine.current === 'off') {
        _this.hud.multiCropperStateMachine.displayOn(this);
        _this.multiCropperState = 'on';
      } else {
        // make sure to force the controls back to auto fade
        _this.forceShowControls = false;
        _this.hud.multiCropperStateMachine.displayOff(this);
        _this.multiCropperState = 'off';
      }
    });

    this.element.find('.multi-cropper-pointer-mode').on('click', function() {
      console.log('c1 pointer');
      // go back to pointer mode
      if (_this.hud.multiCropperStateMachine.current === 'rect') {
        _this.hud.multiCropperStateMachine.choosePointer();
        //go back to allowing the controls to auto fade
        _this.forceShowControls = false;
      }
    });

    this.element.find('.multi-cropper-check_box_outline_blank-mode').on('click', function() {
      console.log('c2 rect');
      if (_this.hud.multiCropperStateMachine.current === 'pointer') {
        _this.hud.multiCropperStateMachine.chooseRect();
      }
      //when a user is in Create mode, don't let the controls auto fade as it could be distracting to the user
      _this.forceShowControls = true;
      _this.element.find(".hud-control").stop(true, true).removeClass('hidden', _this.state.getStateProperty('fadeDuration'));
    });
  };

  const origImageViewPrototypeInitialiseImageCanvas = $.ImageView.prototype.initialiseImageCanvas;
  $.ImageView.prototype.initialiseImageCanvas = function() {
    origImageViewPrototypeInitialiseImageCanvas.call(this);

    // this.addMultiCropperLayer(this.elemMultiCropper);

    const originalState = this.hud.multiCropperStateMachine.current;
    if(originalState === 'none') {
      this.hud.multiCropperStateMachine.startup();
    } else if (originalState === 'off' || this.multiCropperState === 'off') {
      // original state is off, so don't need to do anything
    } else {
      this.hud.multiCropperStateMachine.displayOff();
    }

    if (originalState === 'pointer' || this.multiCropperState === 'on') {
      this.hud.multiCropperStateMachine.displayOn();
    } else if (originalState === 'rect') {
      this.hud.multiCropperStateMachine.displayOn();
      this.hud.multiCropperStateMachine.chooseRect();
    } else {
      // original state is off, so don't need to do anything
    }
  };

  $.ImageView.prototype.addMultiCropperLayer = function(element) {
    const _this = this;
    _this.multiCropperLayer = new $.MultiCropperLayer({
      state: _this.state,
      multiCropperList: _this.state.getWindowMultiCropperList(_this.windowId)
        || [],
      viewer: _this.osd,
      windowId: _this.windowId,
      element: element,
      eventEmitter: _this.eventEmitter
    });
  };

  OpenSeadragon.Viewer.prototype.multiCropperOverlay = function(osdViewerId, windowId, state, eventEmitter) {
    return new $.MultiCropperOverlay(this, osdViewerId, windowId, state, eventEmitter);
  };

  // custom extend: annotationsLayer : todo tocheck
  $.MultiCropperLayer = function(options) {
    $jq.extend(true, this, {
      multiCropperList: null,
      viewer: null,
      drawTool: null,
      selected: null,
      hovered: null,
      windowId: null,
      mode: 'default',
      element: null,
      eventEmitter: null,
    }, options);

    this.init();
  };

  $.MultiCropperLayer.DISPLAY_multiCropper = 'displayMultiCropper';

  $.MultiCropperLayer.prototype = {
    init: function() {
      const _this = this;
      _this.eventEmitter.unsubscribe(('multiCropperModeChange.' + _this.windowId));
      _this.eventEmitter.unsubscribe(('slotLeave.' + _this.windowId));
      _this.eventEmitter.unsubscribe(('slotEnter.' + _this.windowId));

      this.createStateMachine();
      this.createRenderer();
      this.bindEvents();
      this.listenForActions();
    },

    listenForActions: function() {
      const _this = this;

      _this.eventEmitter.subscribe('multiCropperModeChange.' + _this.windowId, function(event, modeName) {
        _this.mode = modeName;
        _this.modeSwitch();
      });

      _this.eventEmitter.subscribe('multiCropperListLoaded.' + _this.windowId, function(event) {
        _this.multiCropperList = _this.state.getWindowMultiCropperList(_this.windowId);
        _this.updateRenderer();
      });

      _this.eventEmitter.subscribe('slotLeave.' + _this.windowId, function(event, eventData) {
        if (_this.layerState.current == "display") {
          _this.layerState.defaultState();
          _this.modeSwitch();
        }
      });

      _this.eventEmitter.subscribe('slotEnter.' + _this.windowId, function(event, eventData) {
        if (_this.element.showMultiCropper && _this.layerState.current == "display") {
          _this.layerState.defaultState();
          _this.modeSwitch();
        }
      });
    },

    bindEvents: function() {
      var _this = this;
    },

    createStateMachine: function() {
      var _this = this;
      this.layerState = StateMachine.create({
        events: [
          { name: 'startup', from: 'none', to: 'default' },
          { name: 'defaultState', from: ['default', 'display', 'create', 'edit'], to: 'default' },
          { name: 'displayMultiCropper', from: ['default', 'display', 'create', 'edit', 'newShape'], to: 'display' },
          { name: 'createMultiCropper', from: ['default','display'], to: 'create' },
          { name: 'createhape', from: 'edit', to: 'newShape'},
          { name: 'editMultiCropper', from: ['default','display', 'newShape'], to: 'edit' }
        ],
        callbacks: {
          onstartup: function(event) {
            _this.drawTool.enterDefault();
          },
          ondefaultState: function(event) {
            _this.drawTool.enterDefault();
          },
          ondisplayMultiCropper: function(event) {
            _this.drawTool.enterDisplayAnnotations();
          },
          oncreateMultiCropper: function(event) {
            _this.drawTool.enterCreateAnnotation();
          },
          oncreateShape: function(event) {
            _this.drawTool.enterCreateShape();
          },
          oneditMultiCropper: function(event) {
            _this.drawTool.enterEditAnnotations();
          }
        }
      });
    },

    createRenderer: function() {
      var _this = this;
      this.drawTool = new $.MultiCropperRegionDrawTool({
        osdViewer: _this.viewer,
        parent: _this,
        list: _this.multiCropperList, // must be passed by reference.
        visible: false,
        windowId: _this.windowId,
        state: _this.state,
        eventEmitter: _this.eventEmitter,
      });
      this.layerState.startup();
    },

    updateRenderer: function() {
      this.drawTool.list = this.multiCropperList;
      // this.modeSwitch();
    },

    modeSwitch: function() {
      if(this.mode === 'displayMultiCropper') {
      // this.layerState.displayMultiCropper();
      } else if(this.mode === 'editingMultiCropper') {
        this.layerState.editMultiCropper();
      } else if(this.mode === 'creatingMultiCropper') {
        if(this.layerState.current !== 'edit') {
          this.layerState.createMultiCropper();
        } else {
          this.layerState.createShape();
        }
      } else if(this.mode === 'default') {
        this.layerState.defaultState();
      } else {}
    }
  };

  // /* custom: annotations */
  // custom: rectangle
  $.MultiCropperRectangle = function(options) {
    $.Rectangle.call(this, options);
  };

  $.MultiCropperRectangle.prototype = Object.create($.Rectangle.prototype, {
  });

  // /* custom: utils */
  // custom: saveController
  const origSaveControllerPrototypeInit = $.SaveController.prototype.init;
  $.SaveController.prototype.init = function(config) {
    origSaveControllerPrototypeInit.apply(this, [config]);

    if(config && config.rectSelector) {
      jQuery.extend(_options, config.multiCropper);
    }
  };

  $.SaveController.prototype.getWindowMultiCropperList = function(windowId) {
    if(this.windowsMultiCropperLists) {
      return this.windowsMultiCropperLists[windowId];
    } else {
      return null;
    }
  };

  // custom extend: OsdRegionDrawTool
  $.MultiCropperRegionDrawTool = function(options) {
    $.OsdRegionDrawTool.call(this, options);
  };

  $.MultiCropperRegionDrawTool.prototype = Object.create($.OsdRegionDrawTool.prototype, {
    listenForActions: {
      value() {
        var _this = this;

        this._thisDestroy = function(){
          _this.destroy();
        };

        _this.osdViewer.addHandler('close', this._thisDestroy);

        this.eventsSubscriptions.push(this.eventEmitter.subscribe('DESTROY_EVENTS.'+this.windowId, function(event) {
          _this.destroy();
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('updateTooltips.' + _this.windowId, function(event, location, absoluteLocation) {
          if (_this.annoTooltip && !_this.annoTooltip.inEditOrCreateMode) {
            _this.showTooltipsFromousePosition(event, location, absoluteLocation);
          }
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('removeTooltips.' + _this.windowId, function() {
          jQuery(_this.osdViewer.element).qtip('destroy', true);
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('disableTooltips.' + _this.windowId, function() {
          if (_this.annoTooltip) {
            _this.annoTooltip.inEditOrCreateMode = true;
          }
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('enableTooltips.' + _this.windowId, function() {
          if (_this.annoTooltip) {
            _this.annoTooltip.inEditOrCreateMode = false;
          }
          _this.svgOverlay.restoreDraftShapes();
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('SET_ANNOTATION_EDITING.' + _this.windowId, function(event, options) {
          jQuery.each(_this.annotationsToShapesMap, function(key, paths) {
            // if we have a matching annotationId, pass the boolean value on for each path, otherwise, always pass false
            if (key === options.annotationId) {
              if (options.isEditable) {
                _this.eventEmitter.publish('SET_OVERLAY_TOOLTIP.' + _this.windowId, {"tooltip" : options.tooltip, "visible" : true, "paths" : paths});
              } else {
                _this.eventEmitter.publish('SET_OVERLAY_TOOLTIP.' + _this.windowId, {"tooltip" : null, "visible" : false, "paths" : []});
              }
              jQuery.each(paths, function(index, path) {
                path.data.editable = options.isEditable;
                if (options.isEditable) {
                  path.strokeWidth = (path.data.strokeWidth + 5) / _this.svgOverlay.paperScope.view.zoom;
                } else {
                  path.strokeWidth = path.data.strokeWidth / _this.svgOverlay.paperScope.view.zoom;
                }
                //just in case, force the shape to be non hovered
                var tool = _this.svgOverlay.getTool(path);
                tool.onHover(false, path, path.strokeWidth);
              });
            } else {
              jQuery.each(paths, function(index, path) {
                path.data.editable = false;
              });
            }
          });
          _this.svgOverlay.paperScope.view.draw();
        }));

        this.eventsSubscriptions.push(_this.eventEmitter.subscribe('refreshOverlay.' + _this.windowId, function (event) {
          _this.render();
        }));

        this.eventsSubscriptions.push(this.eventEmitter.subscribe("enableManipulation",function(event, tool){
          if(tool === 'mirror') {
            _this.horizontallyFlipped = true;
          }
        }));
        this.eventsSubscriptions.push(this.eventEmitter.subscribe("disableManipulation",function(event, tool){
          if(tool === 'mirror') {
            _this.horizontallyFlipped = false;
          }
        }));
      },
    },
  });
})(Mirador, jQuery);

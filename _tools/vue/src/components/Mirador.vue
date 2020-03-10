<template>
  <div id="mirador_canvas_container">
    <div id="mirador_canvas"></div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Mirador',

  props: {
    manifest: String,
    manifestData: Object,
    pos: Number,
    xywh: String,
  },

  mounted() {
    window.Mirador.DEFAULT_SETTINGS.windowSettings.sidePanel = false;
    window.Mirador.DEFAULT_SETTINGS.windowSettings.canvasControls
      .annotations.annotationState = 'on';

    if(this.manifest && this.manifestData) {
      const canvasID = this.manifestData
        .sequences[0].canvases[this.pos]['@id'];
      console.log(this.pos, canvasID)
      const windowOptions = (this.xywh.match(/^(\d+),(\d+),(\d+),(\d+)$/))
        ? {
          osdBounds: {
            x: parseInt(RegExp.$1),
            y: parseInt(RegExp.$2),
            width: parseInt(RegExp.$3),
            height: parseInt(RegExp.$4),
          },
          zoomLevel: 1,
        } : null;

      window.myMiradorInstance = window.Mirador({
        id: 'mirador_canvas',

        data: [{
          manifestUri: this.manifest,
          location: 'jikei',
        }],

        windowObjects: [{
          loadedManifest: this.manifest,
          canvasID: canvasID,
          windowOptions: windowOptions,
        }],
      });
    } else {
      window.myMiradorInstance = window.Mirador({
        id: 'mirador_canvas',
      });
    }
  },
};
</script>

<style>
#mirador_canvas {
  position: relative;
  left: 0;
  right: 0;
  margin: auto;
  max-width: 960px;
  min-height: 500px;
}
#mirador_canvas_container {
  display: block;
  left: 0;
  right: 0;
  margin: auto;
  /* width: 800px; */
  /* height: 500px; */
}
</style>

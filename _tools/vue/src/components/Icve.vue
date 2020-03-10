<template>
  <v-container>
    <v-row class="text-center">
      <v-col cols="12">
        <div id="image_canvas" style="height: 500px;"></div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'Icve',

    props: {
      manifest: String,
    },

    created() {
      console.log('icve');
      window.IIIFCurationViewerEmbedded({
        id: 'image_canvas',
        // data: {
        //   manifest: this.manifest,
        // }
      });
    },
  }
</script>

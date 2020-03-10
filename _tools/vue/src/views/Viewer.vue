<template>
  <div class="viewer">
    <v-layout wrap>
      <v-flex md12 lg8>
        <v-responsive aspect-ratio="16/9">
          <Component :is="viewer"
                      :manifest="manifest"
                      :manifestData="manifestData"
                      :pos="pos"
                      :xywh="xywh"
                      />
        </v-responsive>
      </v-flex>
      <v-flex md12 lg4>
        <Bibliography :data="manifestData"/>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import Icve from '@/components/Icve.vue';
import Mirador from '@/components/Mirador.vue';
import Bibliography from '@/components/Bibliography.vue';
import axios from 'axios';

export default {
  name: 'Viewer',
  components: {
    Icve,
    Mirador,
    Bibliography,
  },

  data() {
    return {
      viewer: '',
      manifest: '',
      manifestData: null,
      pos: 0,
      xywh: '',
    }
  },

  async created() {
    if(this.$route.query.manifest) {
      this.manifest = this.$route.query.manifest;
      if(this.$route.query.pos && this.$route.query.pos.match(/^\d+$/)) {
        this.pos = parseInt(this.$route.query.pos) - 1;
      }
      if(this.$route.query.xywh
        && this.$route.query.xywh.match(/^\d+,\d+,\d+,\d+$/)) {
        this.xywh = this.$route.query.xywh;
      }

      try {
        const res = await axios.get(this.manifest);
        this.manifestData = res.data;
      } catch(err) {
        console.log('Error:', err);
      }
    }

    if(['Icve', 'Mirador'].includes(this.$route.name)) {
      this.viewer = this.$route.name;
    } else {
      this.viewer = 'Mirador';
    }
  },
};
</script>

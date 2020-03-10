<template>
  <v-app>
    <v-app-bar
      app
      color="white"
      dense
      flat
    >
      <v-img
        alt="NIJL Logo"
        class="shrink mr-2"
        contain
        src="./assets/img/nijl_symbolmark.jpg"
        transition="scale-transition"
        width="30"
        ></v-img>

      <router-link
        to="/"
        class="shrink mr-2"
        min-width="100"
        >
        国文研字形検索β
      </router-link>

      <v-spacer></v-spacer>

      <v-text-field
        v-model="chars"
        class="mt-4"
        rounded
        solo
        dense
        clearable
        prepend-inner-icon="mdi-magnify"
        v-on:keypress.enter="onPrepareSearch"
        v-on:keyup.enter="onSearch"
        ></v-text-field>
    </v-app-bar>

    <v-content>
      <router-view/>
    </v-content>

    <v-footer
      app
      flat
      >
      <v-col
        class="text-center"
        cols="12"
        >
        2020&nbsp;&dash;&nbsp;
        <a target="_blank" href="https://www.nijl.ac.jp/">
          国文学研究資料館
        </a>&nbsp;
        <a target="_blank" href="https://www.nijl.ac.jp/cijproject/">
          『日本語の歴史的典籍の国際共同研究ネットワーク構築計画』
        </a>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  name: 'App',

  data() {
    return {
      chars: '',
      prepared: false,
    };
  },

  methods: {
    onPrepareSearch() {
      this.prepared = true;
    },
    onSearch() {
      if(this.prepared) {
        this.$router.push(`/search?q=${this.chars}`)
        this.prepared = false;
      }
    },
  },
};
</script>

<style scoped>
.v-app-bar a {
  color: #eb9396 !important;
  font-weight: bold !important;
  text-decoration: none !important;
}
.v-footer a {
  text-decoration: none !important;
}
</style>

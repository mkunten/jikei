<template>
  <div class="search">
    <v-pagination
      v-if="total > 0"
      v-model="page"
      :length="length"
      :total-visible="7"
      @input="onSearch"
      ></v-pagination>
    <v-lazy
      :options="{ threshold: .5 }"
      >
      <v-container
        fruid
        >
        total: {{ this.total }} ({{ this.page }}/{{ this.length }})
        <transition-group
          tag="div"
          class="layout row wrap"
          >
          <v-flex
            :key="item.id"
            v-for="item in data.list"
            >
            <v-card
              class="ma-4 pa-4"
              width="200px"
              height="200px"
              raised
              >
              <v-layout
                style="height:100px;"
                >
                <router-link
                  class="ma-auto"
                  :to="`/viewer?manifest=${apiBaseURL}/biblio/${item.source.bid}/${encodeURIComponent(item.character)}/manifest&pos=${item.source.pos}&xywh=${item.source.x},${item.source.y},${item.source.width},${item.source.height}`"
                  >
                  <v-img
                    :src="item.thumbnail_url"
                    max-width="100px"
                    max-height="100px"
                    aspect-racio="1"
                    ></v-img>
                </router-link>
              </v-layout>
              <v-card-title
                class="justify-center pa-0"
                >
                {{ item.character }}
              </v-card-title>
              <v-card-text
                class="text-center"
                >
                『{{ item.source.title }}』<br />
                <a target="_blank" title="新日本古典籍総合データベースで開く" :href="`https://kotenseki.nijl.ac.jp/biblio/${item.source.bid}/viewer/${item.source.frame}`">
                  {{ item.source.bid }} {{ item.source.frame }}{{ item.source.side }} <v-icon>open_in_new</v-icon>
                </a>
              </v-card-text>
            </v-card>
          </v-flex>
        </transition-group>
      </v-container>
    </v-lazy>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Search',

  data() {
    return {
      // apiBaseURL: 'http://localhost:58080/jikei/api',
      apiBaseURL: '/jikei/api',
      chars: '',
      page: 1,
      limit: 50,
      length: 0,
      total: 0,
      data: {},
    };
  },

  watch: {
    '$route'(to, from) {
      if(to.fullPath !== from.fullPath) {
        this.page = 1;
        this.onSearch();
      }
    }
  },

  created() {
    this.onSearch();
  },

  methods: {
    onSearch() {
      this.chars = this.$route.query['q'];
      if(this.chars) {
        axios.get(`${this.apiBaseURL}/search?q=${this.chars}&offset=${(this.page - 1) * this.limit}&limit=${this.limit}`)
          .then((res) => {
            console.log(res);
            this.total = res.data.total;
            this.length = Math.ceil(this.total / this.limit);
            this.data = res.data;
          })
          .catch((err) => {
            console.log(err);
          });
      }
    },
  },
};
</script>

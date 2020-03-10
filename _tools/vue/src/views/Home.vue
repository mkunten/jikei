<template>
  <div class="home">
    <v-row align="center">
      <!-- <v&#45;row justify="space&#45;around"> -->
      <!--   <v&#45;form ref="form" @submit.prevent="onSearch"> -->
      <!--     <v&#45;text&#45;field -->
      <!--       v&#45;model="char" -->
      <!--       label="Character" -->
      <!--       required -->
      <!--       ></v&#45;text&#45;field> -->
      <!--   </v&#45;form> -->
      <!--   <v&#45;btn -->
      <!--     color="success" -->
      <!--     @click="onSearch" -->
      <!--     >Search</v&#45;btn> -->
      <!-- </v&#45;row> -->
    </v-row>
  </div>
</template>

<script>
export default {
  name: 'Home',
  components: {
  },

  data() {
    return {
      char: '',
    };
  },

  methods: {
    onSearch() {
      this.$router.push(`/search?q=${this.char}`)
    },
  },
};
</script>

<template>
  <div class="home">
    <v-container>
      <v-row justify="center">
        <v-col sm="10" md="8" lg="6">
          <p>
          国文学研究資料館では、ROIS-DS人文学オープンデータ共同研究センター（CODH）の協力のもと、古典籍に現れる多様な字形のデータセット「日本古典籍くずし字データセット」（CC BY-SA、『日本古典籍くずし字データセット』（国文研ほか所蔵／CODH加工） doi:10.20676/00000340）の公開を進めています。
          </p>

          <p>
          本サイトでは、すでに公開されているこのデータセットのうち、当館が「新日本古典籍総合データベース」で公開している書誌に登場する字形（2020年3月現在、41書誌3,986字種847,918字）を対象に、IIIFでの活用をはじめ、ウェブ上での多様な利用に供するためのAPIの構築などを実験的に行っています。また、API利用の一例として、字形検索を提供しています。右上検索窓よりお試しください。
          </p>

          <p>
          Links:

          <ul>
            <li>
              <a target="_blank" href="https://kotenseki.nijl.ac.jp/">
                新日本古典籍総合データベース
              </a>
            </li>
            <li>
              <a target="_blank" href="http://codh.rois.ac.jp/char-shape/">
                日本古典籍くずし字データセット（ROIS-DS CODH）
              </a>
            </li>
          </ul>
          </p>
        </v-col>
      </v-row>
    </v-container>
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

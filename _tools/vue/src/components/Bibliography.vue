<template>
  <div id="bibliography">
    <div class="link">
      <a v-if="dbLink" :href="dbLink" :alt="messageOpen" target="_blank">
        {{ messageOpen }}
      </a>
    </div>
    <v-simple-table>
      <template v-slot:default>
        <tbody>
          <template v-if="data">
            <tr>
              <th></th>
              <td></td>
            </tr>
            <tr :key="item.label" v-for="item in data.metadata">
              <th>{{ item.label }}</th>
              <td>{{ item.value }}</td>
            </tr>
            <tr>
              <th>LICENSE</th>
              <td>
                <a :href="data.license" :alt="data.license" target="_blank">
                  {{ data.license }}
                </a>
              </td>
            </tr>
            <tr>
              <th>ATTRIBUTION</th>
              <td>{{ data.attribution }}</td>
            </tr>
          </template>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>
  export default {
    name: 'Bibliography',

    data: () => {
      return {
        messageOpen: "新日本古典籍総合DBで開く",
      };
    },

    props: {
      data: Object,
    },

    computed: {
      dbLink() {
        let dbLink = '';
        const metadata = {};
        if(this.data && this.data.metadata) {
          this.data.metadata.forEach((item) => {
            metadata[item.label] = item.value;
          });
          if(metadata['BID'] && metadata['BID'].match(/^\d{9}$/)) {
            if(metadata['FRAME'] && metadata['FRAME'].match(/^\d+$/)) {
              dbLink = `https://kotenseki.nijl.ac.jp/biblio/${metadata['BID']}/viewer/${metadata['FRAME']}`;
            }
            else {
              dbLink = `https://kotenseki.nijl.ac.jp/biblio/${metadata['BID']}/viewer/`;
            }
          }
        }
        return dbLink;
      },
    },
  }
</script>

<style scoped>
#bibliography {
  max-width: 800px;
  left: 1rem;
  right: 1rem;
  margin: 1rem auto;
}

.link {
  margin-left: 1rem;
}
</style>

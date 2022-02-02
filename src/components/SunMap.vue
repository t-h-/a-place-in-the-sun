<template>
  <div class="sunmap">
    <h1>{{ msg }}</h1>
      <l-map style="height: 600px" :zoom="zoom" :center="center">
      <l-tile-layer :url="url" :attribution="attribution"></l-tile-layer>
      <LHeatmap :latLngs=latLngArray :max="maxValue" :radius="15"></LHeatmap>
      </l-map>
  </div>
</template>

<script>
import { LMap, LTileLayer } from 'vue2-leaflet';
import LHeatmap from './Vue2LeafletHeatmap.vue';
// import LHeatmap from './Vue2LeafletHeatmapForked.vue';
// import LHotline from 'vue2-leaflet-hotline';

export default {
  name: 'SunMap',
  props: {
    msg: String,
  },
  components: {
    LMap,
    LTileLayer,
    LHeatmap,
  },
  data() {
    return {
      url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
      attribution:
        '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a> contributors',
      zoom: 15,
      center: [51.50, 30.5],
      latLngArray: [[50.5, 30.5, 40], [50.51, 30.51, 80]],
      maxValue: 100,
    };
  },
};
</script>

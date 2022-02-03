<template>
  <div class="sunmap">
    <h1>{{ msg }}</h1>
      <l-map ref="leafletMap" @moveend="recalc" style="height: 600px" :zoom="zoom" :center="center">
      <l-tile-layer :url="url" :attribution="attribution"></l-tile-layer>
      <LHeatmap ref="heatmapLayer" :latLng=latLngArray :radius=radius :blur=blur :gradient=gradient :max=max></LHeatmap>
      </l-map>
  </div>
</template>

<script>
import { LMap, LTileLayer } from 'vue2-leaflet';
import LHeatmap from './Vue2LeafletHeatmap.vue';

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
      map: null,
      url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
      attribution:
        '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a> contributors',
      zoom: 11,
      center: [50.50, 30.5],
      latLngArray: [
        [50.50, 30.50, 40],
        [50.50, 30.51, 80.80],
        [50.50, 30.52, 80.80],
        [50.50, 30.53, 80.80],
        [50.50, 30.64, 50],
        [50.50, 30.74, 80.80],
      ],
      max: null,
      radius: 100,
      blur: 30,
      gradient: { 0.1: 'green', 0.5: 'yellow', 1.0: 'orange' },
    };
  },
  methods: {
    recalc(mapObject) {
      console.log('recalc');
      console.log(mapObject);
      const bnds = this.$refs.leafletMap.mapObject.getBounds();
      console.log(bnds);
      this.latLngArray.push([50.45, 30.74, 80.80]);
      this.$refs.heatmapLayer.addLatLng([50.55, 30.74, 80.80]);
    },
  },
};
</script>

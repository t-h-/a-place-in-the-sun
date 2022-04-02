<template>
  <div class="sunmap" ref="rootC">
    <h1>{{ msg }}</h1>
      <l-map ref="leafletMap" @moveend="onMoveEnd" style="height: 600px" :zoom="zoom" :center="center">
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
      center: [46.2353601361914, 6.200332068813093],
      latLngArray: [],
      sunIntensities: {},
      max: null,
      radius: 100,
      blur: 30,
      gradient: { 0.1: 'green', 0.5: 'yellow', 1.0: 'orange' },
    };
  },
  async mounted() {
    await this.$nextTick();
    this.onMoveEnd();
  },
  methods: {
    async onMoveEnd(mapObject) {
      const bnds = this.$refs.leafletMap.mapObject.getBounds();
      // const bnds = [lat, lng]; // this.$refs.leafletMap.mapObject.getBounds();
      const southWest = bnds.getSouthWest();
      const northEast = bnds.getNorthEast();
      const pixelsY = this.$refs.rootC.clientHeight;
      const pixelsX = this.$refs.rootC.clientWidth;
      await this.$store.dispatch('queryAllPointsInBounds', {
        southWest, northEast, pixelsX, pixelsY,
      });
      this.$store.dispatch('sampleLatLngArray2', {
        southWest, northEast, pixelsX, pixelsY,
      });
      const lla = this.$store.getters.getLatLngArray;
      this.$refs.heatmapLayer.setLatLngs(lla);
    },
  },
};
</script>

<style>
.leaflet-heatmap-layer {
  opacity: .7;
}
</style>

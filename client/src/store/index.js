import Vue from 'vue';
import Vuex from 'vuex';

const _ = require('lodash');
const axios = require('axios');

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    latLngArray: [],
  },
  mutations: {
    setLatLngArray(state, latLngArray) {
      state.latLngArray = latLngArray;
    },
  },
  actions: {
    async queryAllPointsInBounds(context, {
      northEast, southWest, numPointsX, numPointsY,
    }) {
      const swLat = southWest.lat;
      const swLng = southWest.lng;
      const neLat = northEast.lat;
      const neLng = northEast.lng;

      const cfg = {
        method: 'post',
        url: 'http://localhost:8083/sunnyness/grid',
        headers: { },
        data: {
          box: {
            top_left_lat: neLat,
            top_left_lng: neLng,
            bottom_right_lat: swLat,
            bottom_right_lng: swLng,
          },
          num_points: {
            lat: numPointsX.numPointsX,
            lng: numPointsY.numPointsY,
          },
        },
      };

      // console.log(`querying ${config}`);
      const result = await axios(cfg);

      const latLngArray = [];
      result.data.grid.values.forEach((res) => {
        latLngArray.push([res.lat, res.lng, res.val]);
      });

      context.commit('setLatLngArray', latLngArray);
    },
  },
  modules: {
  },
  getters: {
    getLatLngArray: (state) => state.latLngArray,
  },
});

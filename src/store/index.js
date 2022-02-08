import Vue from 'vue';
import Vuex from 'vuex';

const _ = require('lodash');
const axios = require('axios');

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    sunnynessCache: {},
    latLngArray: [],
  },
  mutations: {
    setSunnyness(state, { lat, lng, sunnyness }) {
      if (!(lat in state.sunnynessCache)) {
        state.sunnynessCache[lat] = {};
      }
      state.sunnynessCache[lat][lng] = sunnyness;
    },
    setLatLngArray(state, latLngArray) {
      console.log('setLatLngArray');
      state.latLngArray = latLngArray;
    },
  },
  actions: {
    async queryAllPointsInBounds(context, {
      southWest, northEast, pixelsX, pixelsY,
    }) {
      console.log('state: queryAllPointsInBounds');
      const queriesPerPixel = 0.1;

      const swLat = southWest.lat;
      const swLng = southWest.lng;
      const neLat = northEast.lat;
      const neLng = northEast.lng;

      const stepLat = Math.max((northEast.lat - southWest.lat) / (pixelsX * queriesPerPixel), 0.1);
      const stepLng = Math.max((northEast.lng - southWest.lng) / (pixelsY * queriesPerPixel), 0.1);

      const allPromises = [];
      _.range(swLat - stepLat, neLat + stepLat, stepLat).forEach((cLat) => {
        _.range(swLng - stepLng, neLng + stepLng, stepLng).forEach((cLng) => {
          const lat = Math.round(cLat * 10) / 10;
          const lng = Math.round(cLng * 10) / 10;
          const pm = this.dispatch('queryCoord', { lat, lng });
          allPromises.push(pm);
        });
      });
      const results = await Promise.all(allPromises);
      results.forEach((res) => {
        const data = res.data;
        const lat = data.location.lat;
        const lng = data.location.lon;
        const sunnyness = 100 - data.current.cloud;
        context.commit('setSunnyness', { lat, lng, sunnyness });
      });
    },
    async queryCoord(context, { lat, lng }) {
      // TODO use real cache
      if (lat in context.state.sunnynessCache && lng in context.state.sunnynessCache[lat]) {
        const sunnyness = context.state.sunnynessCache[lat][lng];
        console.log(`queryCoord: cache hit: ${lat} ${lng} ${sunnyness}`);
        return Promise.resolve({
          data: {
            current: {
              cloud: 100 - sunnyness,
            },
            location: {
              lat: lat,
              lon: lng,
            },
          },
        });
      }

      const config = {
        method: 'get',
        url: `http://api.weatherapi.com/v1/current.json?key=591b7934afcf484fa3191051223101&q=${lat},${lng}&aqi=no`,
        headers: { },
      };

      console.log(`queryCoord: querying ${lat} ${lng}`);
      return axios(config);
    },
    sampleLatLngArray(context, { southWest, northEast }) {
      const latLngArray = [];
      Object.keys(context.state.sunnynessCache)
        .filter((lat) => lat > southWest.lat && lat < northEast.lat)
        .forEach((lat) => {
          const lngObj = context.state.sunnynessCache[lat];
          Object.keys(lngObj)
            .filter((lng) => lng > southWest.lng && lng < northEast.lng)
            .forEach((lng) => {
              const sn = lngObj[lng];
              latLngArray.push([lat, lng, sn]);
            });
        });
      context.commit('setLatLngArray', latLngArray);
    },
  },
  modules: {
  },
  getters: {
    getLatLngArray: (state) => state.latLngArray,
    getSunnynessCache: (state) => state.sunnynessCache,
  },
});

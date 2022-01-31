<template>
  <div class="sun">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld msg="Lets see where the Sun is!"/>
  </div>
</template>

<script>
// @ is an alias to /src
import SunMap from '@/components/SunMap.vue';

export default {
  name: 'Sun',
  components: {
    SunMap,
  },
};
</script>

<template>
  <div class="sun">
    <SunMap msg="Let's see where the Sun is!"/>
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

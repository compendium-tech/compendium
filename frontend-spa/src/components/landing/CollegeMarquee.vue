<template>
  <div class="relative overflow-hidden">
    <Vue3Marquee :clone='true'>
      <div class="flex items-center space-x-16 mr-16">
        <div v-for="(logo, index) in logos" :key="`first-${index}`">
          <img :src="logo.src" :alt="logo.alt"
            :class="['object-contain opacity-80 hover:opacity-100 transition-opacity', logo.heightClass || 'h-12']" />
        </div>
      </div>
    </Vue3Marquee>
  </div>
</template>

<script setup lang="ts">
import { Vue3Marquee } from 'vue3-marquee'

import { ref } from 'vue';

const logos = ref([
  {
    src: '/logos/colleges/harvard.png',
    alt: 'Harvard University',
    heightClass: 'h-16'
  },
  {
    src: '/logos/colleges/stanford.png',
    alt: 'Stanford University',
    heightClass: 'h-30'
  },
  {
    src: '/logos/colleges/ucberkley.png',
    alt: 'UC Berkeley',
    heightClass: 'h-32'
  },
  {
    src: '/logos/colleges/mit.png',
    alt: 'MIT',
    heightClass: 'h-18'
  },
  {
    src: '/logos/colleges/princeton.png',
    alt: 'Princeton',
    heightClass: 'h-32'
  },
  {
    src: '/logos/colleges/oxford.png',
    alt: 'Oxford',
    heightClass: 'h-18'
  },
  {
    src: '/logos/colleges/caltech.png',
    alt: 'Caltech',
    heightClass: 'h-18'
  },
]);
</script>

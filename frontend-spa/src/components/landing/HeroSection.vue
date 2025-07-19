<template>
  <section
    class="max-w-full pt-48 mb-32 flex flex-col justify-center items-center text-center px-4 bg-gradient-to-b from-primary-50 to-white">
    <div class="animate-fade-in-up">
      <div class=" bg-primary-100 text-primary-600 px-6 py-2 rounded-full inline-flex items-center mb-6">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <span>Introducing Compendium</span>
      </div>

      <h1 class="text-3xl md:text-4xl lg:text-6xl font-bold mb-6">
        Your <span class="text-primary-600">College Application</span> Assistant
      </h1>

      <p class="text-lg md:text-xl mb-10 text-gray-600 px-8">
        From university research to essay polishing - we streamline your journey to dream schools in US, Europe, and
        China.
      </p>

      <RouterLink to="/auth/signup">
        <button
          class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-8 rounded-xl text-lg transition-all transform hover:scale-105 shadow-lg hover:shadow-xl">
          Start your academic journey
        </button>
      </RouterLink>

      <CollegeMarquee class="mt-16 max-w-screen" />
    </div>
  </section>
</template>

<script setup lang="ts">
import CollegeMarquee from './CollegeMarquee.vue'
import { RouterLink } from 'vue-router'
</script>

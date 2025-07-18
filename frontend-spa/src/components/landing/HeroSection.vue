<template>
  <section
    class="pt-48 mb-32 flex flex-col justify-center items-center text-center px-4 bg-gradient-to-b from-primary-50 to-white">
    <div class="opacity-0" :class="{ 'animate-fade-in-up': animate }">
      <div class="bg-primary-100 text-primary-600 px-6 py-2 rounded-full inline-flex items-center mb-6">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <span>Introducing Compendium</span>
      </div>

      <h1 class="text-2xl md:text-4xl lg:text-6xl font-bold mb-6">
        Your <span class="text-primary-600">College Application</span> Assistant
      </h1>

      <p class="text-xl max-w-3xl mx-auto mb-10 text-gray-600">
        From university research to essay polishing - we streamline your journey to dream schools in US, Europe, and
        China.
      </p>

      <button
        class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-8 rounded-xl text-lg transition-all transform hover:scale-105 shadow-lg hover:shadow-xl">
        Start Free Trial
      </button>

      <CollegeMarquee class="mt-16" />
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CollegeMarquee from './CollegeMarquee.vue'

const animate = ref(false)

onMounted(() => {
  setTimeout(() => {
    animate.value = true
  }, 100)
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.8s ease-out forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

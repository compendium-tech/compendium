<template>
  <section
    class="min-h-screen flex flex-col justify-center items-center text-center py-20 px-4 bg-gradient-to-b from-primary-50 to-white">
    <div class="opacity-0" :class="{ 'animate-fade-in-up': animate }">
      <div class="bg-primary-100 text-primary-600 px-6 py-2 rounded-full inline-flex items-center mb-6">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <span>Introducing Compendium</span>
      </div>

      <h1 class="text-4xl md:text-6xl lg:text-7xl font-bold mb-6">
        Your <span class="text-primary-600">University Application</span> Assistant
      </h1>

      <p class="text-xl max-w-3xl mx-auto mb-10 text-gray-600">
        From university research to essay polishing - we streamline your journey to dream schools in US, Europe, and
        China.
      </p>

      <button
        class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-8 rounded-full text-lg transition-all transform hover:scale-105 shadow-lg hover:shadow-xl">
        Start Free Trial
      </button>

      <div class="mt-16">
        <AnimatedCard :delay="0.8">
          <div class="bg-white p-8 rounded-2xl shadow-xl border border-primary-100 max-w-4xl mx-auto">
            <div class="flex flex-col md:flex-row items-center justify-center">
              <div
                class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center text-primary-600 mb-4 md:mb-0 md:mr-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <p class="text-left text-lg">
                "Compendium helped me secure offers from 3 Ivy League schools with their essay feedback system!"
              </p>
            </div>
          </div>
        </AnimatedCard>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AnimatedCard from './AnimatedCard.vue'

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

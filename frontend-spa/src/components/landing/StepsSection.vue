<template>
  <section class="py-20 bg-gray-50 relative">
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden opacity-10">
      <div class="absolute border-2 border-primary-400 rounded-full"
        style="width: 180px; height: 180px; top: 15%; left: 5%; animation: pulse 12s infinite">
      </div>
      <div class="absolute border-2 border-primary-400 rounded-full"
        style="width: 220px; height: 220px; top: 25%; left: 75%; animation: pulse 15s infinite">
      </div>
      <div class="absolute border-2 border-primary-400 rounded-full"
        style="width: 150px; height: 150px; top: 65%; left: 10%; animation: pulse 18s infinite">
      </div>
      <div class="absolute border-2 border-primary-400 rounded-full"
        style="width: 250px; height: 250px; top: 45%; left: 60%; animation: pulse 14s infinite">
      </div>
      <div class="absolute border-2 border-primary-400 rounded-full"
        style="width: 120px; height: 120px; top: 75%; left: 80%; animation: pulse 16s infinite">
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 relative z-10">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-4">
        Your Path to <span class="text-primary-600">Admission Success</span>
      </h2>
      <p class="text-xl text-center mb-16 text-gray-600 max-w-3xl mx-auto">
        How Compendium guides you through every step of the application journey
      </p>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 px-8">
        <div v-for="(step, _index) in steps" :key="step.id" :delay="step.delay">
          <div
            class="bg-white p-8 rounded-2xl h-full flex flex-col transition-all transform duration-300 hover:scale-105 hover:border-primary-400 hover:border-solid border-hidden border-2 animate-fade-in-up">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">{{ step.number }}</div>
            </div>
            <h3 class="text-xl font-bold mb-3">{{ step.title }}</h3>
            <p class="text-gray-600 text-left">{{ step.description }}</p>
          </div>
        </div>
      </div>

      <div class="mt-16 text-center">
        <button
          class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-8 rounded-xl text-lg transition-all transform hover:scale-105">
          See How It Works
        </button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue';

const steps = ref([
  {
    id: 1,
    number: 1,
    title: 'University Research',
    description: 'Explore 5,000+ global institutions with detailed profiles to find your perfect academic match.',
    delay: 0.1,
  },
  {
    id: 2,
    number: 2,
    title: 'Application Crafting',
    description: 'Get guidance for standout essays, recommendations, and activity descriptions.',
    delay: 0.3,
  },
  {
    id: 3,
    number: 3,
    title: 'Exams Preparation',
    description: 'Master standardized tests with personalized study plans and proven strategies.',
    delay: 0.5,
  },
]);
</script>

<style>
@keyframes pulse {

  0%,
  100% {
    transform: scale(1);
    opacity: 0.1;
  }

  50% {
    transform: scale(1.05);
    opacity: 0.15;
  }
}
</style>

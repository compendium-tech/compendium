<template>
  <section class="py-20 bg-gray-50 relative">
    <!-- Decorative elements -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden opacity-10">
      <div v-for="i in 10" :key="i"
           class="absolute border-2 border-primary-400 rounded-full"
           :style="{
             width: `${Math.random() * 200 + 100}px`,
             height: `${Math.random() * 200 + 100}px`,
             top: `${Math.random() * 100}%`,
             left: `${Math.random() * 100}%`,
             animation: `pulse ${Math.random() * 10 + 10}s infinite`
           }">
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 relative z-10">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-4">
        Your Path to <span class="text-primary-600">Admission Success</span>
      </h2>
      <p class="text-xl text-center mb-16 text-gray-600 max-w-3xl mx-auto">
        How Compendium guides you through every step of the application journey
      </p>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8">
        <!-- Step 1 -->
        <AnimatedCard :delay="0.1">
          <div class="bg-white p-8 rounded-2xl shadow-lg h-full flex flex-col items-center text-center transition-all transform hover:scale-105 hover:shadow-xl">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">1</div>
            </div>
            <h3 class="text-xl font-bold mb-3">University Research</h3>
            <p class="text-gray-600">Explore our database of 5,000+ institutions with detailed profiles and admission stats</p>
          </div>
        </AnimatedCard>

        <!-- Step 2 -->
        <AnimatedCard :delay="0.2">
          <div class="bg-white p-8 rounded-2xl shadow-lg h-full flex flex-col items-center text-center transition-all transform hover:scale-105 hover:shadow-xl">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">2</div>
            </div>
            <h3 class="text-xl font-bold mb-3">Smart Shortlisting</h3>
            <p class="text-gray-600">AI-powered recommendations based on your profile, preferences, and goals</p>
          </div>
        </AnimatedCard>

        <!-- Step 3 -->
        <AnimatedCard :delay="0.3">
          <div class="bg-white p-8 rounded-2xl shadow-lg h-full flex flex-col items-center text-center transition-all transform hover:scale-105 hover:shadow-xl">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">3</div>
            </div>
            <h3 class="text-xl font-bold mb-3">Essay Crafting</h3>
            <p class="text-gray-600">AI-assisted writing, editing, and feedback for compelling personal statements</p>
          </div>
        </AnimatedCard>

        <!-- Step 4 -->
        <AnimatedCard :delay="0.4">
          <div class="bg-white p-8 rounded-2xl shadow-lg h-full flex flex-col items-center text-center transition-all transform hover:scale-105 hover:shadow-xl">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">4</div>
            </div>
            <h3 class="text-xl font-bold mb-3">Recommendations</h3>
            <p class="text-gray-600">Tools to help you secure strong recommendation letters from teachers</p>
          </div>
        </AnimatedCard>

        <!-- Step 5 -->
        <AnimatedCard :delay="0.5">
          <div class="bg-white p-8 rounded-2xl shadow-lg h-full flex flex-col items-center text-center transition-all transform hover:scale-105 hover:shadow-xl">
            <div class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center mb-6">
              <div class="text-primary-600 text-2xl font-bold">5</div>
            </div>
            <h3 class="text-xl font-bold mb-3">Exam Prep</h3>
            <p class="text-gray-600">Personalized study plans for SAT, ACT, TOEFL, IELTS, and other exams</p>
          </div>
        </AnimatedCard>
      </div>

      <div class="mt-16 text-center">
        <button class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-8 rounded-full text-lg transition-all transform hover:scale-105">
          See How It Works
        </button>
      </div>
    </div>
  </section>
</template>

<script setup>
import AnimatedCard from './AnimatedCard.vue'
</script>

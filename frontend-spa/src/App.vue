<template>
  <div class="font-sans min-h-screen">
    <header class="fixed w-full bg-white/90 backdrop-blur-sm py-4 px-6 z-50">
      <div class="max-w-6xl mx-auto flex justify-between items-center">
        <div class="flex items-center">
          <div class="w-13 h-13 rounded-full flex items-center justify-center">
            <img src="/public/icon.svg" />
          </div>
          <span class="ml-3 text-xl font-bold">Compendium</span>
          <span class="text-xl font-bold text-primary-600">.</span>
        </div>
        <button class="bg-primary-600 hover:bg-primary-700 text-white py-2 px-6 rounded-full transition-all">
          Try Demo
        </button>
      </div>
    </header>

    <main>
      <HeroSection />
      <ProblemStats />
      <PricingCards />
      <Footer />

      <section class="py-20 bg-gradient-to-r from-primary-400 to-primary-600">
        <div class="max-w-4xl mx-auto text-center px-4">
          <h2 class="text-4xl font-bold text-white mb-6">
            Ready to transform your application journey?
          </h2>
          <p class="text-xl text-primary-100 mb-10 max-w-2xl mx-auto">
            Join thousands of students who've secured spots in their dream schools
          </p>
          <div class="flex flex-col sm:flex-row justify-center gap-4">
            <button
              class="bg-white text-primary-600 font-bold py-3 px-8 rounded-full text-lg hover:bg-gray-100 transition-all">
              Start Free Trial
            </button>
            <button
              class="bg-transparent border-2 border-white text-white py-3 px-8 rounded-full text-lg hover:bg-white/10 transition-all">
              Book Demo
            </button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import Footer from './components/Footer.vue'
import HeroSection from './components/HeroSection.vue'
import ProblemStats from './components/ProblemStats.vue'
import PricingCards from './components/PricingCards.vue'
</script>

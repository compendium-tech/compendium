<template>
  <main>
    <Header />
    <HeroSection />
    <ProblemStatsSection />
    <StepsSection />
    <PricingCardsSection />
    <ReadyToStartSection />

    <Footer />
  </main>
</template>

<script setup>
import Header from '../layout/Header.vue'
import Footer from '../layout/Footer.vue'
import StepsSection from "./StepsSection.vue"
import HeroSection from './HeroSection.vue'
import ProblemStatsSection from './ProblemStatsSection.vue'
import PricingCardsSection from './PricingCardsSection.vue'
import ReadyToStartSection from "./ReadyToStartSection.vue"
</script>

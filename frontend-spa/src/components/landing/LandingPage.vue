<template>
  <StandardLayout>
    <div>
      <HeroSection />
      <ProblemStatsSection />
      <StepsSection />
      <PricingCardsSection />
      <ReadyToStartSection />
    </div>
  </StandardLayout>
</template>

<script setup>
import StandardLayout from "../layout/StandardLayout.vue"
import StepsSection from "./StepsSection.vue"
import HeroSection from './HeroSection.vue'
import ProblemStatsSection from './ProblemStatsSection.vue'
import PricingCardsSection from './PricingCardsSection.vue'
import ReadyToStartSection from "./ReadyToStartSection.vue"
</script>

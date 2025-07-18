<template>
  <section class="py-12 bg-white relative overflow-hidden">
    <div class="max-w-7xl mx-auto px-4 mb-12">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-6">
        Trusted by Students Applying to
      </h2>
    </div>

    <!-- Marquee container -->
    <div class="relative h-24 overflow-hidden">
      <!-- First marquee -->
      <div class="absolute top-0 left-0 h-full flex items-center animate-marquee whitespace-nowrap">
        <div v-for="(logo, index) in logos" :key="index" class="mx-8 flex items-center">
          <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16" />
          <span class="ml-3 text-xl font-semibold">{{ logo }}</span>
        </div>
      </div>

      <!-- Second marquee (duplicate for seamless looping) -->
      <div class="absolute top-0 left-full h-full flex items-center animate-marquee whitespace-nowrap">
        <div v-for="(logo, index) in logos" :key="index" class="mx-8 flex items-center">
          <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16" />
          <span class="ml-3 text-xl font-semibold">{{ logo }}</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const logos = [
  "Harvard University",
  "Stanford University",
  "MIT",
  "University of Oxford",
  "ETH Zurich",
  "Sorbonne University",
  "Tsinghua University",
  "University of Tokyo",
  "National University of Singapore",
  "University of Melbourne"
]
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(0%);
  }

  100% {
    transform: translateX(-100%);
  }
}

.animate-marquee {
  animation: marquee 20s linear infinite;
}
</style>

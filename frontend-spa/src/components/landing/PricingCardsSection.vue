<template>
  <section class="py-20 px-4 bg-primary-50">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-3xl md:text-4xl font-bold mb-4 text-center">
        Simple, <span class="text-primary-600">Affordable</span> Pricing
      </h2>
      <p class="text-xl text-center mb-16 text-gray-600 max-w-2xl mx-auto">
        Pay only for what you need. No hidden fees, no surprises.
      </p>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div v-for="(plan, index) in pricing" :key="index">
          <div :class="`bg-white rounded-2xl overflow-hidden shadow-lg border h-full flex flex-col ${plan.highlight
            ? 'border-primary-600 ring-2 ring-primary-600 md:scale-105 z-10'
            : 'border-gray-200'
            }`" class="animate-fade-in-up">
            <div v-if="plan.highlight" class="bg-primary-600 text-white text-center py-2">
              Most Popular
            </div>
            <div class="p-8 flex-grow flex flex-col">
              <h3 class="text-2xl font-bold mb-2">{{ plan.name }}</h3>
              <div class="flex items-baseline mb-4">
                <span class="text-5xl font-bold">{{ plan.price }}</span>
                <span class="text-gray-500">/month</span>
              </div>
              <p class="text-gray-600 mb-6">{{ plan.description }}</p>

              <ul class="space-y-4 mb-8 flex-grow">
                <li v-for="(feature, i) in plan.features" :key="i" class="flex items-start">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500 mr-2 flex-shrink-0 mt-0.5"
                    fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>{{ feature }}</span>
                </li>
              </ul>

              <button :class="`w-full py-3 rounded-xl font-medium mt-auto ${plan.highlight
                ? 'bg-primary-600 hover:bg-primary-700 text-white'
                : 'bg-gray-100 hover:bg-gray-200 text-gray-800'
                }`">
                Get Started
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
interface PricingCard {
  name: string
  price: string
  description: string
  features: string[]
  highlight: boolean
}

const pricing: PricingCard[] = [
  {
    name: "Student",
    price: "$5",
    description: "Perfect for individual students",
    features: [
      "University database access",
      "Essay and extracurricular activity reviews",
      "Exam preparation resources"
    ],
    highlight: false
  },
  {
    name: "Team",
    price: "$10",
    description: "For small groups & counselors",
    features: [
      "Everything in Starter",
      "Invite 15 students",
      "Invite counselors and recommenders to your personalized workspace"
    ],
    highlight: true
  },
  {
    name: "Community",
    price: "$30",
    description: "Schools & large organizations",
    features: [
      "Everything in Pro",
      "Invite 150+ students",
      "Advanced analytics"
    ],
    highlight: false
  }
]
</script>

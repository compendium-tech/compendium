<template>
  <section class="mb-12 px-4 bg-white">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-3xl md:text-4xl font-bold mb-16 text-center">
        The Application <span class="text-primary-600">Struggle</span> Is Real
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <AnimatedCard v-for="(stat, index) in stats" :key="index" :delay="0.1 * index">
          <div class=" p-6 rounded-xl text-center h-full">
            <div class="text-primary-600 mb-4 flex justify-center">
              <component :is="stat.icon" class="h-10 w-10" />
            </div>
            <h3 class="text-3xl font-bold mb-2">
              <span ref="counterEls">{{ stat.initialValue }}</span>{{ stat.unit }}
            </h3>
            <p class="text-gray-600">{{ stat.label }}</p>
          </div>
        </AnimatedCard>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AnimatedCard from './AnimatedCard.vue'

// Icons
const AlertTriangle = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  `
}

const Clock = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  `
}

const Calendar = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
    </svg>
  `
}

const Percent = {
  template: `
    <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 14l6-6m-5.5.5h.01m4.99 5h.01M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16l3.5-2 3.5 2 3.5-2 3.5 2zM10 8.5a.5.5 0 11-1 0 .5.5 0 011 0zm5 5a.5.5 0 11-1 0 .5.5 0 011 0z" />
    </svg>
  `
}

const stats = [
  {
    icon: AlertTriangle,
    value: 73,
    initialValue: 0,
    unit: "%",
    label: "of students feel overwhelmed by applying to colleges",
    duration: 1500
  },
  {
    icon: Clock,
    value: 120,
    initialValue: 0,
    unit: "+ hours",
    label: "wasted on university research",
    duration: 2000
  },
  {
    icon: Calendar,
    value: 56,
    initialValue: 0,
    unit: "%",
    label: "of students have no strategy in preparing for exams",
    duration: 1500
  },
  {
    icon: Percent,
    value: 3.7,
    initialValue: 0,
    unit: "x",
    label: "higher rejection rates for unpolished applications",
    duration: 1800,
    decimalPlaces: 1
  }
]

const counterEls = ref([])

onMounted(() => {
  // Start counters when component mounts
  stats.forEach((stat, index) => {
    animateCounter(
      counterEls.value[index],
      stat.value,
      stat.duration,
      stat.decimalPlaces
    )
  })
})

function animateCounter(el, target, duration, decimalPlaces = 0) {
  const start = 0
  const increment = target / (duration / 16) // 60fps
  let current = start
  const startTime = performance.now()

  function updateCounter(timestamp) {
    const elapsed = timestamp - startTime
    const progress = Math.min(elapsed / duration, 1)

    if (decimalPlaces > 0) {
      current = Number((target * progress).toFixed(decimalPlaces))
    } else {
      current = Math.floor(target * progress)
    }

    el.textContent = current

    if (progress < 1) {
      requestAnimationFrame(updateCounter)
    } else {
      el.textContent = decimalPlaces > 0 ? target.toFixed(decimalPlaces) : target
    }
  }

  requestAnimationFrame(updateCounter)
}
</script>

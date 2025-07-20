<template>
  <section class="mb-12 px-4 bg-white">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-3xl md:text-4xl font-bold mb-16 text-center">
        The Application <span class="text-primary-600">Struggle</span> Is Real
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <div v-for="(stat, index) in stats" :key="index" :delay="0.1 * index">
          <div class="animate-fade-in-up p-6 rounded-xl text-center h-full">
            <div class="text-primary-600 mb-4 flex justify-center">
              <Icon :icon="stat.icon" class="h-10 w-10" />
            </div>
            <h3 class="text-3xl font-bold mb-2 text-primary-600">
              <span ref="counterEls">{{ stat.initialValue }}</span>{{ stat.unit }}
            </h3>
            <p class="text-gray-600">{{ stat.label }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { Icon } from "@iconify/vue"

interface Stat {
  icon: string
  value: number
  initialValue: number
  unit: string
  label: string
  duration: number
  decimalPlaces?: number
}

const stats: Stat[] = [
  {
    icon: "solar:sad-circle-outline",
    value: 73,
    initialValue: 0,
    unit: "%",
    label: "of students feel overwhelmed by applying to colleges",
    duration: 1500,
  },
  {
    icon: "game-icons:extra-time",
    value: 120,
    initialValue: 0,
    unit: "+ hours",
    label: "wasted on university research",
    duration: 2000,
  },
  {
    icon: "streamline-freehand:edit-pen-write-paper",
    value: 56,
    initialValue: 0,
    unit: "%",
    label: "of students have no strategy in preparing for exams",
    duration: 1500,
  },
  {
    icon: "pepicons-pencil:cv-circle-off",
    value: 3.7,
    initialValue: 0,
    unit: "x",
    label: "higher rejection rates for unpolished applications",
    duration: 1800,
    decimalPlaces: 1,
  },
]

const counterEls = ref<HTMLElement[]>([])

onMounted(() => {
  stats.forEach((stat, index) => {
    if (counterEls.value[index]) {
      animateCounter(
        counterEls.value[index],
        stat.value,
        stat.duration,
        stat.decimalPlaces
      )
    }
  })
})

function animateCounter(
  el: HTMLElement,
  target: number,
  duration: number,
  decimalPlaces: number = 0
): void {
  const start: number = 0
  let current: number = start
  const startTime: DOMHighResTimeStamp = performance.now()

  function updateCounter(timestamp: DOMHighResTimeStamp): void {
    const elapsed: number = timestamp - startTime
    const progress: number = Math.min(elapsed / duration, 1)

    if (decimalPlaces > 0) {
      current = Number((target * progress).toFixed(decimalPlaces))
    } else {
      current = Math.floor(target * progress)
    }

    el.textContent = current.toString()

    if (progress < 1) {
      requestAnimationFrame(updateCounter)
    } else {
      el.textContent = decimalPlaces > 0 ? target.toFixed(decimalPlaces) : target.toString()
    }
  }

  requestAnimationFrame(updateCounter)
}
</script>

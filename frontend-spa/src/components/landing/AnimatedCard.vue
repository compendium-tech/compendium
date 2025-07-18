<template>
  <div class="fade-in" :style="{ 'animation-delay': `${delay}s` }" :class="{ 'appear': appear }"
    @vue:mounted="startAnimation">
    <slot />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  delay: {
    type: Number,
    default: 0
  }
})

const appear = ref(false)

const startAnimation = () => {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        setTimeout(() => {
          appear.value = true
        }, props.delay * 300)
        observer.unobserve(entry.target)
      }
    })
  }, { threshold: 0.1 })

  observer.observe(document.querySelector('.fade-in'))
}
</script>

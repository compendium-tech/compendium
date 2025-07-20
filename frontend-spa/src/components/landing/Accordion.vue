<template>
  <div class="space-y-4">
    <div v-for="(item, index) in items" :key="index"
      class="border border-gray-200 rounded-lg overflow-hidden shadow-sm">
      <button
        class="flex justify-between items-center w-full p-4 bg-gray-50 hover:bg-gray-100 transition-colors duration-200 ease-in-out"
        @click="toggleItem(index)" :aria-expanded="openIndex === index ? 'true' : 'false'"
        :aria-controls="`accordion-content-${index}`">
        <span class="text-xl text-left flex-grow font-medium text-gray-700">{{ item.title }}</span>
        <svg :class="`w-6 h-6 text-gray-600 transform transition-transform duration-200 ease-in-out ${openIndex === index ? 'rotate-180' : ''
          }`" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
        </svg>
      </button>

      <div v-if="openIndex === index" :id="`accordion-content-${index}`" role="region"
        class="p-4 bg-white text-gray-600 border-t border-gray-200 animate-fade-in text-lg">
        {{ item.content }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export interface AccordionItem {
  title: string
  content: string
}
</script>

<script setup lang="ts">
import { ref } from "vue"

interface AccordionItems {
  items: AccordionItem[]
}

defineProps<AccordionItems>()

const openIndex = ref(-1)

const toggleItem = (index: number) => {
  openIndex.value = openIndex.value === index ? -1 : index
};
</script>

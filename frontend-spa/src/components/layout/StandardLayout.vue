<template>
  <main>
    <StandardHeader />
    <slot />
    <StandardFooter />
  </main>
</template>

<script setup>
import StandardHeader from "../layout/StandardHeader.vue"
import StandardFooter from "../layout/StandardFooter.vue"
</script>

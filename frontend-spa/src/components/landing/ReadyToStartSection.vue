<template>
  <section class="py-20 bg-gradient-to-r from-primary-400 to-primary-600">
    <div class="max-w-4xl mx-auto text-center px-4">
      <h2 class="md:text-4xl text-3xl font-bold text-white mb-6">
        Ready to transform your application journey?
      </h2>
      <p class="text-xl text-primary-100 mb-10 max-w-2xl mx-auto">
        Join thousands of students who've secured spots in their dream schools.
      </p>
      <div class="flex flex-col sm:flex-row justify-center gap-4">
        <RouterLink to="/auth/signin">
          <BaseButton variant="outline" class="text-white hover:text-primary-400" size="md">
            Let's do this!
          </BaseButton>
        </RouterLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router"
import BaseButton from "../ui/BaseButton.vue"
</script>

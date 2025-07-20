<template>
  <div class="bg-gray-900 text-gray-300 pt-20 pb-12 relative overflow-hidden">
    <div class="absolute top-0 left-0 w-full h-20 overflow-hidden">
      <svg class="absolute top-0 left-0 w-full h-20" viewBox="0 0 1200 120" preserveAspectRatio="none">
        <path
          d="M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z"
          class="fill-primary-900 opacity-20" />
        <path
          d="M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z"
          class="fill-primary-700 opacity-30" />
      </svg>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-10">
        <div class="lg:col-span-2">
          <RouterLink to="/">
            <div class="flex items-center mb-6">
              <div class="w-13 h-13 rounded-full flex items-center justify-center mr-4">
                <img src="/icon.svg" alt="Compendium Logo" />
              </div>
              <div>
                <h2 class="text-2xl font-semibold text-white">compendium<span class="text-primary-400">.</span></h2>
                <p class="text-primary-400">Your AI Education Partner</p>
              </div>
            </div>
          </RouterLink>
          <p class="mb-6 max-w-sm">
            AI-powered university application assistant helping students worldwide get into their dream schools.
          </p>

          <div class="flex space-x-4 mb-8">
            <a v-for="social in socialLinks" :key="social.name" :href="social.link" target="_blank"
              :aria-label="social.name"
              class="bg-gray-800 hover:bg-primary-600 p-3 rounded-full transition-all duration-300 transform hover:-translate-y-1">
              <Icon :icon="social.icon" class="h-5 w-5" />
            </a>
          </div>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-white mb-6">Solutions</h3>
          <ul class="space-y-4">
            <li v-for="solution in solutionLinks" :key="solution.name">
              <a :href="solution.link" class="hover:text-primary-400 transition-colors duration-300 flex items-center">
                <Icon icon="mdi:chevron-right" class="h-4 w-4 mr-2 text-primary-500" />
                {{ solution.name }}
              </a>
            </li>
          </ul>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-white mb-6">Resources</h3>
          <ul class="space-y-4">
            <li v-for="resource in resourceLinks" :key="resource.name">
              <a :href="resource.link" class="hover:text-primary-400 transition-colors duration-300 flex items-center">
                <Icon icon="mdi:chevron-right" class="h-4 w-4 mr-2 text-primary-500" />
                {{ resource.name }}
              </a>
            </li>
          </ul>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-white mb-6">Stay Updated</h3>
          <p class="mb-4">Subscribe to our newsletter for application tips and scholarship opportunities.</p>
          <form @submit.prevent="subscribe" class="flex flex-col space-y-3">
            <input type="email" v-model="email" placeholder="Your email address" required
              class="bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
            <button type="submit"
              class="bg-primary-600 hover:bg-primary-700 text-white font-medium py-3 px-6 rounded-xl transition-colors duration-300">
              Subscribe
            </button>
          </form>
        </div>
      </div>

      <div class="border-t border-gray-800 mt-16 pt-8 flex flex-col md:flex-row justify-between items-center">
        <p class="text-gray-500 text-sm">© 2025 Compendium. All rights reserved.</p>
        <div class="flex space-x-6 mt-4 md:mt-0">
          <a v-for="policy in policyLinks" :key="policy.name" :href="policy.link"
            class="text-gray-500 hover:text-gray-300 text-sm">
            {{ policy.name }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue"
import { Icon } from "@iconify/vue"

const email = ref("")

const subscribe = () => {
  if (email.value) {
    alert(`Subscribing with: ${email.value}`)
    email.value = ""
  }
}

const socialLinks = ref([
  { name: "Facebook", icon: "mdi:facebook", link: "#!" },
  { name: "Twitter", icon: "mdi:twitter", link: "#!" },
  { name: "Instagram", icon: "mdi:instagram", link: "https://www.instagram.com/compendium_tech/" },
  { name: "LinkedIn", icon: "mdi:linkedin", link: "https://www.linkedin.com/company/ccompendium" },
])

const solutionLinks = ref([
  { name: "University Research", link: "#!" },
  { name: "Essay Checking", link: "#!" },
  { name: "Scholarship Matching", link: "#!" },
  { name: "Application Tracking", link: "#!" },
  { name: "Interview Preparation", link: "#!" },
])

const resourceLinks = ref([
  { name: "Blog & Guides", link: "#!" },
  { name: "Success Stories", link: "#!" },
  { name: "University Database", link: "#!" },
  { name: "Scholarship Directory", link: "#!" },
  { name: "Webinars & Events", link: "#!" },
])

const policyLinks = ref([
  { name: "Privacy Policy", link: "#!" },
  { name: "Terms of Service", link: "#!" },
  { name: "Cookie Policy", link: "#!" },
])
</script>

<style scoped>
a {
  transition: transform 0.3s ease, color 0.3s ease
}

a:hover {
  transform: translateX(5px)
}
</style>

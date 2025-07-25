<template>
  <StandardLayout>
    <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-4xl w-full bg-white p-8 space-y-8 rounded-lg">
        <div class="text-center">
          <h1 class="text-4xl font-extrabold text-gray-900">Your Dashboard</h1>
          <p class="mt-2 text-lg text-gray-600">Manage your account and active sessions.</p>
        </div>

        <AccountInfo />

        <div class="mt-8">
          <ActiveSessions />
        </div>
      </div>
    </div>
  </StandardLayout>
</template>

<script setup lang="ts">
import { onMounted, Ref, ref } from 'vue'
import StandardLayout from '../layout/StandardLayout.vue'
import AccountInfo from './AccountInfoSection.vue'
import ActiveSessions from './ActiveSessionsSection.vue'

const globalError: Ref<string | null> = ref(null)

const w: any = window
const paddle: any = w.Paddle

const handleSubscribe = () => {
  if (paddle) {
    paddle.Checkout.open({
      settings: {
        displayMode: 'overlay',
        theme: 'light',
        locale: 'en',
        allowLogout: false,
      },
      product: 'pro_01k0qbrebvvgfwacb4qqgyh0bx',
      customer: {
        email: '',
      },
      items: [
        {
          priceID: 'pri_01k0qbs1mgx0dnjd0zytj23zm7',
          quantity: 1,
        },
      ],
      successCallback: (data: any) => {
        console.log('Paddle checkout successful:', data)
      },
      closeCallback: () => {
        console.log('Paddle checkout closed.')
      },
    })
  } else {
    console.error('Paddle.js not loaded!')
    globalError.value = 'Billing service is not available. Please try again later.'
  }
}

onMounted(() => {
  if (paddle) {
    paddle.Environment.set('sandbox')
    paddle.Setup({
      token: 'test_bf5c18ea62fd1d30c00bc5c2821',
      debug: true,
    })
  } else {
    console.error('Paddle.js script not found. Make sure it is included in your index.html')
    globalError.value = 'Billing service initialization failed.'
  }
})
</script>

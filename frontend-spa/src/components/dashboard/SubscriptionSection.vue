<template>
  <div class="mt-8">
    <h2 class="text-2xl font-semibold text-gray-800 mb-4">Your Subscription</h2>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-6">
      <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
      <p class="mt-2 text-md text-primary-600">Loading subscription info...</p>
    </div>
    <div v-if="subscriptionError" class="text-center text-red-600">{{ subscriptionError }}</div>
    <div v-if="globalError" class="text-center text-red-600">{{ globalError }}</div>

    <div v-if="subscriptionResponse && subscriptionResponse.subscription && !isLoading">
      <div class="bg-gray-50 rounded-lg p-6 shadow-sm">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Subscription Level</label>
            <p class="mt-1 text-lg text-gray-900">{{ formatTier(subscriptionResponse.subscription.tier) }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Active Since</label>
            <p class="mt-1 text-lg text-gray-900">{{ formatSubscriptionDate(subscriptionResponse.subscription.since) }}
            </p>
          </div>
          <div class="col-span-1 md:col-span-2">
            <label class="block text-sm font-medium text-gray-700">Renews / Expires</label>
            <p class="mt-1 text-lg text-gray-900">{{ formatSubscriptionDate(subscriptionResponse.subscription.till) }}
            </p>
          </div>
        </div>
        <div class="mt-6 text-right">
          <BaseButton @click="handleCancelSubscription" variant="secondary" :is-loading="isCancelling">
            Cancel Subscription
          </BaseButton>
        </div>
      </div>
    </div>

    <div v-if="!subscriptionResponse || !subscriptionResponse.isActive && !isLoading && !subscriptionError">
      <div class="my-8">
        <p class="text-lg text-gray-700">You don't have an active subscription. Choose a plan to get started!</p>
      </div>
      <section class="py-20 px-4">
        <div class="max-w-6xl mx-auto">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div v-for="(plan, index) in pricing" :key="index">
              <div :class="`bg-white rounded-2xl overflow-hidden shadow-lg border-3 h-full flex flex-col ${plan.highlight
                ? 'border-primary-600 md:scale-105 z-10'
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

                  <BaseButton :variant="plan.highlight ? 'primary' : 'secondary'" @click="handleSubscribe(plan.name)"
                    class="mt-auto">
                    Get Started
                  </BaseButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, Ref, ref } from 'vue'
import BaseButton from '../ui/BaseButton.vue'
import { subscriptionService, SubscriptionResponse, Tier } from '../../api/subscription.ts'

const subscriptionResponse: Ref<SubscriptionResponse | null> = ref(null)
const isLoading = ref(true)
const subscriptionError: Ref<string | null> = ref(null)
const globalError: Ref<string | null> = ref(null)
const isCancelling = ref(false)

const w: any = window
const paddle: any = w.Paddle

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

const productMap: Record<string, { monthly: string, annually: string }> = {
  Student: {
    monthly: 'pri_01k0qbs1mgx0dnjd0zytj23zm7',
    annually: 'pri_01k0qbt9jwq824bhec8edv4gh1',
  },
  Team: {
    monthly: 'pri_01k0qbwbhpa8jzs3z21md7ytx9',
    annually: 'pri_01k0qbx136dka17x086crra4kq',
  },
  Community: {
    monthly: 'pri_01k0qbytrbsfdft9ty91bng7sr',
    annually: 'pri_01k0qbzs3rch23hx6p9ge1sa5b',
  },
};


const fetchSubscription = async () => {
  isLoading.value = true
  subscriptionError.value = null
  try {
    const data = await subscriptionService.getSubscription()
    subscriptionResponse.value = data
  } catch (err: any) {
    subscriptionError.value = err.message || 'Failed to fetch subscription details.'
    subscriptionResponse.value = { isActive: false };
  } finally {
    isLoading.value = false
  }
}

const formatTier = (tier: Tier): string => {
  switch (tier) {
    case 'student':
      return 'Student Tier'
    case 'team':
      return 'Team Tier'
    case 'community':
      return 'Community Tier'
    default:
      return 'Unknown Tier'
  }
}

const formatSubscriptionDate = (dateString: string): string => {
  try {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  } catch (e) {
    return 'N/A'
  }
}

const handleSubscribe = (planName: string) => {
  const priceId = productMap[planName]?.monthly;

  if (!priceId) {
    console.error('Invalid plan name or price ID not found.');
    globalError.value = 'Could not initiate checkout for the selected plan.';
    return;
  }

  if (paddle) {
    paddle.Checkout.open({
      settings: {
        displayMode: 'overlay',
        theme: 'light',
        locale: 'en',
        allowLogout: false,
      },
      items: [
        {
          priceId: priceId,
          quantity: 1,
        },
      ],
      successCallback: (data: any) => {
        console.log('Paddle checkout successful:', data);
        fetchSubscription();
      },
      closeCallback: () => {
        console.log('Paddle checkout closed.');
      },
    })
  } else {
    console.error('Paddle.js not loaded!');
    globalError.value = 'Billing service is not available. Please try again later.';
  }
}

const handleCancelSubscription = async () => {
  isCancelling.value = true;
  try {
    await subscriptionService.cancelSubscription()
    await fetchSubscription()
  } catch (err: any) {
    subscriptionError.value = err.message || 'Failed to cancel subscription.';
  } finally {
    isCancelling.value = false;
  }
}

onMounted(() => {
  fetchSubscription()

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

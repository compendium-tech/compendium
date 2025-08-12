<template>
  <h2 class="text-2xl font-semibold text-gray-800">Your Subscription</h2>

  <div v-if="globalError" class="text-red-600">{{ globalError }}</div>

  <div v-if="isLoading.subscription" class="flex flex-col items-center justify-center">
    <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
    <p class="mt-2 text-md text-primary-600">Loading subscription info...</p>
  </div>

  <div v-if="subscriptionResponse?.subscription && !isLoading.subscription">
    <div class="bg-gray-50 rounded-lg p-6 shadow-sm">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700">Subscription Level</label>
          <p class="mt-1 text-lg text-gray-900">{{ formatTier(subscriptionResponse.subscription.tier) }}</p>
        </div>
        <div v-if="!isPayer">
          <label class="block text-sm font-medium text-gray-700">Paid by</label>
          <p class="mt-1 text-lg text-gray-900" v-if="payerInfo">{{ payerInfo.name }} &lt;{{ payerInfo.email }}&gt;</p>
          <p v-else class="mt-1 text-lg text-gray-900">Unknown</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Active Since</label>
          <p class="mt-1 text-lg text-gray-900">{{ dateToString(subscriptionResponse.subscription.since) }}</p>
        </div>
        <div class="col-span-1 md:col-span-2">
          <label class="block text-sm font-medium text-gray-700">Renews / Expires</label>
          <p class="mt-1 text-lg text-gray-900">{{ dateToString(subscriptionResponse.subscription.till) }}</p>
        </div>
      </div>
      <div class="mt-6 text-right" v-if="isPayer">
        <BaseButton class="flex space-x-2" @click="showCancelConfirm = true" variant="outline"
          :is-loading="isLoading.cancellingSubscription" size="sm" hover-effect="scale">
          <Icon icon="lets-icons:cancel" class="h-5 w-5" />
          <span>Cancel Subscription</span>
        </BaseButton>
      </div>

      <div v-if="showCancelConfirm"
        class="fixed inset-0 bg-opacity-50 backdrop-blur-sm overflow-y-auto h-full w-full flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Confirm Cancellation</h3>
          <p class="text-gray-600 mb-6">Are you sure you want to cancel your subscription? This action cannot be undone.
          </p>
          <div class="flex justify-end space-x-4">
            <BaseButton class="flex space-x-2" variant="secondary" @click="showCancelConfirm = false" size="sm">
              <Icon icon="arcticons:no-fap" class="h-5 w-5" />
              <span>Cancel</span>
            </BaseButton>
            <BaseButton class="flex space-x-2" @click="handleCancelSubscription"
              :is-loading="isLoading.cancellingSubscription" size="sm">
              <Icon icon="line-md:circle-to-confirm-circle-transition" class="h-5 w-5" />
              <span>Confirm</span>
            </BaseButton>
          </div>
        </div>
      </div>

      <div v-if="hasTeamFeatures" class="mt-8 pt-6 border-t border-gray-200">
        <div class="mb-6">
          <h3 class="text-xl font-medium text-gray-800 mb-4">Invitation Code</h3>
          <div v-if="isLoading.invitationCode" class="flex items-center space-x-2 text-primary-600">
            <div class="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-primary-600"></div>
            <span>Loading code...</span>
          </div>
          <div v-else-if="invitationCode">
            <div class="flex flex-col md:flex-row md:items-center space-y-2 md:space-y-0 space-x-4 mb-4">
              <p class="text-2xl font-bold tracking-widest text-gray-700 bg-gray-100 px-4 py-2 rounded-lg">
                {{ invitationCode }}
              </p>
              <div class="flex space-x-2">
                <BaseButton variant="outline" class="flex space-x-2" @click="updateInvitationCode"
                  :is-loading="isLoading.updatingInvitationCode" hover-effect="scale" style="min-width: 0;">
                  <Icon icon="mdi:refresh" class="h-5 w-5" />
                  <span>Update</span>
                </BaseButton>
                <BaseButton variant="secondary" class="flex space-x-2" @click="removeInvitationCode"
                  :is-loading="isLoading.removingInvitationCode" style="min-width: 0;">
                  <Icon icon="mdi:delete" class="h-5 w-5" />
                  <span>Remove</span>
                </BaseButton>
              </div>
            </div>
          </div>
          <div v-else>
            <p class="text-gray-600 mb-4">No active invitation code. Generate one to invite members to your plan.</p>
            <BaseButton @click="updateInvitationCode" class="flex space-x-2"
              :is-loading="isLoading.updatingInvitationCode" size="sm">
              <Icon icon="mdi:refresh" class="h-5 w-5" />
              <span>Generate</span>
            </BaseButton>
          </div>
        </div>

        <div v-if="subscriptionResponse.subscription.members?.length">
          <h3 class="text-xl font-medium text-gray-800 mb-4">Members</h3>
          <ul class="space-y-4">
            <li v-for="member in subscriptionResponse.subscription.members" :key="member.userId"
              class="flex items-center justify-between p-4 bg-gray-100 rounded-md">
              <div class="flex items-center space-x-2">
                <p class="font-semibold text-gray-900">
                  {{ member.name }} <span v-if="member.email">&lt;{{ member.email }}&gt;</span>
                  <span v-if="member.role === 'payer'"
                    class="ml-2 text-xs font-medium text-primary-600 bg-primary-100 px-2 py-1 rounded-full">
                    payer
                  </span>
                </p>
              </div>
              <div class="flex items-center space-x-2">
                <span v-if="member.role === 'member'">
                  <BaseButton @click="handleRemoveMember(member.userId)" variant="secondary" size="sm"
                    :is-loading="isLoading.idOfMemberBeingRemoved === member.userId">
                    Remove
                  </BaseButton>
                </span>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>

  <div v-if="!subscriptionResponse?.isActive && !isLoading.subscription">
    <p class="text-lg text-gray-700">You don't have an active subscription. Choose a plan to get started!</p>

    <div class="mt-8 bg-white rounded-lg p-6 shadow-sm border border-primary-200">
      <h3 class="text-xl font-medium text-gray-800 mb-4">Join a collective plan</h3>
      <p class="text-gray-600 mb-4">If you've been invited to a collective plan, enter your invitation code below.</p>
      <div class="flex items-center space-x-2">
        <BaseInput v-model="joinInvitationCode" type="text" placeholder="Enter invitation code"
          :disabled="isLoading.joiningSubscription" />
        <BaseButton @click="handleJoinSubscription" class="flex space-x-2" :is-loading="isLoading.joiningSubscription"
          size="sm">
          <Icon icon="mdi:human-hello-variant" class="h-5 w-5" />
          <span>Join</span>
        </BaseButton>
      </div>
    </div>

    <section class="py-20 px-4">
      <div class="max-w-6xl mx-auto">
        <div class="flex justify-center mb-8">
          <div class="inline-flex rounded-md shadow-sm" role="group">
            <button type="button" @click="selectedBillingCycle = 'monthly'"
              :class="['px-4 py-2 text-sm font-medium border rounded-l-lg', selectedBillingCycle === 'monthly' ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-gray-900 border-gray-200 hover:bg-gray-100']">
              Monthly
            </button>
            <button type="button" @click="selectedBillingCycle = 'annually'"
              :class="['px-4 py-2 text-sm font-medium border rounded-r-lg', selectedBillingCycle === 'annually' ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-gray-900 border-gray-200 hover:bg-gray-100']">
              Yearly
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-0 md:gap-x-8 gap-y-8">
          <div v-for="(plan, index) in pricing" :key="index">
            <div
              :class="`bg-white rounded-2xl overflow-hidden shadow-lg border-3 h-full flex flex-col ${plan.highlight ? 'border-primary-600 md:scale-105 z-10' : 'border-gray-200'}`">
              <div v-if="plan.highlight" class="bg-primary-600 text-white text-center py-2">
                Most Popular
              </div>
              <div class="p-8 flex-grow flex flex-col">
                <h3 class="text-2xl font-bold mb-2">{{ plan.name }}</h3>
                <div class="flex items-baseline mb-4">
                  <span class="text-5xl font-bold">
                    {{ selectedBillingCycle === 'monthly' ? plan.priceMonthly : plan.priceYearly }}
                  </span>
                  <span class="text-gray-500">/month</span>
                </div>
                <p class="text-gray-600 mb-6">{{ plan.description }}</p>
                <ul class="space-y-4 mb-8 flex-grow">
                  <li v-for="(feature, i) in plan.features" :key="i" class="flex space-x-2 items-start">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500 flex-shrink-0 mt-0.5"
                      fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <span>{{ feature }}</span>
                  </li>
                </ul>
                <BaseButton :variant="plan.highlight ? 'primary' : 'secondary'"
                  @click="handleSubscribe(plan.name, selectedBillingCycle)" class="mt-auto" size="md">
                  Get Started
                </BaseButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, Ref } from 'vue'
import { userService, AccountDetails } from '../../api/user.ts'
import { subscriptionService, SubscriptionResponse, Tier } from '../../api/subscription.ts'
import { Icon } from '@iconify/vue'
import BaseInput from '../ui/BaseInput.vue'
import BaseButton from '../ui/BaseButton.vue'
import { ApiError } from '../../api/base.ts'
import { dateToString } from '../../utils/date'
import { pricing } from '../../pricing.ts'

interface ProductMap {
  [key: string]: {
    monthly: string
    annually: string
  }
}

const user: Ref<AccountDetails | null> = ref(null)
const subscriptionResponse: Ref<SubscriptionResponse | null> = ref(null)
const globalError = ref('')
const showCancelConfirm = ref(false)
const isLoading = ref({
  account: true,
  subscription: true,
  invitationCode: false,
  updatingInvitationCode: false,
  removingInvitationCode: false,
  joiningSubscription: false,
  cancellingSubscription: false,
  idOfMemberBeingRemoved: ''
})

const selectedBillingCycle: Ref<'monthly' | 'annually'> = ref('monthly')
const invitationCode: Ref<string | undefined> = ref(undefined)
const joinInvitationCode = ref('')

const productMap: ProductMap = {
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
    annually: 'pri_01k13a0e9hda6hmjck2ecn0jq9',
  },
};

const isPayer = computed(() => subscriptionResponse.value?.subscription?.role === 'payer')
const payerInfo = computed(() => subscriptionResponse.value?.subscription?.members?.find(member => member.role === 'payer'))
const hasTeamFeatures = computed(() => isPayer.value && (subscriptionResponse.value?.subscription?.tier === 'team' || subscriptionResponse.value?.subscription?.tier === 'community'))

const fetchUserData = async () => {
  isLoading.value.account = true
  globalError.value = ''

  try {
    user.value = await userService.getAccountDetails()
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.account = false
  }
}

const fetchSubscription = async () => {
  isLoading.value.subscription = true
  globalError.value = ''

  try {
    subscriptionResponse.value = await subscriptionService.getSubscription()
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
    subscriptionResponse.value = { isActive: false }
  } finally {
    isLoading.value.subscription = false
  }
}

const fetchInvitationCode = async () => {
  isLoading.value.invitationCode = true
  globalError.value = ''

  try {
    const data = await subscriptionService.getSubscriptionInvitationCode()
    invitationCode.value = data.invitationCode
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
    invitationCode.value = undefined
  } finally {
    isLoading.value.invitationCode = false
  }
}

const updateInvitationCode = async () => {
  isLoading.value.updatingInvitationCode = true
  globalError.value = ''

  try {
    const data = await subscriptionService.updateSubscriptionInvitationCode()
    invitationCode.value = data.invitationCode
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.updatingInvitationCode = false
  }
}

const removeInvitationCode = async () => {
  isLoading.value.removingInvitationCode = true
  globalError.value = ''

  try {
    await subscriptionService.removeSubscriptionInvitationCode()
    invitationCode.value = undefined
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.removingInvitationCode = false
  }
}

const handleJoinSubscription = async () => {
  if (!joinInvitationCode.value) {
    globalError.value = 'Please enter an invitation code.'
    return
  }

  isLoading.value.joiningSubscription = true
  globalError.value = ''

  try {
    subscriptionResponse.value = await subscriptionService.joinSubscription(joinInvitationCode.value)
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.joiningSubscription = false
  }
}

const handleRemoveMember = async (memberId: string) => {
  isLoading.value.idOfMemberBeingRemoved = memberId
  globalError.value = ''

  try {
    await subscriptionService.removeSubscriptionMember(memberId)
    if (subscriptionResponse.value?.subscription?.members) {
      subscriptionResponse.value.subscription.members = subscriptionResponse.value.subscription.members.filter(
        (member) => member.userId !== memberId
      )
    }
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.idOfMemberBeingRemoved = ''
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

const handleSubscribe = (planName: string, billingCycle: 'monthly' | 'annually') => {
  const priceId = productMap[planName]?.[billingCycle]

  if (!priceId) {
    globalError.value = 'Could not initiate checkout for the selected plan and billing cycle.'
    return
  }

  const paddle = (window as any).Paddle
  if (paddle) {
    paddle.Checkout.open({
      settings: {
        displayMode: 'overlay',
        theme: 'light',
        locale: 'en',
        allowLogout: false,
        variant: 'one-page'
      },
      items: [
        {
          priceId,
          quantity: 1,
        },
      ],
      customer: {
        email: user.value?.email,
      },
      customData: {
        userId: user.value?.id,
      },
      successCallback: (data: unknown) => {
        console.log('Paddle checkout successful:', data)
        fetchSubscription()
      },
      closeCallback: () => {
        console.log('Paddle checkout closed.')
      },
    })
  } else {
    globalError.value = 'Billing service is not available. Please try again later.'
  }
}

const handleCancelSubscription = async () => {
  isLoading.value.cancellingSubscription = true
  globalError.value = ''

  try {
    await subscriptionService.cancelSubscription()
    await fetchSubscription()
  } catch (error) {
    if (error instanceof ApiError) {
      globalError.value = error.message
    }
  } finally {
    isLoading.value.cancellingSubscription = false
    showCancelConfirm.value = false
  }
}

onMounted(async () => {
  await Promise.all([fetchUserData(), fetchSubscription()])

  if (hasTeamFeatures.value) {
    await fetchInvitationCode()
  }

  const paddle = (window as any).Paddle
  if (paddle) {
    paddle.Environment.set('sandbox')
    paddle.Setup({
      token: 'test_bf5c18ea62fd1d30c00bc5c2821',
      debug: true,
    })
  } else {
    globalError.value = 'Billing service initialization failed.'
  }
})
</script>

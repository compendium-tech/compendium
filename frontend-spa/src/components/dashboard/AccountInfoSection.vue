<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-800 mb-4">Account Information</h2>

    <div v-if="isLoadingAccount" class="flex flex-col items-center justify-center py-6">
      <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
      <p class="mt-2 text-md text-primary-600">Loading account info...</p>
    </div>
    <div v-if="accountError" class="text-center text-red-600">{{ accountError }}</div>

    <div v-if="user && !isLoadingAccount">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">Email Address</label>
          <p class="mt-1 text-lg text-gray-900">{{ user.email }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Account Created</label>
          <p class="mt-1 text-lg text-gray-900">{{ formattedCreationDate }}</p>
        </div>
        <div class="col-span-1 md:col-span-2">
          <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
          <div class="mt-1 flex items-center space-x-2" v-if="!isEditingName">
            <span class="text-lg">{{ user.name }}</span>
            <BaseButton variant="outline" class="px-2 flex items-center" @click="toggleEditName" size="sm"
              hover-effect="scale">
              <Icon icon="mdi:pencil" class="h-5 w-5 mr-2" />
              <span>Edit</span>
            </BaseButton>
          </div>
          <div class="mt-1 flex items-center space-x-2" v-else>
            <BaseInput v-model="editableName" type="text" id="name" :disabled="isSavingName" />
            <BaseButton class="px-4 flex items-center" @click="toggleEditName" size="sm" hover-effect="scale"
              :is-loading="isSavingName">
              <Icon icon="mdi:content-save" class="h-5 w-5 mr-2" />
              <span>Save</span>
            </BaseButton>
          </div>
          <BaseTransitioningText>
            <p v-if="nameUpdateMessage" :class="nameUpdateSuccess ? 'text-green-600' : 'text-red-600'"
              class="mt-2 text-sm">
              {{ nameUpdateMessage }}
            </p>
          </BaseTransitioningText>
        </div>
      </div>
    </div>

    <div class="mt-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Your Subscription</h2>

      <div v-if="isLoadingSubscription" class="flex flex-col items-center justify-center py-6">
        <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
        <p class="mt-2 text-md text-primary-600">Loading subscription info...</p>
      </div>
      <div v-if="subscriptionError" class="text-center text-red-600">{{ subscriptionError }}</div>
      <div v-if="globalError" class="text-center text-red-600">{{ globalError }}</div>

      <div v-if="subscriptionResponse && subscriptionResponse.subscription && !isLoadingSubscription">
        <div class="bg-gray-50 rounded-lg p-6 shadow-sm">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Subscription Level</label>
              <p class="mt-1 text-lg text-gray-900">{{ formatTier(subscriptionResponse.subscription.tier) }}</p>
            </div>
            <div v-if="!isPayer">
              <label class="block text-sm font-medium text-gray-700">Paid by</label>
              <p class="mt-1 text-lg text-gray-900" v-if="payerInfo">{{ payerInfo.name }} &lt;{{ payerInfo.email }}&gt;
              </p>
              <p v-else class="mt-1 text-lg text-gray-900">Unknown</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Active Since</label>
              <p class="mt-1 text-lg text-gray-900">{{ formatSubscriptionDate(subscriptionResponse.subscription.since)
              }}
              </p>
            </div>
            <div class="col-span-1 md:col-span-2">
              <label class="block text-sm font-medium text-gray-700">Renews / Expires</label>
              <p class="mt-1 text-lg text-gray-900">{{ formatSubscriptionDate(subscriptionResponse.subscription.till) }}
              </p>
            </div>
          </div>
          <div class="mt-6 text-right" v-if="isPayer">
            <BaseButton @click="handleCancelSubscription" variant="secondary" :is-loading="isCancelling">
              Cancel Subscription
            </BaseButton>
          </div>

          <div v-if="hasTeamFeatures" class="mt-8 pt-6 border-t border-gray-200">
            <div class="mb-6">
              <h3 class="text-xl font-medium text-gray-800 mb-4">Invitation Code</h3>
              <div v-if="isLoadingInvitationCode" class="flex items-center space-x-2 text-primary-600">
                <div class="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-primary-600"></div>
                <span>Loading code...</span>
              </div>
              <div v-else-if="invitationCode">
                <div class="flex items-center space-x-4 mb-4">
                  <p class="text-2xl font-bold tracking-widest text-gray-700 bg-gray-100 px-4 py-2 rounded-lg">
                    {{ invitationCode }}
                  </p>
                  <BaseButton variant="outline" class="flex" size="sm" @click="updateInvitationCode"
                    :is-loading="isUpdatingInvitationCode" hover-effect="scale">
                    <Icon icon="mdi:refresh" class="h-5 w-5 mr-2" />
                    <span>Update</span>
                  </BaseButton>
                  <BaseButton variant="secondary" class="flex" size="sm" @click="removeInvitationCode"
                    :is-loading="isRemovingInvitationCode">
                    <Icon icon="mdi:delete" class="h-5 w-5 mr-2" />
                    <span>Remove</span>
                  </BaseButton>
                </div>
              </div>
              <div v-else>
                <p class="text-gray-600 mb-4">No active invitation code. Generate one to invite members to your plan.
                </p>
                <BaseButton @click="updateInvitationCode" class="flex" :is-loading="isUpdatingInvitationCode" size="sm">
                  <Icon icon="mdi:refresh" class="h-5 w-5 mr-2" />
                  Generate
                </BaseButton>
              </div>
              <div v-if="invitationCodeError" class="mt-2 text-red-600 text-sm">
                {{ invitationCodeError }}
              </div>
            </div>

            <div
              v-if="subscriptionResponse.subscription.members && subscriptionResponse.subscription.members.length > 0">
              <h3 class="text-xl font-medium text-gray-800 mb-4">Members</h3>
              <ul class="space-y-4">
                <li v-for="member in subscriptionResponse.subscription.members" :key="member.userId"
                  class="flex items-center justify-between p-4 bg-gray-100 rounded-md">
                  <div class="flex items-center justify-between space-x-2">
                    <p class="font-semibold text-gray-900">{{ member.name }} <span v-if="member.email">&lt;{{
                      member.email }}&gt;</span></p>

                    <span class="text-xs font-medium text-primary-600 bg-primary-100 px-2 py-1 rounded-full mt-1">
                      {{ member.role }}
                    </span>
                  </div>
                  <div v-if="member.role === 'member'">
                    <BaseButton @click="handleRemoveMember(member.userId)" variant="secondary" size="sm"
                      :is-loading="isRemovingMember === member.userId">
                      Remove
                    </BaseButton>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div
        v-if="!subscriptionResponse || !subscriptionResponse.isActive && !isLoadingSubscription && !subscriptionError">
        <div class="my-8">
          <p class="text-lg text-gray-700">You don't have an active subscription. Choose a plan to get started!</p>
        </div>

        <div class="mt-8 bg-white rounded-lg p-6 shadow-sm border border-primary-200">
          <h3 class="text-xl font-medium text-gray-800 mb-4">Join a collective plan</h3>
          <p class="text-gray-600 mb-4">If you've been invited to a collective plan, enter your invitation code below.
          </p>
          <div class="flex items-center space-x-2">
            <BaseInput v-model="joinInvitationCode" type="text" placeholder="Enter invitation code"
              :disabled="isJoiningSubscription" />
            <BaseButton @click="handleJoinSubscription" class="flex" :is-loading="isJoiningSubscription" size="sm">
              <Icon icon="mdi:human-hello-variant" class="h-5 w-5 mr-2" />
              Join
            </BaseButton>
          </div>
          <div v-if="joinSubscriptionError" class="mt-2 text-red-600 text-sm">
            {{ joinSubscriptionError }}
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
                <button type="button" @click="selectedBillingCycle = 'yearly'"
                  :class="['px-4 py-2 text-sm font-medium border', selectedBillingCycle === 'yearly' ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-gray-900 border-gray-200 hover:bg-gray-100']">
                  Yearly
                </button>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div v-for="(plan, index) in pricing" :key="index">
                <div
                  :class="`bg-white rounded-2xl overflow-hidden shadow-lg border-3 h-full flex flex-col ${plan.highlight ? 'border-primary-600 md:scale-105 z-10' : 'border-gray-200'}`">
                  <div v-if="plan.highlight" class=" bg-primary-600 text-white text-center py-2">
                    Most Popular
                  </div>
                  <div class="p-8 flex-grow flex flex-col">
                    <h3 class="text-2xl font-bold mb-2">{{ plan.name }}</h3>
                    <div class="flex items-baseline mb-4">
                      <span class="text-5xl font-bold">
                        {{
                          selectedBillingCycle === 'monthly'
                            ? plan.priceMonthly
                            : selectedBillingCycle === 'yearly'
                              ? plan.priceYearly
                              : plan.priceOneTime
                        }}
                      </span>
                      <span class="text-gray-500">/month</span>
                    </div>
                    <p class="text-gray-600 mb-6">{{ plan.description }}</p>

                    <ul class="space-y-4 mb-8 flex-grow">
                      <li v-for="(feature, i) in plan.features" :key="i" class="flex items-start">
                        <svg xmlns="[http://www.w3.org/2000/svg](http://www.w3.org/2000/svg)"
                          class="h-6 w-6 text-green-500 mr-2 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24"
                          stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        <span>{{ feature }}</span>
                      </li>
                    </ul>

                    <BaseButton :variant="plan.highlight ? 'primary' : 'secondary'"
                      @click="handleSubscribe(plan.name, selectedBillingCycle)" class="mt-auto">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, Ref } from 'vue';
import { userService, AccountDetails } from '../../api/user';
import { subscriptionService, SubscriptionResponse, Tier, SubscriptionRole, SubscriptionMember } from '../../api/subscription.ts';
import { Icon } from '@iconify/vue';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseTransitioningText from '../ui/BaseTransitioningText.vue';

const user: Ref<AccountDetails | null> = ref(null);
const isLoadingAccount = ref(true);
const accountError: Ref<string | null> = ref(null);
const isEditingName = ref(false);
const editableName = ref('');
const nameUpdateMessage = ref('');
const nameUpdateSuccess = ref(false);
const isSavingName = ref(false);

const subscriptionResponse: Ref<SubscriptionResponse | null> = ref(null);
const isLoadingSubscription = ref(true);
const subscriptionError: Ref<string | null> = ref(null);
const globalError: Ref<string | null> = ref(null);
const isCancelling = ref(false);

const selectedBillingCycle: Ref<'monthly' | 'yearly'> = ref('monthly');
const invitationCode: Ref<string | undefined> = ref(undefined);
const isLoadingInvitationCode = ref(true);
const invitationCodeError: Ref<string | null> = ref(null);
const isUpdatingInvitationCode = ref(false);
const isRemovingInvitationCode = ref(false);
const joinInvitationCode = ref('');
const isJoiningSubscription = ref(false);
const joinSubscriptionError: Ref<string | null> = ref(null);
const isRemovingMember: Ref<string | null> = ref(null);

const w: any = window;
const paddle: any = w.Paddle;

interface PricingCard {
  name: string;
  priceMonthly: string;
  priceYearly: string;
  priceOneTime?: string;
  description: string;
  features: string[];
  highlight: boolean;
}

const pricing: PricingCard[] = [
  {
    name: "Student",
    priceMonthly: "$5",
    priceYearly: "$2.5",
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
    priceMonthly: "$10",
    priceYearly: "$4.17",
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
    priceMonthly: "$30",
    priceYearly: "$15",
    priceOneTime: "$200",
    description: "Schools & large organizations",
    features: [
      "Everything in Pro",
      "Invite 150+ students",
      "Advanced analytics"
    ],
    highlight: false
  }
];

const productMap: Record<string, { monthly: string, annually: string, onetime?: string }> = {
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

const formattedCreationDate = computed(() => {
  if (user.value && user.value.createdAt) {
    return new Date(user.value.createdAt).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }
  return 'N/A';
});

const isPayer = computed(() => {
  return subscriptionResponse.value?.subscription?.role === 'payer';
});

const payerInfo = computed(() => {
  return subscriptionResponse.value?.subscription?.members?.find(
    member => member.role === 'payer'
  );
});

const hasTeamFeatures = computed(() => {
  const tier = subscriptionResponse.value?.subscription?.tier;
  return isPayer.value && (tier === 'team' || tier === 'community');
});

const fetchUserData = async () => {
  isLoadingAccount.value = true;
  accountError.value = null;
  try {
    const data = await userService.getAccountDetails();
    user.value = data;
    editableName.value = data.name;
  } catch (err: any) {
    accountError.value = err.message || 'Failed to fetch user data.';
  } finally {
    isLoadingAccount.value = false;
  }
};

const toggleEditName = async () => {
  if (!user.value) return;

  if (isEditingName.value) {
    isSavingName.value = true;
    try {
      const updatedUser = await userService.updateName(editableName.value);
      user.value.name = updatedUser.name;
      nameUpdateMessage.value = 'Name updated successfully!';
      nameUpdateSuccess.value = true;
    } catch (err: any) {
      nameUpdateMessage.value = 'Failed to update name: ' + (err.message || 'Unknown error');
      nameUpdateSuccess.value = false;
    } finally {
      isSavingName.value = false;
      setTimeout(() => {
        nameUpdateMessage.value = '';
      }, 3000);
    }
  }
  isEditingName.value = !isEditingName.value;
};

const fetchSubscription = async () => {
  isLoadingSubscription.value = true;
  subscriptionError.value = null;
  try {
    const data = await subscriptionService.getSubscription();
    subscriptionResponse.value = data;
  } catch (err: any) {
    subscriptionError.value = err.message || 'Failed to fetch subscription details.';
    subscriptionResponse.value = { isActive: false };
  } finally {
    isLoadingSubscription.value = false;
  }
};

const fetchInvitationCode = async () => {
  isLoadingInvitationCode.value = true;
  invitationCodeError.value = null;
  try {
    const data = await subscriptionService.getSubscriptionInvitationCode();
    invitationCode.value = data.invitationCode;
  } catch (err: any) {
    invitationCodeError.value = err.message || 'Failed to fetch invitation code.';
    invitationCode.value = undefined;
  } finally {
    isLoadingInvitationCode.value = false;
  }
};

const updateInvitationCode = async () => {
  isUpdatingInvitationCode.value = true;
  invitationCodeError.value = null;
  try {
    const data = await subscriptionService.updateSubscriptionInvitationCode();
    invitationCode.value = data.invitationCode;
  } catch (err: any) {
    invitationCodeError.value = err.message || 'Failed to update invitation code.';
  } finally {
    isUpdatingInvitationCode.value = false;
  }
};

const removeInvitationCode = async () => {
  isRemovingInvitationCode.value = true;
  invitationCodeError.value = null;
  try {
    await subscriptionService.removeSubscriptionInvitationCode();
    invitationCode.value = undefined;
  } catch (err: any) {
    invitationCodeError.value = err.message || 'Failed to remove invitation code.';
  } finally {
    isRemovingInvitationCode.value = false;
  }
};

const handleJoinSubscription = async () => {
  if (!joinInvitationCode.value) {
    joinSubscriptionError.value = 'Please enter an invitation code.';
    return;
  }
  isJoiningSubscription.value = true;
  joinSubscriptionError.value = null;
  try {
    const data = await subscriptionService.joinSubscription(joinInvitationCode.value);
    subscriptionResponse.value = data;
  } catch (err: any) {
    joinSubscriptionError.value = err.message || 'Failed to join subscription. Check the code and try again.';
  } finally {
    isJoiningSubscription.value = false;
  }
};

const handleRemoveMember = async (memberId: string) => {
  isRemovingMember.value = memberId;
  try {
    await subscriptionService.removeSubscriptionMember(memberId);
    // Optimistically update the member list
    if (subscriptionResponse.value?.subscription?.members) {
      subscriptionResponse.value.subscription.members = subscriptionResponse.value.subscription.members.filter(
        (member) => member.userId !== memberId
      );
    }
  } catch (err: any) {
    subscriptionError.value = err.message || 'Failed to remove member.';
  } finally {
    isRemovingMember.value = null;
  }
};

const formatTier = (tier: Tier): string => {
  switch (tier) {
    case 'student':
      return 'Student Tier';
    case 'team':
      return 'Team Tier';
    case 'community':
      return 'Community Tier';
    default:
      return 'Unknown Tier';
  }
};

const formatSubscriptionDate = (dateString: string): string => {
  try {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  } catch (e) {
    return 'N/A';
  }
};

const handleSubscribe = (planName: string, billingCycle: 'monthly' | 'yearly') => {
  let priceId: string | undefined;

  if (billingCycle === 'monthly') {
    priceId = productMap[planName]?.monthly;
  } else if (billingCycle === 'yearly') {
    priceId = productMap[planName]?.annually;
  }

  if (!priceId) {
    console.error('Invalid plan name or price ID not found for the selected billing cycle.');
    globalError.value = 'Could not initiate checkout for the selected plan and billing cycle.';
    return;
  }

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
          priceId: priceId,
          quantity: 1,
        },
      ],
      customer: {
        email: user.value?.email,
      },
      customData: {
        userId: user.value?.id,
      },
      successCallback: (data: any) => {
        console.log('Paddle checkout successful:', data);
        fetchSubscription();
      },
      closeCallback: () => {
        console.log('Paddle checkout closed.');
      },
    });
  } else {
    console.error('Paddle.js not loaded!');
    globalError.value = 'Billing service is not available. Please try again later.';
  }
};

const handleCancelSubscription = async () => {
  isCancelling.value = true;
  try {
    await subscriptionService.cancelSubscription();
    await fetchSubscription();
  } catch (err: any) {
    subscriptionError.value = err.message || 'Failed to cancel subscription.';
  } finally {
    isCancelling.value = false;
  }
};

onMounted(async () => {
  await fetchUserData();
  await fetchSubscription();

  // If the user has a subscription and is the payer, fetch the invitation code
  if (hasTeamFeatures.value) {
    await fetchInvitationCode();
  }

  if (paddle) {
    paddle.Environment.set('sandbox');
    paddle.Setup({
      token: 'test_bf5c18ea62fd1d30c00bc5c2821',
      debug: true,
    });
  } else {
    console.error('Paddle.js script not found. Make sure it is included in your index.html');
    globalError.value = 'Billing service initialization failed.';
  }
});
</script>

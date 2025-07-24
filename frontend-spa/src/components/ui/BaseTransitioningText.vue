<template>
  <transition enter-active-class="transition-opacity duration-500 ease-out"
    leave-active-class="transition-opacity duration-500 ease-in" enter-from-class="opacity-0"
    leave-to-class="opacity-0">
    <slot></slot>
  </transition>
</template>

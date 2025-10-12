<!-- AI-assisted code -->
<script setup>
import { ref, onMounted } from 'vue'
import { api } from './api.js'

const dark = ref(false)

onMounted(() => {
  dark.value = localStorage.getItem('librelog_dark') === '1'
  applyTheme()
})

function applyTheme() {
  document.documentElement.classList.toggle('dark', dark.value)
}

function toggleDark() {
  dark.value = !dark.value
  localStorage.setItem('librelog_dark', dark.value ? '1' : '0')
  applyTheme()
}
</script>

<template>
  <div class="shell">
    <button class="theme-toggle" @click="toggleDark">{{ dark ? 'light' : 'dark' }}</button>
    <router-view />
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
}
.theme-toggle {
  position: fixed;
  top: 1rem;
  right: 1rem;
  font-size: 0.8rem;
  color: var(--muted);
  z-index: 10;
}
</style>

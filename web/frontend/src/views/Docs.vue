<!-- AI-assisted code -->
<template>
  <div class="docs">
    <div class="docs-head">
      <span class="brand"><img src="/logo.png" alt="" class="brand-logo" />LibreLog</span>
      <router-link to="/app" v-if="hasToken" class="back">app</router-link>
      <router-link to="/" v-else class="back">home</router-link>
    </div>

    <nav class="docs-nav">
      <router-link
        v-for="p in pages"
        :key="p.slug"
        :to="`/docs/${p.slug}`"
        :class="{ active: current?.slug === p.slug }"
      >{{ p.title }}</router-link>
    </nav>

    <div class="docs-content" v-html="rendered"></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../api.js'
import { pages, getPage } from '../docs.js'
import { renderMarkdown } from '../markdown.js'

const route = useRoute()
const hasToken = !!api.getToken()

const current = computed(() => {
  return route.params.page ? getPage(route.params.page) : pages[0]
})

const rendered = computed(() => {
  return current.value ? renderMarkdown(current.value.content) : ''
})
</script>

<style scoped>
.docs {
  max-width: 720px;
  margin: 0 auto;
  padding: 2rem 1.5rem 4rem;
}
.docs-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border);
}
.brand {
  font-family: 'Young Serif', serif;
  font-size: 1.1rem;
  color: var(--accent);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.brand-logo {
  width: 24px;
  height: 24px;
}
.back {
  font-size: 0.85rem;
  color: var(--muted);
}
.docs-nav {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2.5rem;
  font-size: 0.9rem;
}
.docs-nav a {
  color: var(--muted);
}
.docs-nav a.active,
.docs-nav a.router-link-active {
  color: var(--accent);
}

.docs-content :deep(h1) {
  font-size: 2rem;
  color: var(--accent);
  margin-bottom: 1rem;
}
.docs-content :deep(h2) {
  font-size: 1.3rem;
  color: var(--accent);
  margin: 2.5rem 0 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}
.docs-content :deep(h3) {
  font-family: 'Courier New', monospace;
  font-size: 0.95rem;
  font-weight: normal;
  color: var(--text);
  margin: 2rem 0 0.5rem;
}
.docs-content :deep(h4) {
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  font-weight: normal;
  color: var(--text);
  margin: 1.5rem 0 0.5rem;
}
.docs-content :deep(p) {
  color: var(--muted);
  font-size: 0.9rem;
  margin-bottom: 0.75rem;
  line-height: 1.6;
}
.docs-content :deep(ul) {
  color: var(--muted);
  font-size: 0.9rem;
  margin-bottom: 0.75rem;
  padding-left: 1.25rem;
  line-height: 1.8;
}
.docs-content :deep(strong) {
  color: var(--text);
}
.docs-content :deep(code) {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  color: var(--accent);
}
.docs-content :deep(pre) {
  background: var(--surface);
  padding: 0.75rem 1rem;
  border-radius: 3px;
  font-size: 0.8rem;
  font-family: 'Courier New', monospace;
  overflow-x: auto;
  margin-bottom: 0.75rem;
  white-space: pre-wrap;
  word-break: break-all;
}
.docs-content :deep(pre code) {
  color: inherit;
  font-size: inherit;
}
.docs-content :deep(table) {
  margin-bottom: 0.75rem;
  font-size: 0.9rem;
}
.docs-content :deep(th) {
  text-align: left;
  padding: 0.4rem 1.5rem 0.4rem 0;
  color: var(--text);
  font-weight: normal;
}
.docs-content :deep(td) {
  padding: 0.4rem 1.5rem 0.4rem 0;
  color: var(--muted);
}
.docs-content :deep(td:first-child) {
  white-space: nowrap;
}
</style>

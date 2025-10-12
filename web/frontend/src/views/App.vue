<!-- AI-assisted code -->
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api.js'

const router = useRouter()

const view = ref('logsets')
const logsets = ref([])
const selected = ref(null)
const logs = ref([])
const viewMode = ref('table')
const showNewForm = ref(false)
const newForm = ref({ name: '', description: '' })
const editing = ref(false)
const editForm = ref({ name: '', description: '' })
const error = ref('')
const logsLoading = ref(false)
const lastFetchCount = ref(0)
const LIMIT = 100

const apiKeys = ref([])
const showNewKey = ref(false)
const newKeyName = ref('')
const createdKey = ref(null)

const columns = computed(() => {
  const keys = new Set()
  for (const entry of logs.value) {
    try {
      const parsed = JSON.parse(entry.data)
      if (typeof parsed === 'object' && parsed !== null) {
        Object.keys(parsed).forEach(k => keys.add(k))
      }
    } catch {}
  }
  return Array.from(keys)
})

const hasMore = computed(() => lastFetchCount.value === LIMIT)

onMounted(loadLogsets)

watch(selected, async (s) => {
  if (s) {
    logs.value = []
    editing.value = false
    await loadLogs()
  }
})

async function loadLogsets() {
  logsets.value = await api.get('/api/datasets')
}

async function loadLogs() {
  if (!selected.value) return
  logsLoading.value = true
  const data = await api.get(`/api/datasets/${selected.value.log_id}/logs?limit=${LIMIT}`)
  logs.value = data
  lastFetchCount.value = data.length
  logsLoading.value = false
}

async function loadMore() {
  if (!selected.value || logs.value.length === 0) return
  const last = logs.value[logs.value.length - 1].recv_time
  const more = await api.get(`/api/datasets/${selected.value.log_id}/logs?limit=${LIMIT}&before=${encodeURIComponent(last)}`)
  logs.value = [...logs.value, ...more]
  lastFetchCount.value = more.length
}

async function createLogset() {
  error.value = ''
  try {
    const created = await api.post('/api/datasets', newForm.value)
    newForm.value = { name: '', description: '' }
    showNewForm.value = false
    await loadLogsets()
    selected.value = created
  } catch (e) {
    error.value = e.message
  }
}

function startEdit() {
  editForm.value = { name: selected.value.name, description: selected.value.description }
  editing.value = true
}

async function saveEdit() {
  error.value = ''
  try {
    const updated = await api.put(`/api/datasets/${selected.value.log_id}`, editForm.value)
    selected.value = updated
    editing.value = false
    await loadLogsets()
  } catch (e) {
    error.value = e.message
  }
}

async function deleteLogset() {
  if (!confirm(`Delete "${selected.value.name}"?`)) return
  await api.del(`/api/datasets/${selected.value.log_id}`)
  selected.value = null
  logs.value = []
  await loadLogsets()
}

function logout() {
  api.post('/api/logout').catch(() => {})
  api.clearToken()
  router.push('/')
}

function parseData(entry) {
  try { return JSON.parse(entry.data) } catch { return {} }
}

function formatTime(t) {
  return new Date(t).toLocaleString()
}

function showKeys() {
  view.value = 'keys'
  selected.value = null
  loadKeys()
}

async function loadKeys() {
  apiKeys.value = await api.get('/api/tokens')
}

async function createKey() {
  error.value = ''
  try {
    const data = await api.post('/api/tokens', { name: newKeyName.value })
    createdKey.value = data
    newKeyName.value = ''
    showNewKey.value = false
    await loadKeys()
  } catch (e) {
    error.value = e.message
  }
}

async function revokeKey(hash) {
  if (!confirm('Revoke this key?')) return
  await api.del(`/api/tokens/${hash}`)
  await loadKeys()
}

function copyToken(token) {
  navigator.clipboard.writeText(token)
}
</script>

<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-head">
        <span class="brand"><img src="/logo.png" alt="" class="brand-logo" />LibreLog</span>
        <button class="quiet" @click="logout">logout</button>
      </div>

      <div class="logset-list">
        <button
          v-for="ls in logsets"
          :key="ls.log_id"
          class="logset-item"
          :class="{ active: view === 'logsets' && selected?.log_id === ls.log_id }"
          @click="view = 'logsets'; selected = ls"
        >
          {{ ls.name }}
        </button>
      </div>

      <div v-if="showNewForm" class="new-form">
        <input v-model="newForm.name" placeholder="Name" @keyup.enter="createLogset" />
        <input v-model="newForm.description" placeholder="Description (optional)" />
        <div class="new-form-actions">
          <button class="cta" @click="createLogset">Create</button>
          <button @click="showNewForm = false">Cancel</button>
        </div>
      </div>
      <button v-else class="new-logset" @click="showNewForm = true">+ new logset</button>

      <div class="sidebar-footer">
        <button class="nav-item" :class="{ active: view === 'keys' }" @click="showKeys">api keys</button>
      </div>
    </aside>

    <main class="content">
      <div v-if="view === 'keys'" class="keys-panel">
        <h2>API Keys</h2>
        <p class="desc">Keys for external apps and scripts. Each key works like a login token.</p>

        <div v-if="createdKey" class="key-created">
          <p>Key created. Copy it now, you won't see it again.</p>
          <div class="key-display">
            <code>{{ createdKey.token }}</code>
            <button @click="copyToken(createdKey.token)">copy</button>
          </div>
          <button @click="createdKey = null">done</button>
        </div>

        <div class="key-actions">
          <div v-if="showNewKey" class="new-form">
            <input v-model="newKeyName" placeholder="Key name (e.g. vscode-laptop)" @keyup.enter="createKey" />
            <p v-if="error" class="error">{{ error }}</p>
            <div class="new-form-actions">
              <button class="cta" @click="createKey">Create</button>
              <button @click="showNewKey = false">Cancel</button>
            </div>
          </div>
          <button v-else class="new-logset" @click="showNewKey = true">+ new key</button>
        </div>

        <table v-if="apiKeys.length > 0" class="keys-table">
          <thead>
            <tr>
              <th>name</th>
              <th>prefix</th>
              <th>created</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in apiKeys" :key="key.token_hash">
              <td>{{ key.name }}</td>
              <td><code>{{ key.prefix }}...</code></td>
              <td class="time-cell">{{ formatTime(key.created_at) }}</td>
              <td><button class="quiet" @click="revokeKey(key.token_hash)">revoke</button></td>
            </tr>
          </tbody>
        </table>

        <div v-else class="empty-state">
          <p>No API keys yet.</p>
        </div>
      </div>

      <div v-else-if="!selected" class="empty-state">
        <p>Select a logset to view its data, or create a new one.</p>
      </div>

      <div v-else>
        <header class="logset-header">
          <div v-if="!editing">
            <h2>{{ selected.name }}</h2>
            <p class="desc" v-if="selected.description">{{ selected.description }}</p>
            <div class="header-actions">
              <button @click="startEdit">rename</button>
              <button class="quiet" @click="deleteLogset">delete</button>
            </div>
          </div>
          <div v-else class="edit-inline">
            <input v-model="editForm.name" placeholder="Name" />
            <input v-model="editForm.description" placeholder="Description" />
            <p v-if="error" class="error">{{ error }}</p>
            <div class="header-actions">
              <button class="cta" @click="saveEdit">Save</button>
              <button @click="editing = false">cancel</button>
            </div>
          </div>
        </header>

        <div class="log-toolbar">
          <button @click="loadLogs">refresh</button>
          <button @click="viewMode = viewMode === 'table' ? 'json' : 'table'">
            {{ viewMode === 'table' ? 'json' : 'table' }}
          </button>
        </div>

        <div v-if="logsLoading" class="empty-state"><p>Loading...</p></div>

        <div v-else-if="logs.length === 0" class="empty-state">
          <p>No logs yet. Push some data to this logset to see it here.</p>
        </div>

        <table v-else-if="viewMode === 'table'">
          <thead>
            <tr>
              <th>time</th>
              <th v-for="col in columns" :key="col">{{ col }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(entry, i) in logs" :key="i">
              <td class="time-cell">{{ formatTime(entry.recv_time) }}</td>
              <td v-for="col in columns" :key="col">{{ parseData(entry)[col] ?? '' }}</td>
            </tr>
          </tbody>
        </table>

        <div v-else class="json-list">
          <pre v-for="(entry, i) in logs" :key="i" class="json-entry"><span class="time-cell">{{ formatTime(entry.recv_time) }}</span>
{{ JSON.stringify(JSON.parse(entry.data), null, 2) }}</pre>
        </div>

        <div v-if="logs.length > 0 && hasMore" class="load-more">
          <button @click="loadMore">older entries &darr;</button>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 240px;
  border-right: 1px solid var(--border);
  padding: 1rem;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border);
}
.brand {
  font-family: 'Young Serif', serif;
  font-size: 1.1rem;
  color: var(--accent);
  font-weight: normal;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.brand-logo {
  width: 24px;
  height: 24px;
}

.logset-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}
.logset-item {
  text-align: left;
  padding: 0.5rem 0.75rem;
  border-radius: 3px;
  color: var(--text);
  font-size: 0.9rem;
}
.logset-item:hover {
  background: var(--surface);
}
.logset-item.active {
  background: var(--surface);
  color: var(--accent);
}

.new-logset {
  margin-top: 0.5rem;
  text-align: left;
  padding: 0.5rem 0.75rem;
  font-size: 0.85rem;
  color: var(--muted);
}
.new-logset:hover {
  color: var(--accent);
}

.new-form {
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.75rem;
  background: var(--surface);
  border-radius: 3px;
}
.new-form-actions {
  display: flex;
  gap: 0.5rem;
}

.content {
  flex: 1;
  padding: 2rem;
  overflow-x: auto;
}

.empty-state {
  color: var(--muted);
  padding: 4rem 0;
  text-align: center;
  font-style: italic;
}

.logset-header {
  margin-bottom: 2rem;
}
.logset-header h2 {
  font-weight: normal;
  font-size: 1.5rem;
}
.desc {
  color: var(--muted);
  font-size: 0.9rem;
}
.header-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 0.5rem;
}
.edit-inline {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-width: 360px;
}
.error {
  color: var(--error);
  font-size: 0.9rem;
}

.log-toolbar {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border);
}
.log-toolbar button {
  font-size: 0.85rem;
  color: var(--muted);
  padding: 0;
}
.log-toolbar button:hover {
  color: var(--accent);
}

.time-cell {
  white-space: nowrap;
  color: var(--muted);
  font-size: 0.85rem;
}

.json-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.json-entry {
  background: var(--surface);
  padding: 0.75rem 1rem;
  border-radius: 3px;
  font-size: 0.85rem;
  font-family: 'Courier New', monospace;
  overflow-x: auto;
}

.load-more {
  text-align: center;
  padding: 1.5rem 0;
}
.load-more button {
  color: var(--muted);
  font-size: 0.85rem;
}
.load-more button:hover {
  color: var(--accent);
}

.sidebar-footer {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border);
}
.nav-item {
  text-align: left;
  padding: 0.5rem 0.75rem;
  border-radius: 3px;
  color: var(--muted);
  font-size: 0.85rem;
  width: 100%;
}
.nav-item:hover {
  color: var(--accent);
}
.nav-item.active {
  background: var(--surface);
  color: var(--accent);
}

.keys-panel h2 {
  font-weight: normal;
  font-size: 1.5rem;
  margin-bottom: 0.25rem;
}
.keys-panel > .desc {
  color: var(--muted);
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
}
.key-actions {
  margin-bottom: 1.5rem;
}
.key-created {
  background: var(--surface);
  padding: 1rem;
  border-radius: 3px;
  margin-bottom: 1.5rem;
}
.key-created p {
  color: var(--error);
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}
.key-display {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}
.key-display code {
  font-size: 0.8rem;
  font-family: 'Courier New', monospace;
  word-break: break-all;
  flex: 1;
}
.keys-table code {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  color: var(--muted);
}
</style>

<!-- AI-assisted code -->
<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api.js'

const router = useRouter()

if (api.getToken()) router.replace('/app')

const mode = ref('none')
const error = ref('')

const loginForm = ref({ account_number: '', password: '' })
const signupForm = ref({ password: '', name: '' })
const createdAccount = ref('')

async function login() {
  error.value = ''
  try {
    const data = await api.post('/api/login', loginForm.value)
    api.setToken(data.token)
    router.push('/app')
  } catch (e) {
    error.value = e.message
  }
}

async function signup() {
  error.value = ''
  try {
    const data = await api.post('/api/signup', signupForm.value)
    createdAccount.value = data.account_number
    mode.value = 'created'
  } catch (e) {
    error.value = e.message
  }
}

function toLogin() {
  loginForm.value.account_number = createdAccount.value
  createdAccount.value = ''
  mode.value = 'login'
}
</script>

<template>
  <div class="landing">
    <section class="hero">
      <div class="hero-brand">
        <img src="/logo.png" alt="LibreLog" class="hero-logo" />
        <h1>LibreLog</h1>
      </div>
      <p class="tagline">A place to keep track of things.</p>
      <p class="sub">
        Most apps that hold your personal data don't let you really own it.
        LibreLog is a log store that stays out of your way. Send data in,
        pull it out, do whatever you want with it.
        Open source, self-hostable, privacy respecting.
      </p>
      <div class="hero-actions" v-if="mode === 'none'">
        <button class="cta" @click="mode = 'signup'">Sign up</button>
        <button @click="mode = 'login'">Log in</button>
      </div>
    </section>

    <section class="auth" v-if="mode === 'created'">
      <p>Your account number:</p>
      <code class="account-display">{{ createdAccount }}</code>
      <p class="warn">Write this down. It's the only time you'll see it.</p>
      <button class="cta" @click="toLogin">Continue</button>
    </section>

    <section class="auth" v-else-if="mode === 'login'">
      <form @submit.prevent="login">
        <input v-model="loginForm.account_number" placeholder="Account number" required />
        <input v-model="loginForm.password" type="password" placeholder="Password" required />
        <p v-if="error" class="error">{{ error }}</p>
        <button type="submit" class="cta">Log in</button>
        <button type="button" @click="mode = 'signup'">Need an account?</button>
      </form>
    </section>

    <section class="auth" v-else-if="mode === 'signup'">
      <form @submit.prevent="signup">
        <input v-model="signupForm.name" placeholder="A name for your account (optional)" />
        <input v-model="signupForm.password" type="password" placeholder="Password" required />
        <p v-if="error" class="error">{{ error }}</p>
        <button type="submit" class="cta">Create account</button>
        <button type="button" @click="mode = 'login'">Already have one?</button>
      </form>
    </section>

    <section class="about">
      <div class="about-grid">
        <div>
          <h3>Logsets</h3>
          <p>Organize entries into logsets. Exercise, nutrition, coding time, finances, reading, sleep, plant watering.</p>
        </div>
        <div>
          <h3>Send data in</h3>
          <p>Anything that can make an HTTP request can write to LibreLog. A script, a cron job, an app you built.</p>
        </div>
        <div>
          <h3>Pull data out</h3>
          <p>Query your logs through the API. Build charts, run stats, feed it into other software.</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.landing {
  max-width: 640px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.hero {
  padding: 5rem 0 3rem;
}
.hero-brand {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.75rem;
}
.hero-logo {
  width: 56px;
  height: 56px;
}
.hero h1 {
  font-size: 3rem;
  font-weight: normal;
  letter-spacing: -0.02em;
  color: var(--accent);
}
.tagline {
  font-size: 1.25rem;
  margin-bottom: 1rem;
}
.sub {
  color: var(--muted);
  font-size: 0.95rem;
  margin-bottom: 2rem;
  max-width: 480px;
}
.hero-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.auth {
  padding: 2rem 0;
  max-width: 320px;
}
.auth form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.auth form button {
  align-self: flex-start;
}
.error {
  color: var(--error);
  font-size: 0.9rem;
}
.account-display {
  display: block;
  font-size: 2rem;
  font-family: 'Georgia', serif;
  padding: 1rem 0;
  letter-spacing: 0.08em;
  color: var(--accent);
}
.warn {
  color: var(--error);
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
}

.about {
  padding: 4rem 0 2rem;
  border-top: 1px solid var(--border);
  margin-top: 3rem;
}
.about-grid {
  display: grid;
  gap: 2rem;
}
.about h3 {
  font-size: 1rem;
  font-weight: normal;
  color: var(--accent);
  margin-bottom: 0.25rem;
}
.about p {
  color: var(--muted);
  font-size: 0.9rem;
}
</style>

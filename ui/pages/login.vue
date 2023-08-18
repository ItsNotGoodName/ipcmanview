<script lang="ts" setup>
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  layout: "center",
  middleware: [
    'auth-restrict',
  ]
});

const authStore = useAuthStore()

const usernameOrEmail = ref("")
const password = ref("")

const {
  mutate: login,
  loading: loginLoading,
  error: loginError
} = useMutation(() => authStore.login(usernameOrEmail.value, password.value))
</script>

<template>
  <div class="w-full max-w-md pt-16 px-2">
    <UiCard>
      <template #header>
        <div class="flex justify-between">
          <div class="text-lg overflow-hidden text-ellipsis">
            IPCManView
          </div>
          <ThemeSwitcher />
        </div>
      </template>
      <form class="flex flex-col p-4 gap-4" @submit.prevent="login">
        <UiInput v-model="usernameOrEmail" :disabled="loginLoading" label="Username or Email"
          placeholder="Username or Email" autocomplete="username" />
        <UiInput v-model="password" :disabled="loginLoading" label="Password" placeholder="Password" type="password"
          autocomplete="current-password" />
        <UiButton :disabled="loginLoading" type="submit">Log in</UiButton>
        <p v-if="loginError" class="text-Red">{{ loginError.cause || loginError.message }}</p>
      </form>
    </UiCard>
    <div class="flex justify-center mt-2">
      <NuxtLink class="text-Blue" to="/register">Register</NuxtLink>
    </div>
  </div>
</template>

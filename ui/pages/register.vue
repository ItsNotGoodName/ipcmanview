<script lang="ts" setup>
import { UserRegister } from '~/core/client.gen';

definePageMeta({
  layout: "center",
  middleware: [
    'auth-restrict',
  ]
});

const { $authService } = useNuxtApp()

const req = reactive<UserRegister>({
  email: "",
  username: "",
  password: "",
  passwordConfirm: ""
})

const {
  mutate: register,
  loading: registerLoading,
  error: registerError
} = useMutation(() => $authService.register({ user: req }).
  then(() => { navigateTo('/') }))
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
      <form class="flex flex-col p-4 gap-4" @submit.prevent="register">
        <UiInput :disabled="registerLoading" v-model="req.email" label="Email" placeholder="Email" type="email" />
        <UiInput :disabled="registerLoading" v-model="req.username" label="Username" placeholder="Username" />
        <UiInput :disabled="registerLoading" v-model="req.password" label="Password" placeholder="Password"
          type="password" />
        <UiInput :disabled="registerLoading" v-model="req.passwordConfirm" label="Password confirm"
          placeholder="Password confirm" type="password" />
        <UiButton :disabled="registerLoading" type="submit">Register</UiButton>
        <p v-if="registerError" class="text-Red">{{ registerError.cause || registerError.message }}</p>
      </form>
    </UiCard>
    <div class="flex justify-center mt-2">
      <NuxtLink class="text-Blue" to="/login">Login</NuxtLink>
    </div>
  </div>
</template>

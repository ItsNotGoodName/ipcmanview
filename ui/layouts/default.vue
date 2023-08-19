<script setup lang="ts">
import { useAuthStore } from '~/stores/auth';
import { formatInitials } from '~/utils'

const authStore = useAuthStore();

const { mutate: logout, loading: logoutLoading } = useMutation(() => authStore.logout())
</script>

<template>
  <NuxtLoadingIndicator class="from-Mauve to-Mauve-200 bg-gradient-to-r" color="" />
  <div class="bg-Base text-Text fixed inset-0">
    <div class="flex h-full flex-col">
      <!-- Header -->
      <div>
        <div class="bg-Crust border-b-Overlay flex h-11 justify-between gap-2 border-b p-2">
          <div class="flex items-center overflow-hidden text-2xl">
            <div class="truncate">
              IPCManView
            </div>
          </div>
          <div class="flex gap-1">
            <!-- Theme switcher button -->
            <ThemeSwitcher />
            <!-- Profile dropdown -->
            <HeadlessMenu as="div">
              <HeadlessMenuButton title="Profile"
                class="hover:text-Mauve hover:fill-Mauve flex items-center justify-center">
                <div class="bg-Surface flex h-7 w-7 items-center justify-center rounded-full">
                  {{ formatInitials(authStore.user.username) }}
                </div>
              </HeadlessMenuButton>
              <div class="flex flex-row-reverse">
                <transition enter-active-class="transition duration-100 ease-out"
                  enter-from-class="transform scale-95 opacity-0" enter-to-class="transform scale-100 opacity-100"
                  leave-active-class="transition duration-75 ease-out" leave-from-class="transform scale-100 opacity-100"
                  leave-to-class="transform scale-95 opacity-0">
                  <HeadlessMenuItems
                    class="bg-Surface border-Overlay right absolute right-0 z-10 mr-2 mt-2 flex w-32 origin-top-right flex-col rounded p-1 shadow-lg">
                    <HeadlessMenuItem v-slot="{ active, close }">
                      <NuxtLink class="rounded" :class='{ "bg-Mauve text-Crust": active }' to="/profile">
                        <button class="flex w-full items-center gap-2 p-1 text-left" @click="close">
                          <Icon name="ri:user-line" class="h-5 w-5" />
                          Profile
                        </button>
                      </NuxtLink>
                    </HeadlessMenuItem>
                    <HeadlessMenuItem v-slot="{ active }">
                      <div class="rounded" :class='{ "bg-Red text-Crust": active }'>
                        <button class="flex w-full items-center gap-2 p-1" @click="logout" :disabled="logoutLoading">
                          <Icon name="ri:logout-circle-r-line" class="h-5 w-5" />
                          Log out
                        </button>
                      </div>
                    </HeadlessMenuItem>
                  </HeadlessMenuItems>
                </transition>
              </div>
            </HeadlessMenu>
          </div>
        </div>
      </div>
      <div class="flex h-full flex-col overflow-hidden md:flex-row">
        <!-- Nav -->
        <div>
          <div
            class="bg-Mantle border-b-Overlay md:border-r-Overlay flex h-11 gap-1 overflow-x-auto border-b p-2 md:h-full md:w-11 md:flex-col md:border-b-0 md:border-r ">
            <!-- Home link -->
            <NuxtLink title="Home" to="/" active-class="text-Mauve border-Mauve fill-Mauve"
              class="hover:text-Mauve hover:fill-Mauve flex items-center justify-center">
              <Icon name="ri:home-5-line" class="h-7 w-7" />
            </NuxtLink>
            <!-- Palette link -->
            <NuxtLink title="Palette" to="/palette" active-class="text-Mauve border-Mauve fill-Mauve"
              class="hover:text-Mauve hover:fill-Mauve flex items-center justify-center">
              <Icon name="ri:paint-brush-line" class="h-7 w-7" />
            </NuxtLink>
          </div>
        </div>
        <!-- Content -->
        <div class="h-full w-full overflow-auto">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

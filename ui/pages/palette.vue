<script lang="ts" setup>
definePageMeta({
  middleware: [
    'auth'
  ]
})

// Card
const showCard = ref(true)

// Modal
const modalOpen = ref(false)
const setModelOpen = (open: boolean) => {
  modalOpen.value = open
}
</script>

<template>
  <div class="flex flex-col gap-2 p-2">
    <UiStatusBanner status="idle">
      Idle
    </UiStatusBanner>
    <UiStatusBanner status="pending">
      Pending
    </UiStatusBanner>
    <UiStatusBanner status="success">
      Success
    </UiStatusBanner>
    <UiStatusBanner status="error">
      Error
    </UiStatusBanner>
    <hr class="text-Overlay">
    <!-- Buttons -->
    <div class="flex flex-wrap gap-2">
      <UiButton color="primary">Primary</UiButton>
      <UiButton color="secondary">Secondary</UiButton>
      <UiButton color="error">Error</UiButton>
      <UiButton color="success">Success</UiButton>
      <UiButton color="info">Info</UiButton>
      <UiButton color="warning">Warning</UiButton>
    </div>
    <div class="flex flex-wrap gap-2">
      <UiButton color="primary" size="sm">Primary</UiButton>
      <UiButton color="secondary" size="sm">Secondary</UiButton>
      <UiButton color="error" size="sm">Error</UiButton>
      <UiButton color="success" size="sm">Success</UiButton>
      <UiButton color="info" size="sm">Info</UiButton>
      <UiButton color="warning" size="sm">Warning</UiButton>
    </div>
    <hr class="text-Overlay">
    <!-- Card -->
    <UiCard v-if="showCard">
      <template #header>
        <div class="flex flex-row items-center justify-between">
          <div>Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.
          </div>
          <button class="hover:bg-Surface-200 text-Red flex items-center rounded-full" @click="showCard = false">
            <Icon class="h-5 w-5" name="ri:close-line" />
          </button>
        </div>
      </template>
      <div class="p-2">
        Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem
        pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud
        ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut
        officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris
        cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia
        eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat
        ullamco ut ea consectetur et est culpa et culpa duis.
      </div>
    </UiCard>
    <hr class="text-Overlay">
    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="border-Overlay bg-Surface w-full border">
        <caption class="mb-2 caption-top">
          Cameras
        </caption>
        <thead class="bg-Surface-100">
          <tr>
            <th class="border-Overlay w-0 border p-1">
              <div class="flex items-center">
                <input class="accent-Mauve bg-Crust text-Base h-5 w-5" type="checkbox">
              </div>
            </th>
            <th class="border-Overlay border-y p-1">
              <button class="flex w-full items-center justify-center gap-1">
                <div>
                  Address
                </div>
                <Icon class="h-5 w-5" name="ri:sort-asc" />
              </button>
            </th>
            <th class="border-Overlay border-y p-1">
              <button class="flex w-full items-center justify-center gap-1">
                <div>
                  Username
                </div>
                <Icon class="h-5 w-5" name="ri:sort-desc" />
              </button>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="border-Overlay w-0 border p-1">
              <div class="flex items-center">
                <input class="accent-Mauve bg-Crust text-Base h-5 w-5" type="checkbox">
              </div>
            </td>
            <td class="border-Overlay border-y p-1">192.168.1.108</td>
            <td class="border-Overlay border-y p-1">admin</td>
          </tr>
          <tr>
            <td class="border-Overlay w-0 border p-1">
              <div class="flex items-center">
                <input class="accent-Mauve bg-Crust text-Base h-5 w-5" type="checkbox">
              </div>
            </td>
            <td class="border-Overlay border-y p-1">192.168.1.108</td>
            <td class="border-Overlay border-y p-1">
              <a class="text-Blue hover:text-Blue-200 hover:underline" href="#">admin</a>
            </td>
          </tr>
          <tr>
            <td class="border-Overlay w-0 border p-1">
              <div class="flex items-center">
                <input class="accent-Mauve bg-Crust text-Base h-5 w-5" type="checkbox">
              </div>
            </td>
            <td class="border-Overlay border-y p-1">192.168.1.108</td>
            <td class="border-Overlay border-y">
              <UiButton class="w-full" size="sm" color="warning">Button</UiButton>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <hr class="text-Overlay">
    <!-- Modal -->
    <UiButton @click="setModelOpen(true)">
      Open Modal
    </UiButton>
    <ClientOnly>
      <HeadlessTransitionRoot appear :show="modalOpen" as="template">
        <HeadlessDialog :open="modalOpen" @close="setModelOpen" class="relative z-10">
          <HeadlessTransitionChild as="template" enter="duration-200 ease-out" enter-from="opacity-0"
            enter-to="opacity-100">
            <div class="bg-Mantle fixed inset-0 bg-opacity-25" />
          </HeadlessTransitionChild>
          <div class="fixed inset-0 overflow-y-auto">
            <HeadlessTransitionChild as="template" enter="duration-200 ease-out" enter-from="opacity-0 scale-95"
              enter-to="opacity-100 scale-100">
              <HeadlessDialogPanel class="text-Text mx-auto mt-32 flex max-w-md">
                <UiCard class="w-full shadow-lg">
                  <template #header>
                    Modal
                  </template>
                  <div class="flex flex-col gap-2 p-2">
                    <p>
                      Are you sure you want to deactivate your account? All of your data will be
                      permanently removed. This action cannot be undone.
                    </p>
                    <div class="flex gap-2">
                      <UiButton color="error" @click="setModelOpen(false)">Deactivate</UiButton>
                      <UiButton color="secondary" @click="setModelOpen(false)">Cancel</UiButton>
                    </div>
                  </div>
                </UiCard>
              </HeadlessDialogPanel>
            </HeadlessTransitionChild>
          </div>
        </HeadlessDialog>
      </HeadlessTransitionRoot>
    </ClientOnly>
    <hr class="text-Overlay">
    <!-- Property -->
    <UiCard>
      <template #header>
        Detail
      </template>
      <div class="divide-Overlay divide-y">
        <div class="flex justify-between gap-2 p-2">
          <div>Vendor</div>
          <div>Lorem</div>
        </div>
        <div class="flex justify-between gap-2 p-2">
          <div>SN</div>
          <div>XXXX-XXX-XXXX</div>
        </div>
      </div>
    </UiCard>
    <!-- Input -->
    <UiInput placeholder="Error"
      message="Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."
      color="error" label="Title" />
    <UiInput placeholder="Success"
      message="Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."
      label="Title" />
    <UiInput placeholder="Normal"
      message="Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."
      color="success" label="Title" />
  </div>
</template>

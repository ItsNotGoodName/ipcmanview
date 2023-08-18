<script lang="ts" setup>
definePageMeta({
  middleware: [
    'auth'
  ]
})

const { $dahuaService } = useNuxtApp()

const cameraCount = useAsyncData(() => $dahuaService.cameraCount().then(res => res.count))
const activeScanCount = useAsyncData(() => $dahuaService.activeScannerCount().then(res => "" + res.activeScanners + "/" + res.totalScanners))
</script>

<template>
  <div class="flex flex-col p-4 sm:flex-row gap-4">
    <DataStatCard :data="cameraCount" label="Cameras" icon="bx:bxs-cctv" />
    <DataStatCard :data="activeScanCount" label="Scanners" icon="tabler:zoom-filled" />
  </div>
</template>

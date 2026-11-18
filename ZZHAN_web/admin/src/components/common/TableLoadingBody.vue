<script setup lang="ts">
defineProps<{ columns: number; label: string }>()
</script>

<template>
  <tbody>
    <tr>
      <td :colspan="columns" class="table-loading-cell"><span role="status">{{ label }}</span></td>
    </tr>
  </tbody>
</template>

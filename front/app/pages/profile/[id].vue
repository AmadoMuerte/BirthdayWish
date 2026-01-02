<script setup lang="ts">
const route = useRoute()
const id = route.params.id as string
</script>

<template>
 <div>{{ id }}</div>
</template>
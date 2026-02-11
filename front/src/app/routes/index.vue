<script setup lang="ts">


</script>

<template>
    <div class="container">
        <h1>Home</h1>
    </div>
</template>

<style scoped lang="scss">
.container {
    max-width: var(--content);
    margin: 0 auto;
    padding: var(--space-4);
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: var(--space-6);
}
</style>
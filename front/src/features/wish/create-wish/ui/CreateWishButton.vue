<script setup lang="ts">
import { IconPlus } from '@tabler/icons-vue';
import { Text } from '@/shared/ui/Typography';

</script>

<template>
    <button class="create-wish-button">
        <div class="create-wish-button__icon">
            <IconPlus size="18" />
        </div>
            <Text size="sm" weight="medium">Create Wish</Text>
    </button>
</template>

<style lang="scss" scoped>
.create-wish-button {
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    border-radius: var(--radius-xl);
    border: 1px dashed var(--border);
    background-color: var(--card);
    transition: var(--transition-normal);
    color: var(--muted);
}

.create-wish-button__icon {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--space-2);
    border-radius: var(--radius-full);
    border: 1px dashed var(--border);
}

.create-wish-button:hover {
    color: var(--foreground);
}
</style>


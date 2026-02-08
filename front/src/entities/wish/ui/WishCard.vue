<script setup lang="ts">
import { Heading } from '@/shared/ui/Typography';
import { Text } from '@/shared/ui/Typography';
import { Button } from '@/shared/ui/Button';
import type { Wish } from '../model/wish.types';
import { IconHeart, IconLink, IconX } from '@tabler/icons-vue';

interface Props {
    wish: Wish
}

const { wish } = defineProps<Props>()
</script>

<template>
    <div class="wish-card">
        <img class="wish-card__image" :src="wish.image" alt="">
        <div class="wish-card__content">
            <div>
                <Text class="wish-card__marketplace" size="sm" weight="medium">{{ wish.marketplace }}</Text>
                <Heading class="wish-card__title" size="md">{{ wish.title }}</Heading>
            </div>

            <div class="wish-card__price-container">
                <Text class="wish-card__price" size="md" weight="medium">{{ wish.price }} {{ wish.currency }}</Text>
                <Button size="xs" radius="full" variant="outline">
                    <IconLink size="18" />Link
                </Button>
            </div>
        </div>
        <div class="wish-card__actions">
            <slot name="actions" />
            <Button size="sm" radius="full" variant="outline" :square="true">
                <IconX size="18" />
            </Button>
            <Button size="sm" radius="full" variant="outline" :square="true">
                <IconHeart size="18" />
            </Button>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.wish-card {
    position: relative;
    display: flex;
    flex-direction: column;
    padding: var(--space-2);
    gap: var(--space-1);
    border: 1px solid var(--border);
    border-radius: var(--radius-xl);
    background-color: var(--card);
    box-shadow: var(--shadow-sm);
    transition: var(--transition-normal);
}

.wish-card:hover {
    box-shadow: var(--shadow-lg);
}

.wish-card__content {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-2);
}

.wish-card__marketplace {
    color: var(--muted);
}

.wish-card__image {
    width: 100%;
    aspect-ratio: 1;
    object-fit: cover;
    border-radius: var(--radius-lg);
}

.wish-card__price-container {
    display: flex;
    gap: var(--space-2);
}

.wish-card__price {
    color: var(--muted);
}

.wish-card__actions {
    position: absolute;
    top: 0;
    right: 0;
    left: 0;
    z-index: 1;
    padding: var(--space-4);

    display: flex;
    gap: var(--space-2);
}
</style>
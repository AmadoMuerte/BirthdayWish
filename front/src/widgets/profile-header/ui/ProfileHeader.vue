<script setup lang="ts">
import { Heading, Text } from '@/shared/ui/Typography';

const { user } = defineProps<{
    user: User,
}>()
</script>

<template>
    <div class="profile-header">
        <img class="profile-header__banner" :src="user.banner" :alt="user.name + ' banner'">
        <div class="profile-header__content">
            <div class="profile-header__avatar">
                <img :src="user.avatar" :alt="user.name + ' avatar'">
            </div>
            <div class="profile-header__info">
                <Heading size="xl">{{ user.name }}</Heading>
                <Text size="sm" weight="medium" color="muted">@{{ user.username }}</Text>
            </div>
            <div class="profile-header__actions">
                <slot name="actions" />
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.profile-header {
    display: flex;
    flex-direction: column;
    border-radius: var(--radius-xl);
    overflow: hidden;
    padding: var(--space-2);
    padding-bottom: var(--space-0);
    background-color: var(--card);
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--border);
}

.profile-header__banner {
    width: 100%;
    height: 240px;
    object-fit: cover;
    border-radius: var(--radius-md);
}

.profile-header__content {
    position: relative;
    display: flex;
    gap: var(--space-6);
    padding-left: var(--space-6);
    padding-right: var(--space-4);
}

.profile-header__avatar {
    position: relative;
    width: 200px;
    margin: var(--space-4) 0;

    img {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        aspect-ratio: 1;
        object-fit: cover;
        border-radius: 50%;
        border: 4px solid var(--card);
        box-shadow: var(--shadow-md);
    }
}

.profile-header__info {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: var(--space-4) 0;
}

.profile-header__actions {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-6) 0;
}
</style>

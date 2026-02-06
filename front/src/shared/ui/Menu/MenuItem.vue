<script setup lang="ts">
import {
  DropdownMenuItem,
} from 'reka-ui'

</script>

<template>
    <DropdownMenuItem class="menu-item">
        <div class="menu-item__left">
            <slot name="left" />
        </div>
        <div class="menu-item__content">
            <slot />
        </div>
        <div class="menu-item__right">
            <slot name="right" />
        </div>
    </DropdownMenuItem>
</template>

<style lang="scss" scoped>
.menu-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    height: 38px;
    min-width: 200px;
    background-color: var(--card);
    padding: 0 var(--space-3);
    border-radius: var(--radius-sm);
    font-weight: var(--font-medium);
    font-size: var(--text-sm);

    &__left {
        display: flex;
        align-items: center;
    }

    &__content {
        display: flex;
        align-items: center;
    }

    &__right {
        display: flex;
        align-items: center;
    }

    &:hover {
        background-color: var(--outline-hover);
    }
}
</style>
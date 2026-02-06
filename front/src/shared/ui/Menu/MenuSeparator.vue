<script setup lang="ts">
import {
  DropdownMenuSeparator,
} from 'reka-ui'

</script>

<template>
    <DropdownMenuSeparator class="menu-separator" />
</template>

<style lang="scss" scoped>
.menu-separator {
    height: 1px;
    margin: 0 var(--space-1);
    background-color: var(--border);
    border-radius: var(--radius-full);
}
</style>
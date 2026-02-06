<script setup lang="ts">
import {
  DropdownMenuArrow,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from 'reka-ui'
</script>

<template>
  <DropdownMenuRoot class="menu-root">
    <DropdownMenuTrigger asChild>
        <slot name="trigger" />
    </DropdownMenuTrigger>
    <DropdownMenuPortal>
      <DropdownMenuContent class="menu-content" align="end" :sideOffset="8" :arrowPadding="10">
        <slot name="content" />
      </DropdownMenuContent>
    </DropdownMenuPortal>
  </DropdownMenuRoot>
</template>

<style lang="scss" scoped>
.menu-root {
    display: flex;
    align-items: center;
}

:deep(.menu-content) {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-1);
    border-radius: var(--radius-md);
    background-color: var(--card);
    box-shadow: var(--shadow-md);
    border: 1px solid var(--border);
}
</style>

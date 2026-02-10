<script setup lang="ts">
import Group from '@/shared/ui/Group.vue';
import { UserMenu } from '@/features/user/user-menu';
import { Button } from '@/shared/ui/Button';
import { IconPlus } from '@tabler/icons-vue';
</script>

<template>
    <header class="header">
        <Group class="header__content" justify="between" align="center">
            <div class="header__title">Wish</div>
            <div class="header__actions">
                <Button size="md" radius="sm" variant="outline">
                    
                        <IconPlus size="20" />
                  
                    Add wish
                </Button>
                <UserMenu />
            </div>
        </Group>
    </header>
</template>

<style lang="scss" scoped>
.header {
    width: 100%;
    background: var(--background);
    margin-bottom: var(--space-8);

    &__content {
        max-width: var(--content);
        padding: 0 var(--space-4);
        margin: 0 auto;
        height: 64px;
    }

    &__title {
        font-size: 32px;
        font-weight: bold;
    }

    &__actions {
        display: flex;
        align-items: center;
        gap: var(--space-2);
    }
}
</style>

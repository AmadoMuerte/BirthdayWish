<script setup lang="ts">
import { Heading, Text } from '@/shared/ui/Typography';
import { Menu, MenuItem, MenuSeparator } from '@/shared/ui/Menu';
import { IconUser, IconSettings, IconLogout, IconMoon } from '@tabler/icons-vue';

const user = {
    name: 'Mazde',
    username: 'mazde',
    avatar: 'https://images.unsplash.com/photo-1574158622682-e40e69881006?q=80&w=1160&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
}
</script>

<template>
    <Menu>
        <template #trigger>
            <button class="user-menu">
                <img class="user-menu__avatar" :src="user.avatar" :alt="user.name + ' avatar'" />
                <Heading size="md" class="user-menu__name">{{ user.name }}</Heading>
            </button>
        </template>
        <template #content>
            <MenuItem>
                <template #left>
                    <IconUser size="20" />
                </template>
                Profile
            </MenuItem>
            <MenuItem>
                <template #left>
                    <IconSettings size="20" />
                </template>
                Settings
            </MenuItem>
            <MenuItem>
                <template #left>
                    <IconMoon size="20" />
                </template>
                Theme
            </MenuItem>
            <MenuSeparator />
            <MenuItem class="text-danger">
                <template #left>
                    <IconLogout size="20" />
                </template>
                Logout
            </MenuItem>
        </template>
    </Menu>
</template>

<style lang="scss" scoped>
.user-menu {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    border: none;
    background: none;
    cursor: pointer;
}

.user-menu__info {
    display: flex;
    flex-direction: column;
}

.user-menu__avatar {
    width: 42px;
    height: 42px;
    border-radius: 50%;
    object-fit: cover;
}

.text-danger {
    color: var(--danger-500);
}
</style>
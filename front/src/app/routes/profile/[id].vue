<script setup lang="ts">
import { WishCard } from '@/entities/wish';
import { CreateWishButton } from '@/features/wish/create-wish';
import { ProfileHeader } from '@/widgets/profile-header';
import { IconHeart, IconShare } from '@tabler/icons-vue';
import { Button } from '@/shared/ui/Button';

const wish = {
    id: '1',
    title: 'iPhone 15 Pro Max',
    marketplace: 'amazon.com',
    image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?q=80&w=1160&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
    price: 100,
    currency: 'RUB',
}

const user = {
    id: '1',
    name: 'Mazde',
    username: 'mazde',
    avatar: 'https://images.unsplash.com/photo-1574158622682-e40e69881006?q=80&w=1160&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
    banner: 'https://images.unsplash.com/photo-1603097122177-5d9820645ac5?q=80&w=2574&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
}
</script>

<template>
    <div class="profile">
        <ProfileHeader :user="user">
            <template #actions>
                <Button size="md" radius="sm">
                    <IconHeart size="24" />Follow
                </Button>
                <Button size="md" radius="sm" variant="outline">
                    <IconShare size="24" />Share
                </Button>
            </template>
        </ProfileHeader>
        <div class="container">
            <CreateWishButton />
            <WishCard :wish="wish" />
            <WishCard :wish="wish" />
            <WishCard :wish="wish" />
            <WishCard :wish="wish" />
            <WishCard :wish="wish" />
            <WishCard :wish="wish" />
        </div>
    </div>
</template>

<style scoped lang="scss">
.profile {
    max-width: var(--content);
    margin: 0 auto;
    padding: 0 var(--space-4);

    display: flex;
    flex-direction: column;
    gap: var(--space-12);
}

.container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: var(--space-6);
}
</style>
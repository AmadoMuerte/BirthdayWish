<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority'
import { Primitive, type PrimitiveProps } from 'reka-ui'

const avatarClasses = cva('avatar', {
    variants: {
        size: {
            xs: 'size-xs',
            sm: 'size-sm',
            md: 'size-md',
            lg: 'size-lg',
            xl: 'size-xl',
        },
        variant: {
            image: 'variant-image',
            fallback: 'variant-fallback',
        },
        shape: {
            circle: 'shape-circle',
            square: 'shape-square',
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'image',
        shape: 'circle',
    },
})

type AvatarClasses = VariantProps<typeof avatarClasses>

interface Props extends PrimitiveProps {
    size?: AvatarClasses['size']
    variant?: AvatarClasses['variant']
    shape?: AvatarClasses['shape']
    src?: string
    alt?: string
    fallback?: string
    as?: string
}

const props = withDefaults(defineProps<Props>(), {
    size: 'md',
    variant: 'image',
    shape: 'circle',
    as: 'div',
})
</script>

<template>
    <Primitive 
        :as="props.as"
        :class="avatarClasses({ 
            size: props.size, 
            variant: props.variant, 
            shape: props.shape 
        })"
        :data-src="props.src"
        v-bind="$attrs"
    >
        <!-- Image fallback -->
        <img 
            v-if="props.src && props.variant === 'image'"
            :src="props.src" 
            :alt="props.alt || 'Avatar'"
            class="avatar-image"
        >
        
        <!-- Fallback content -->
        <slot v-else name="fallback" :default="props.fallback">
            <span class="avatar-fallback">{{ props.fallback || 'JD' }}</span>
        </slot>
    </Primitive>
</template>

<style scoped lang="scss">
.avatar {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
    flex-shrink: 0;
    user-select: none;

    /* SIZES */
    &.size-xs {
        width: var(--size-6); height: var(--size-6);
        font-size: var(--text-xs);
    }
    &.size-sm {
        width: var(--size-8); height: var(--size-8);
        font-size: var(--text-sm);
    }
    &.size-md {
        width: var(--size-10); height: var(--size-10);
        font-size: var(--text-sm);
    }
    &.size-lg {
        width: var(--size-12); height: var(--size-12);
        font-size: var(--text-md);
    }
    &.size-xl {
        width: var(--size-14); height: var(--size-14);
        font-size: var(--text-lg);
    }

    /* SHAPE */
    &.shape-circle {
        border-radius: var(--radius-full);
    }
    &.shape-square {
        border-radius: var(--radius-md);
    }

    /* VARIANTS */
    &.variant-image {
        background: var(--background-muted);
        border: 2px solid var(--border);
    }
    &.variant-fallback {
        background: var(--primary);
        color: var(--primary-foreground);
        font-weight: var(--font-semibold);
    }

    .avatar-image {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .avatar-fallback {
        line-height: 1;
        font-weight: var(--font-semibold);
    }

    &:hover:not([disabled]) {
        transform: scale(1.05);
        transition: transform 0.2s ease;
    }
}
</style>


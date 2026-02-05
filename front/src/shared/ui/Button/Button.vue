<script setup lang="ts">
import type { ButtonHTMLAttributes } from 'vue'
import { cva, type VariantProps } from 'class-variance-authority'
import { Primitive, type PrimitiveProps } from 'reka-ui'

const buttonClasses = cva('button', {
    variants: {
        variant: {
            default: 'variant-default',
            outline: 'variant-outline',
        },
        size: {
            xs: 'size-xs',
            sm: 'size-sm',
            md: 'size-md',
            lg: 'size-lg'
        },
        radius: {
            sm: 'radius-sm',
            md: 'radius-md',
            lg: 'radius-lg',
            full: 'radius-full',
        },
        square: {
            true: 'is-square'
        }
    },
})

type ButtonClasses = VariantProps<typeof buttonClasses>

interface Props extends /* @vue-ignore */ ButtonHTMLAttributes, /* @vue-ignore */ PrimitiveProps {
    variant?: ButtonClasses['variant']
    size?: ButtonClasses['size']
    radius?: ButtonClasses['radius']
    square?: ButtonClasses['square']
    as?: string
}

defineOptions({
    name: 'Button',
})

const props = withDefaults(defineProps<Props>(), {
    variant: 'default',
    size: 'md',
    radius: 'md',
    as: 'button',
    square: false,
})
</script>

<template>
    <Primitive :as="props.as" :class="buttonClasses({ variant, size, radius, square })">
        <slot name="prefix" />
        <slot />
        <slot name="suffix" />
    </Primitive>
</template>

<style scoped lang="scss">
.button {
    font-family: var(--font-sans);
    font-weight: var(--font-medium);
    cursor: pointer;
    border: none;
    outline: none;
    display: flex;
    align-items: center;
    gap: var(--space-2);

    &:focus-visible:not([disabled]) {
        outline-offset: .125rem;
        outline: var(--focus);
    }

    /** variants */
    &.variant {
        &-default {
            background: var(--primary);
            color: var(--primary-foreground);

            &:hover:not([disabled]) {
                background: var(--primary-hover);
            }
        }

        &-outline {
            background: var(--card);
            border: 1px solid var(--border);

            &:hover:not([disabled]) {
                background: var(--outline-hover);
            }
        }
    }

    /** sizes */
    &.size {
        &-xs {
            font-size: var(--text-xs);
            height: var(--size-xs);
            padding: 0 var(--space-2);
        }
        
        &-sm {
            font-size: var(--text-sm);
            height: var(--size-sm);
            padding: 0 var(--space-4);
        }

        &-md {
            font-size: var(--text-md);
            height: var(--size-md);
            padding: 0 var(--space-6);
        }

        &-lg {
            font-size: var(--text-lg);
            height: var(--size-lg);
            padding: 0 var(--space-7);
        }
    }

    /* radius */
    &.radius {
        &-sm {
            border-radius: var(--radius-sm);
        }

        &-md {
            border-radius: var(--radius-md);
        }

        &-lg {
            border-radius: var(--radius-lg);
        }

        &-full {
            border-radius: var(--radius-full);
        }
    }

    &.is-square {
        aspect-ratio: 1;
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    /* attributes */
    &[disabled] {
        cursor: not-allowed;
        user-select: none;
        opacity: var(--opacity-50);
    }
}
</style>
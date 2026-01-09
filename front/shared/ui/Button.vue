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
            sm: 'size-sm',
            md: 'size-md',
            lg: 'size-lg'
        },
        radius: {
            sm: 'radius-sm',
            md: 'radius-md',
            lg: 'radius-lg',
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
        <slot />
    </Primitive>
</template>

<style scoped lang="scss">
.button {
    font-family: var(--font-sans);
    font-weight: var(--font-medium);
    cursor: pointer;
    border: none;
    outline: none;

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
            background: var(--background);
            border: 1px solid var(--border);

            &:hover:not([disabled]) {
                background: var(--outline-hover);
            }
        }
    }

    /** sizes */
    &.size {
        &-sm {
            font-size: var(--text-sm);
            height: var(--size-sm);
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
    }

    &.is-square {
        aspect-ratio: 1;
        padding: 0;
    }

    /* attributes */
    &[disabled] {
        cursor: not-allowed;
        user-select: none;
        opacity: var(--opacity-50);
    }
}
</style>
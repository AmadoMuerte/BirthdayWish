<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority';
const textClasses = cva('text', {
    variants: {
        size: {
            sm: 'size-sm',
            md: 'size-md',
            lg: 'size-lg',
        },
        weight: {
            normal: 'weight-normal',
            medium: 'weight-medium',
            semibold: 'weight-semibold',
            bold: 'weight-bold',
        },
        color: {
            default: 'color-default',
            muted: 'color-muted',
        },
    },
})

type TextClasses = VariantProps<typeof textClasses>

interface Props {
    tag?: 'p' | 'span' | 'div',
    size?: TextClasses['size'],
    weight?: TextClasses['weight'],
    color?: TextClasses['color'],
}

const props = withDefaults(defineProps<Props>(), {
    size: 'md',
    tag: 'p',
    color: 'default',
})
</script>

<template>
    <component :is="props.tag" :class="textClasses({ size, weight, color })">
        <slot />
    </component>
</template>

<style lang="scss" scoped>
.text {
    &.size {
        &-sm {
            font-size: var(--text-sm);
        }

        &-md {
            font-size: var(--text-md);
        }

        &-lg {
            font-size: var(--text-lg);
        }
    }

    &.weight {
        &-normal {
            font-weight: var(--font-normal);
        }

        &-medium {
            font-weight: var(--font-medium);
        }

        &-semibold {
            font-weight: var(--font-semibold);
        }

        &-bold {
            font-weight: var(--font-bold);
        }
    }

    &.color {
        &-default {
            color: var(--text);
        }

        &-muted {
            color: var(--muted);
        }
    }
}
</style>
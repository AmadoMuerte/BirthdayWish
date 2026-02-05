<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority';

const headingClasses = cva('heading', {
    variants: {
        size: {
            sm: 'size-sm',
            md: 'size-md',
            lg: 'size-lg',
            xl: 'size-xl',
        },
    },
})
type HeadingClasses = VariantProps<typeof headingClasses>

interface Props {
    size?: HeadingClasses['size'],
    tag?: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6',
}

const props = withDefaults(defineProps<Props>(), {
    size: 'md',
    tag: 'h2',
})
</script>

<template>
    <component :is="props.tag" :class="headingClasses({ size })">
        <slot />
    </component>
</template>

<style lang="scss" scoped>
.heading {
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

        &-xl {
            font-size: var(--text-xl);
        }
    }
}
</style>
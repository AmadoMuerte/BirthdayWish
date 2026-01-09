<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority'
import { Primitive, type PrimitiveProps } from 'reka-ui'

const stackClasses = cva('stack', {
    variants: {
        gap: {
            xs: 'gap-xs',
            sm: 'gap-sm',
            md: 'gap-md',
            lg: 'gap-lg',
            xl: 'gap-xl',
        },
        align: {
            start: 'items-start',
            center: 'items-center',
            end: 'items-end',
            stretch: 'items-stretch',
        },
        justify: {
            start: 'justify-start',
            end: 'justify-end',
            center: 'justify-center',
            between: 'justify-between',
        },
    },
})

type StackClasses = VariantProps<typeof stackClasses>

interface Props extends /* @vue-ignore */ PrimitiveProps {
    gap?: StackClasses['gap']
    align?: StackClasses['align']
    justify?: StackClasses['justify']
    as?: string
}

const props = withDefaults(defineProps<Props>(), {
    gap: 'md',
    align: 'start',
    as: 'div',
})
</script>

<template>
    <Primitive :as="props.as" :class="stackClasses({ gap, align: props.align, justify: props.justify })"
        v-bind="$attrs">
        <slot />
    </Primitive>
</template>

<style scoped lang="scss">
.stack {
    display: flex;
    flex-direction: column;

    /* GAP */
    &.gap {
        &-xs {
            gap: var(--space-1);
        }

        &-sm {
            gap: var(--space-2);
        }

        &-md {
            gap: var(--space-3);
        }

        &-lg {
            gap: var(--space-4);
        }

        &-xl {
            gap: var(--space-5);
        }
    }

    /* ALIGN */
    &.items {
        &-start {
            align-items: flex-start;
        }

        &-center {
            align-items: center;
        }

        &-end {
            align-items: flex-end;
        }

        &-stretch {
            align-items: stretch;
        }
    }

    /* JUSTIFY */
    &.justify {
        &-start {
            justify-content: flex-start;
        }

        &-end {
            justify-content: flex-end;
        }

        &-center {
            justify-content: center;
        }

        &-between {
            justify-content: space-between;
        }
    }
}
</style>

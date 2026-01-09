<script setup lang="ts">
import { cva, type VariantProps } from 'class-variance-authority'
import { Primitive, type PrimitiveProps } from 'reka-ui'

const groupClasses = cva('group', {
    variants: {
        gap: {
            xs: 'gap-xs',
            sm: 'gap-sm',
            md: 'gap-md',
            lg: 'gap-lg',
            xl: 'gap-xl',
        },
        justify: {
            start: 'justify-start',
            end: 'justify-end',
            center: 'justify-center',
            between: 'justify-between',
            around: 'justify-around',
        },
    },
})

type GroupClasses = VariantProps<typeof groupClasses>

interface Props extends /* @vue-ignore */ PrimitiveProps {
    gap?: GroupClasses['gap']
    justify?: GroupClasses['justify']
    as?: string
}

const props = withDefaults(defineProps<Props>(), {
    gap: 'md',
    as: 'div',

})
</script>

<template>
    <Primitive :as="props.as" :class="groupClasses({ gap, justify })">
        <slot />
    </Primitive>
</template>

<style scoped lang="scss">
.group {
    display: flex;

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

    /* JUSTIFY */
    &.justify {
        &-start {
            justify-content: start;
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

        &-around {
            justify-content: space-around;
        }
    }
}
</style>
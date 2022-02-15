<script lang="ts">
  import { removeToast, storeToast } from "@/store/toast";
  import Icon from "@/components/ui/Icon.svelte";

  const close = (id: number) => {
    removeToast(id);
  };
  const customSlide = (el: HTMLElement, { duration }) => {
    return {
      duration,
      css: (_: any, u: any) => `transform: translateX(${u * 400}px)`,
    };
  };
</script>

<div class="fixed bottom-0 right-0 z-50">
  {#each $storeToast as toast (toast.id)}
    <div
      class={`${toast.type} flex p-2 h-12 items-center border-l border-t border-gray-700 w-80`}
      transition:customSlide={{ duration: 250 }}
    >
      <button on:click={() => close(toast.id)}>
        <Icon icon="right" />
      </button>
      <div class="pl-2">
        <span>{toast.message}</span>
      </div>
    </div>
  {/each}
</div>

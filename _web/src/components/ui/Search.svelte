<script lang="ts">
  import Icon from "./Icon.svelte";

  export let search = "";
  export let searchFn: (_: string) => void;
</script>

<form
  on:submit|preventDefault={() => searchFn(search)}
  class="flex items-center"
>
  <input
    bind:value={search}
    placeholder="search"
    class="py-1 px-2 text-black"
  />
  <button
    class="bg-gray-200 p-1 font-bold inline-block hover:bg-green-500 hover:fill-white"
    type="submit"
  >
    <Icon icon="search" class="pointer-event-none" />
  </button>
</form>

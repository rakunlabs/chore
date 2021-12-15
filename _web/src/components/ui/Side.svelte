<script lang="ts">
  import active from "svelte-spa-router/active";

  export let links = [] as Array<string>;

  let className = "";
  export { className as class };
</script>

<navbar class={`${className} bg-gray-600 text-gray-100`}>
  {#each links as link}
    <button
      data-side={link}
      use:active={{
        path: new RegExp(`/${link}(/(.*))*`),
        className: "sidebar-active",
        inactiveClassName: "sidebar-inactive",
      }}
      class="capitalize">{link}</button
    >
  {/each}
</navbar>

<style lang="scss">
  button {
    @apply pl-2 w-full h-8 text-left;
  }

  :global(.sidebar-active) {
    @apply text-gray-800 bg-gray-100;
  }

  :global(.sidebar-inactive:hover) {
    @apply bg-gray-200 text-gray-800;
  }
</style>

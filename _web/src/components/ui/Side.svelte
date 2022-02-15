<script lang="ts">
  import SideLink from "./SideLink.svelte";

  export let links = [] as Array<string | object>;

  let className = "";
  export { className as class };
</script>

<navbar
  class={`${className} bg-slate-500 text-gray-100 flex flex-col justify-between`}
>
  <div>
    {#each links as link}
      {#if typeof link == "string"}
        <SideLink
          {link}
          class="block capitalize h-8 border-b border-gray-100 pl-2 w-full py-1 text-left"
        />
      {:else if typeof link == "object"}
        <div>
          {#each Object.keys(link) as linkKey}
            <span
              class="block capitalize h-8 bg-slate-600 border-b border-black p-2 w-full py-1 text-left hover:bg-gray-300 hover:text-black"
            >
              {linkKey}
            </span>
            <div class="border-l-4 border-indigo-400">
              {#each link[linkKey] as l}
                <SideLink
                  link={l}
                  class="block capitalize h-8 border-b border-gray-100 pl-2 w-full py-1 text-left"
                />
              {/each}
            </div>
          {/each}
        </div>
      {/if}
    {/each}
  </div>
  <div>
    <button
      data-action="sidebar"
      data-side="logout"
      class="block capitalize border-t border-gray-100 hover:bg-red-500 bg-red-400 pl-2 w-full py-1 text-left"
    >
      Logout
    </button>
  </div>
</navbar>

<style lang="scss">
  :global(.sidebar-active) {
    @apply text-gray-800 bg-gray-100 border-black;
  }

  :global(.sidebar-inactive:hover) {
    @apply bg-gray-50 text-gray-800 border-black;
  }
</style>

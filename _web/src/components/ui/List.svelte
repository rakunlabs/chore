<script lang="ts">
  import Item from "@/components/ui/Item.svelte";
  import type { itemType } from "@/models/template";
  import { onDestroy, onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import path from "path-browserify";

  export let items = [] as Array<itemType>;
  export let prefix = "/";

  let listDiv: HTMLElement;

  const deleteIt = (name: string) => {};

  const catchItem = (e: CustomEvent) => {
    const name = e.detail.name as string;

    // delete request
    if (e.detail.action == "delete") {
      deleteIt(name);
      return;
    }

    // update URL
    push(path.join(prefix, name));
  };

  onMount(() => {
    listDiv.addEventListener("item", catchItem);
  });

  onDestroy(() => {
    listDiv.removeEventListener("item", catchItem);
  });
</script>

<div bind:this={listDiv} class="bg-white">
  {#each items as item}
    <Item
      name={item.name}
      show={item.item}
      type={item.name[item.name.length - 1] == "/" ? "folder" : "file"}
    />
  {/each}
</div>

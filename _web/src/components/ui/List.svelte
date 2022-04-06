<script lang="ts">
  import Item from "@/components/ui/Item.svelte";
  import type { itemType } from "@/models/template";
  import { onDestroy, onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import path from "path-browserify";
  import { requestSender } from "@/helper/api";
  import NoData from "./NoData.svelte";

  export let items = [] as Array<itemType>;
  export let prefix = "/";

  let className = "";

  export { className as class };

  let listDiv: HTMLElement;

  const deleteIt = async (name: string) => {
    if (confirm(`Are you sure to delete ${name}?`)) {
      try {
        await requestSender("template", { name }, "DELETE", null, true);
        // console.log(l);
        items = items.filter((i) => i.name != name);
      } catch (error) {
        console.error(error);
      }
    }
  };

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

<div bind:this={listDiv} class={`bg-white ${className}`}>
  {#each items as item (item.name)}
    <Item
      name={item.name}
      show={item.item}
      type={item.name[item.name.length - 1] == "/" ? "folder" : "file"}
    />
  {/each}
  <NoData hide={!!items.length} />
</div>

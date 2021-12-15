<script lang="ts">
  import Item from "@/components/ui/Item.svelte";
  import { deleteItem } from "@/helper/api";
  import { onDestroy, onMount } from "svelte";
  import { push } from "svelte-spa-router";

  export let items = [] as Array<string>;
  export let input = "";

  let listDiv: HTMLElement;

  const deleteIt = (name: string) => {
    const i = name.indexOf("/");
    try {
      deleteItem(name.slice(0, i), name.slice(i));
      items = items.filter((v) => v != name);
    } catch (error) {
      console.log(error);
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
    push(name[0] == "/" ? name : "/" + name);
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
      name={item}
      show={item.substring(input.length - 1)}
      type={item[item.length - 1] == "/" ? "folder" : "file"}
    />
  {/each}
</div>

<script lang="ts">
  import Icon from "./Icon.svelte";

  type contentTypes = "file" | "folder";

  export let name = "";
  export let show = "";
  export let type: contentTypes = "file";

  let item: HTMLElement;

  const catchClick = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];
    const event = new CustomEvent("item", {
      bubbles: true,
      detail: {
        type: type,
        action: action,
        name: name,
      },
    });

    item.dispatchEvent(event);
  };
</script>

<div
  bind:this={item}
  on:click|stopPropagation={catchClick}
  class="flex flex-row w-full h-10 px-2 items-center justify-between border-b border-gray-400 hover:bg-gray-300 cursor-pointer"
>
  <div>
    <Icon icon={type} class="float-left pr-1" />
    <span>{show}</span>
  </div>
  <div>
    <button
      data-action="delete"
      class="w-20 px-4 bg-transparent border-2 border-red-500 text-red-500 text-sm hover:bg-red-500 hover:text-gray-100"
    >
      Delete
    </button>
  </div>
</div>

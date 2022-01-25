<script lang="ts">
  import Bread from "@/components/ui/Bread.svelte";
  import List from "@/components/ui/List.svelte";
  import Editor from "@/components/ui/Editor.svelte";
  import Loading from "@/components/ui/Loading.svelte";
  import { b64ToUtf8 } from "@/helper/codec";
  import { requestSender } from "@/helper/api";
  import type { itemType } from "@/models/template";

  export let params = {} as Record<string, string>;

  let items: Array<itemType> = [];
  let view = "loading";

  // view
  let data: any;
  let title = "";

  let input = params.input;

  const reg = /^\/templates[/]?/i;

  const getInfo = async (v: string) => {
    view = "loading";

    v = v.replace(reg, "");
    if (v == "/") {
      v = "";
    }

    if (!v || v[v.length - 1] == "/") {
      input = params.input;
      try {
        const l = await requestSender(
          "templates",
          { folder: v },
          "GET",
          null,
          true
        );
        // console.log(l);
        items = l ? l.data.data : [];
      } catch (error) {
        items = [];
      }
      view = "list";
    } else {
      try {
        title = v;
        data = b64ToUtf8(
          (await requestSender("template", { name: v }, "GET", null, true)).data
            .data.content
        );
      } catch (error) {
        data = null;
      }
      view = "data";
    }
  };

  $: getInfo(params.input);

  const addNewItem = () => {
    view = "add";
  };
</script>

<div class="grid [grid-template-rows:auto_1fr] h-full">
  <div
    class="pr-2 py-1 border-b border-gray-400 hover:bg-gray-300 flex justify-between items-center"
  >
    <Bread url={params.input} />
    <button
      on:click|stopPropagation={addNewItem}
      class="w-20 bg-transparent border-2 border-gray-500 text-gray-500 text-sm hover:bg-gray-500 hover:text-gray-100"
    >
      Add
    </button>
  </div>

  <div class="pt-1 h-full">
    {#if view == "list"}
      <List {items} prefix="/templates" />
    {:else if view == "data"}
      <Editor {data} {title} area="templates" />
    {:else if view == "add"}
      <Editor
        title={input.replace(reg, "")}
        area="templates"
        editableTitle={true}
      />
    {:else}
      <Loading />
    {/if}
  </div>
</div>

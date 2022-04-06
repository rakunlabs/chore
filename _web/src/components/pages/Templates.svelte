<script lang="ts">
  import Bread from "@/components/ui/Bread.svelte";
  import List from "@/components/ui/List.svelte";
  import Editor from "@/components/ui/Editor.svelte";
  import Loading from "@/components/ui/Loading.svelte";
  import { b64ToUtf8 } from "@/helper/codec";
  import { requestSender } from "@/helper/api";
  import type { itemType } from "@/models/template";
  import { storeHead } from "@/store/store";

  export let params = {} as Record<string, string>;

  storeHead.set("Templates");

  let items: Array<itemType> = [];
  let view = "loading";

  // view
  let data: any;
  let title = "";

  let input = params.input;

  const reg = /^\/templates[/]?/i;

  const fetchList = async (v: string) => {
    try {
      const l = await requestSender(
        "templates",
        {
          folder: v,
          limit: -1,
        },
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
  };

  const getInfo = async (v: string) => {
    view = "loading";

    v = v.replace(reg, "");
    if (v == "/") {
      v = "";
    }

    if (!v || v[v.length - 1] == "/") {
      input = params.input;
      fetchList(v);
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

  const closed = () => {
    getInfo(params.input);
  };

  const reload = () => {
    getInfo(params.input);
  };
</script>

<div class="h-full w-full grid [grid-template-rows:auto_1fr]">
  <div class="border-b border-black flex justify-between items-center mb-1">
    <Bread url={params.input} />
    <div>
      <button
        on:click|stopPropagation={reload}
        class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
      >
        Reload
      </button>
      <button
        on:click|stopPropagation={addNewItem}
        class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
      >
        Add
      </button>
    </div>
  </div>

  {#if view == "list"}
    <List
      {items}
      prefix="/templates"
      class="h-full min-h-full overflow-x-auto"
    />
  {:else if view == "data"}
    <Editor {data} {title} />
  {:else if view == "add"}
    <Editor title={input.replace(reg, "")} editableTitle={true} {closed} />
  {:else}
    <Loading />
  {/if}
</div>

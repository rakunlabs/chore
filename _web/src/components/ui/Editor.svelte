<script lang="ts">
  import { onMount } from "svelte";
  import CodeMirror from "codemirror";
  import { requestSender } from "@/helper/api";
  import { addToast } from "@/store/toast";
  import axios from "axios";

  let code: HTMLElement;
  export let title = "title";
  export let editableTitle = false;
  export let data = "";
  export let closed = () => void {};
  let className = "";

  let editor: CodeMirror.Editor;

  export { className as class };

  let readOnly = !editableTitle;

  onMount(() => {
    editor = CodeMirror(code, {
      mode: "javascript",
      lineNumbers: true,
      tabSize: 2,
      scrollbarStyle: "native",
      readOnly: readOnly,
      lineWrapping: true,
    });
    editor.setSize("100%", "100%");
    if (readOnly) {
      editor.getWrapperElement().classList.add("bg-yellow-50");
    }
    editor.setValue(data);
  });

  const toggleReadOnly = (v: boolean | undefined = undefined) => {
    editor.setOption("readOnly", v ?? !readOnly);
    readOnly = v ?? !readOnly;

    if (readOnly) {
      editor.getWrapperElement().classList.add("bg-yellow-50");
    } else {
      editor.getWrapperElement().classList.remove("bg-yellow-50");
    }
  };

  const save = async () => {
    try {
      await requestSender(
        "template",
        {
          name: title,
        },
        editableTitle ? "POST" : "PATCH",
        editor.getValue(),
        true
      );
      toggleReadOnly(true);
      addToast("saved template", "info");
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };
</script>

<div class={`h-full w-full grid [grid-template-rows:auto_1fr] ${className}`}>
  <div
    class="px-1 pb-1 bg-gray-100 border-b border-gray-200 flex flex-row items-center justify-between"
  >
    {#if editableTitle}
      <input bind:value={title} />
    {:else}
      <span>{title}</span>
    {/if}
    <div>
      {#if editableTitle}
        <button
          on:click={closed}
          class="px-4 bg-transparent border-2 text-sm border-gray-500 hover:bg-gray-500 hover:text-gray-100"
        >
          Close
        </button>
      {/if}
      <button
        on:click|stopPropagation={() => toggleReadOnly()}
        class={`px-4 bg-transparent border-2 text-sm ${
          readOnly
            ? "border-gray-500 hover:bg-gray-500 hover:text-gray-100"
            : "text-green-500 border-green-500 hover:bg-green-500 hover:text-gray-100"
        }`}
      >
        Edit
      </button>
      <button
        on:click|stopPropagation={save}
        class="px-4 bg-transparent border-2 text-sm border-gray-500 hover:bg-gray-500 hover:text-gray-100"
      >
        Save
      </button>
    </div>
  </div>
  <code class="bg-gray-400 overflow-x-auto" bind:this={code} />
</div>

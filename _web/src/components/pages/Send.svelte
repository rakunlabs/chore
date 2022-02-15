<script lang="ts">
  import { requestSender } from "@/helper/api";

  import { formToObject } from "@/helper/codec";
  import { storeHead } from "@/store/store";
  import axios, { CancelTokenSource } from "axios";
  import CodeMirror from "codemirror";
  import { onMount } from "svelte";

  storeHead.set("Send request to controlflow");

  let formEdit: HTMLFormElement;
  let codeInput: HTMLElement;
  let codeOutput: HTMLElement;
  let editorInput: CodeMirror.Editor;
  let editorOutput: CodeMirror.Editor;

  let error = "";

  let working = false;

  const showData = {
    status: "",
    color: "yellow" as "yellow" | "red" | "green",
  };

  let source: CancelTokenSource;

  const send = async () => {
    if (working) {
      source.cancel("request canceled by user");
      working = false;
      return;
    }

    working = true;
    source = axios.CancelToken.source();

    const data = formToObject(formEdit);

    if (!data["control"] || !data["endpoint"]) {
      error = "control and endpoint required";
      working = false;
      return;
    }
    error = "";

    let msg: unknown;

    try {
      const responseGet = await requestSender(
        "send",
        data,
        "POST",
        editorInput.getValue(),
        true,
        {
          notTransformResponse: true,
          cancelToken: source.token,
        }
      );

      // editorOutput.setValue(responseGet.data);
      msg = responseGet.data;
      showData.status = [
        String(responseGet.status),
        responseGet.statusText,
      ].join(" - ");
      showData.color = "green";
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        msg = reason?.response?.data ?? reason.message;
        showData.status = [
          String(reason?.response?.status),
          reason?.response?.statusText,
        ].join(" - ");
        showData.color = "red";
      } else {
        msg = reason;
        showData.status = "";
        showData.color = "yellow";
      }
    }

    editorOutput.setValue(String(msg));
    working = false;
  };

  onMount(() => {
    editorInput = CodeMirror(codeInput, {
      mode: "text/yaml",
      lineNumbers: true,
      tabSize: 2,
      readOnly: false,
      lineWrapping: true,
    });
    editorInput.setSize("100%", "100%");

    editorOutput = CodeMirror(codeOutput, {
      mode: "text/plain",
      lineNumbers: true,
      tabSize: 2,
      readOnly: true,
      lineWrapping: true,
    });
    editorOutput.setSize("100%", "100%");

    editorOutput.getWrapperElement().classList.add("bg-yellow-50");
    // editor.setValue(data);
  });
</script>

<div class="grid h-full grid-rows-[auto_1fr]">
  <div class="bg-slate-50 p-5 mb-3">
    <div class="flex flex-row flex-wrap justify-between gap-4 items-start">
      <div class="flex-1">
        <form bind:this={formEdit}>
          <label class="mb-1 flex">
            <span class="w-20 inline-block">Control</span>
            <input
              type="text"
              name="control"
              placeholder="mycontrolflow"
              class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
            />
          </label>
          <label class="mb-1 flex">
            <span class="w-20 inline-block">Endpoint</span>
            <input
              type="text"
              name="endpoint"
              placeholder="create"
              class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
            />
          </label>
          <div
            class={`mt-2 bg-red-200 w-full h-6 ${
              error != "" ? "" : "invisible"
            }`}
          >
            <span class="break-all">{error}</span>
          </div>
        </form>
      </div>
      <div class="flex-1 self-stretch">
        <button
          class={`bg-gray-200 p-1 font-bold inline-block w-60 h-full float-right ${
            working
              ? "hover:bg-red-500 hover:text-white"
              : "hover:bg-yellow-200"
          }`}
          on:click={send}>{working ? "Cancel" : "Send"}</button
        >
      </div>
    </div>
  </div>

  <div class="flex gap-3 h-full min-h-full overflow-x-auto">
    <code
      class="flex-1 bg-gray-400 h-full min-h-full overflow-x-auto"
      bind:this={codeInput}
    />
    <div class="flex-1 grid h-full grid-rows-[auto_1fr]">
      <div class="bg-slate-50 p-5 mb-3">
        <label class="flex">
          <span class="w-20 inline-block"> Status </span>
          <input
            type="text"
            class={`flex-grow border border-gray-300 focus:border-gray-300 focus:outline-none focus:ring focus:ring-gray-200 focus:ring-opacity-50 status-${showData.color}`}
            readonly
            value={showData.status}
          />
        </label>
      </div>
      <code
        class="bg-gray-400 h-full min-h-full overflow-x-auto"
        bind:this={codeOutput}
      />
    </div>
  </div>
</div>

<style lang="scss">
  :global(.status-yellow) {
    @apply bg-yellow-50;
  }
  :global(.status-green) {
    @apply bg-green-100;
  }
  :global(.status-red) {
    @apply bg-red-100;
  }
</style>

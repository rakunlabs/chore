<script lang="ts">
  import { requestSender } from "@/helper/api";

  import { formToObject } from "@/helper/codec";
  import { storeHead } from "@/store/store";
  import axios, { CancelTokenSource } from "axios";
  import CodeMirror from "codemirror";
  import { onMount } from "svelte";
  import Code from "@/components/ui/Code.svelte";
  import { fullScreenKeys } from "@/helper/code";

  storeHead.set("Send request to controlflow");

  let URL = window.location.origin + window.location.pathname;
  if (URL[URL.length - 1] == "/") {
    URL = URL.slice(0, URL.length - 1);
  }

  let formEdit: HTMLFormElement;
  let codeInput: HTMLElement;
  let codeOutput: HTMLElement;
  let editorInput: CodeMirror.Editor;
  let editorOutput: CodeMirror.Editor;

  let error = "";

  let working = false;

  let control = "";
  let endpoint = "";

  let beutify = false;
  let originalValue = "";

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
      // reset settings
      beutify = false;
      showData.status = "";
      showData.color = "yellow";
      editorOutput.setValue("");

      const responseGet = await requestSender(
        "send",
        data,
        "POST",
        editorInput.getValue(),
        true,
        {
          notTransformResponse: true,
          cancelToken: source.token,
          timeout: 0,
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

  const showBeautify = () => {
    if (beutify) {
      editorOutput.setValue(originalValue);
      beutify = false;
    } else {
      originalValue = editorOutput.getValue();
      try {
        const parsedValue = JSON.parse(originalValue);
        editorOutput.setValue(JSON.stringify(parsedValue, null, "  "));
      } catch (v: unknown) {
        return;
      }
      beutify = true;
    }
  };

  onMount(() => {
    editorInput = CodeMirror(codeInput, {
      mode: "yaml",
      lineNumbers: true,
      tabSize: 2,
      readOnly: false,
      lineWrapping: true,
      styleActiveLine: true,
      matchBrackets: true,
      showTrailingSpace: true,
      placeholder:
        "input value could be anything\nyaml/json supported in template\n\nF11 full-screen",
      extraKeys: fullScreenKeys,
    });
    editorInput.setSize("100%", "100%");

    editorOutput = CodeMirror(codeOutput, {
      mode: "yaml",
      lineNumbers: true,
      tabSize: 2,
      readOnly: true,
      lineWrapping: true,
      styleActiveLine: true,
      matchBrackets: true,
      showTrailingSpace: true,
      placeholder: "output of respond\n\nF11 full-screen",
      extraKeys: fullScreenKeys,
    });
    editorOutput.setSize("100%", "100%");

    editorOutput.getWrapperElement().classList.add("bg-yellow-50");
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
              bind:value={control}
              class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
            />
          </label>
          <label class="mb-1 flex">
            <span class="w-20 inline-block">Endpoint</span>
            <input
              type="text"
              name="endpoint"
              placeholder="create"
              bind:value={endpoint}
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
    <Code
      lang="sh"
      code={`curl -ksSL -X POST -H "Authorization: Bearer \${TOKEN}" --data-binary @filename "${URL}/api/v1/send?control=${control}&endpoint=${endpoint}"`}
    />
  </div>

  <div class="flex gap-3 h-full min-h-full overflow-x-auto">
    <code class="flex-1 bg-gray-400 overflow-x-auto" bind:this={codeInput} />
    <div class="flex-1 grid h-full grid-rows-[auto_1fr]">
      <div class="bg-slate-50 p-5 mb-3 flex gap-1">
        <label class="flex-1 flex">
          <span class="pr-2">Status</span>
          <input
            type="text"
            class={`flex-1 border border-gray-300 focus:border-gray-300 focus:outline-none focus:ring focus:ring-gray-200 focus:ring-opacity-50 status-${showData.color}`}
            readonly
            value={showData.status}
          />
        </label>
        <button
          class="px-2 border border-gray-300 bg-yellow-200 hover:bg-green-300"
          on:click={showBeautify}>{beutify ? "Original" : "Beutify"}</button
        >
      </div>
      <code class="bg-gray-400 overflow-x-auto" bind:this={codeOutput} />
    </div>
  </div>
</div>

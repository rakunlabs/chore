<script lang="ts">
  import { onMount } from "svelte";
  import CodeMirror from "codemirror";
  import { fullScreenKeys } from "@/helper/code";
  import { requestSender } from "@/helper/api";
  import type { CancelTokenSource } from "axios";
  import axios from "axios";
  import { Base64 } from "js-base64";

  let codeEditorScript: CodeMirror.Editor;
  let codeEditorInputs: CodeMirror.Editor;
  let codeEditorOutput: CodeMirror.Editor;
  let codeElementScript: HTMLElement;
  let codeElementInput: HTMLElement;
  let codeElementOutput: HTMLElement;
  let codeChanged = false;
  let information = "";

  export let showEditor = false;
  export let showEditorChange: (v: boolean) => void;
  export let codeEditorSave: (script: string, inputs: string) => void;

  export const setCodeEditorValue = (
    script: string | null,
    inputs: string | null,
    info: string
  ) => {
    codeEditorScript?.setValue(script ?? "");
    codeEditorInputs?.setValue(inputs ?? "");
    codeChanged = false;

    information = info;
  };

  const codeSetting = {
    mode: "javascript",
    lineNumbers: true,
    tabSize: 2,
    lineWrapping: true,
    styleActiveLine: true,
    matchBrackets: true,
    showTrailingSpace: true,
    placeholder: "javascript\n\nF11 full-screen",
    extraKeys: fullScreenKeys,
  };

  let working = false;
  let source: CancelTokenSource;

  const callScript = async () => {
    if (working) {
      source.cancel("request canceled by user");
      working = false;
      return;
    }

    working = true;
    source = axios.CancelToken.source();

    const script = codeEditorScript.getValue();
    const inputs = codeEditorInputs.getValue();

    let output: unknown;

    try {
      const responseRun = await requestSender(
        "run/js",
        null,
        "POST",
        {
          script: Base64.encode(script),
          inputs: Base64.encode(inputs),
        },
        true,
        {
          notTransformResponse: true,
          cancelToken: source.token,
          timeout: 0,
        }
      );

      output = responseRun.data;
    } catch (reason: unknown) {
      console.log(reason);
      if (axios.isAxiosError(reason)) {
        output = reason?.response?.data ?? reason.message;
      } else {
        output = reason;
      }
    }

    codeEditorOutput.setValue(String(output));
    working = false;
  };

  onMount(() => {
    // set code editor script
    codeEditorScript = CodeMirror(codeElementScript, codeSetting);
    codeEditorScript.setSize("100%", "100%");

    codeEditorScript.on("change", () => {
      codeChanged = true;
    });

    // set code editor input
    codeEditorInputs = CodeMirror(codeElementInput, codeSetting);
    codeEditorInputs.setSize("100%", "100%");

    codeEditorInputs.on("change", () => {
      codeChanged = true;
    });

    // set code editor output
    codeEditorOutput = CodeMirror(codeElementOutput, codeSetting);
    codeEditorOutput.setSize("100%", "100%");
  });
</script>

<div
  class={`block absolute h-full w-full z-40 left-0 top-0 bg-slate-500 bg-opacity-40 ${
    showEditor ? "" : "invisible"
  }`}
>
  <div class="w-full m-auto h-full grid grid-rows-[auto_1fr]">
    <div
      class={`flex justify-between border-b border-black ${
        codeChanged ? "bg-yellow-200" : "bg-green-200"
      }`}
    >
      <span class="px-1 text-black">{information}</span>
      <div class="text-[0]">
        <button
          class="w-40 text-base bg-gray-100 hover:bg-red-500 hover:text-white border-l border-black"
          on:click|stopPropagation={() => {
            codeEditorSave(
              codeEditorScript.getValue(),
              codeEditorInputs.getValue()
            );
            codeChanged = false;
          }}>Save</button
        >
        <button
          class="w-40 text-base bg-gray-100 hover:bg-red-500 hover:text-white border-l border-black"
          on:click|stopPropagation={() => showEditorChange(false)}>Close</button
        >
      </div>
    </div>
    <div class="w-full min-h-full h-full grid grid-cols-2">
      <div class="w-full h-full overflow-auto grid grid-rows-[auto_1fr]">
        <div class="bg-gray-200 pl-1 flex flex-row justify-between">
          <span>Script</span>
          <button
            class="px-4 bg-yellow-200 hover:bg-green-400"
            on:click|stopPropagation={callScript}
          >
            Run
          </button>
        </div>
        <code bind:this={codeElementScript} class="overflow-auto" />
      </div>
      <div
        class="w-full h-full overflow-auto grid grid-rows-2 border-l border-black"
      >
        <div class="w-full grid grid-rows-[auto_1fr]">
          <div class="bg-gray-200 px-1">
            <span>Input</span>
          </div>
          <code bind:this={codeElementInput} class="overflow-auto" />
        </div>
        <div
          class="w-full h-full overflow-auto grid grid-rows-[auto_1fr] border-t border-black"
        >
          <div class="bg-gray-200 px-1">
            <span>Output</span>
          </div>
          <code bind:this={codeElementOutput} class="overflow-auto" />
        </div>
      </div>
    </div>
  </div>
</div>

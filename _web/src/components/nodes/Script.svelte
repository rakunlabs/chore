<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { scriptData } from "@/models/nodes/script";
  import NodeSave from "../ui/NodeSave.svelte";
  import CodeEditor from "../ui/CodeEditor.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: scriptData;
  let countOfInputs: number;

  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as scriptData;
    countOfInputs = Object.keys(nodeV.inputs).length;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.script = formData.get("script") as string;
    v.info = formData.get("info") as string;
    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };

  let showEditor = false;
  let showEditorChange = (v: boolean) => {
    showEditor = v;
  };

  let codeEditorSave = (script: string, inputs: string) => {
    data.script = script;
    data.inputs = inputs;
  };

  let setCodeEditorValue: (
    script: string,
    inputs: string,
    info: string
  ) => void;

  const openEditor = () => {
    setCodeEditorValue(data.script, data.inputs, data.info);
    showEditor = true;
  };
</script>

<CodeEditor
  bind:setCodeEditorValue
  {codeEditorSave}
  {showEditor}
  {showEditorChange}
/>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>Info for UI</p>
  <input type="text" placeholder="info" name="info" bind:value={data.info} />
  <button
    on:click={openEditor}
    class="w-full bg-yellow-50 hover:bg-blue-200 border border-gray-400"
    >Open Editor</button
  >
  <p>Enter script</p>
  <textarea
    class="h-56"
    placeholder="script"
    name="script"
    bind:value={data.script}
  />
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { scriptData } from "@/models/nodes/script";
  import NodeSave from "../ui/NodeSave.svelte";
  import CodeEditor from "../ui/CodeEditor.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;
  export let nodeUnselected: () => void;

  let data: scriptData;
  let countOfInputs: number;

  let inputCount: number = 0;
  let setInputCount: boolean = false;

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

    if (setInputCount) {
      const count = +(formData.get("input_count") as string);

      editor.updateNodeDataFromId(node.id, v);

      const nodeNew = editor.getNodeFromId(node.id);
      editor.removeNodeId("node-" + node.id);

      const newID = editor.addNode(
        "script",
        count,
        Object.keys(nodeNew.outputs).length,
        nodeNew.pos_x,
        nodeNew.pos_y,
        nodeNew.class,
        nodeNew.data,
        nodeNew.html,
        nodeNew.typenode
      );

      editor.getNodeFromId(newID).inputs["input_1"].connections = [
        ...nodeNew.inputs["input_1"].connections,
      ];

      Object.keys(nodeNew.inputs).forEach((input, index) => {
        nodeNew.inputs[input].connections.forEach((connection) => {
          if (index >= count) {
            return;
          }

          editor.addConnection(connection.node, newID, connection.input, input);
        });
      });

      Object.keys(nodeNew.outputs).forEach((output) => {
        nodeNew.outputs[output].connections.forEach((connection) => {
          editor.addConnection(
            newID,
            connection.node,
            output,
            (connection as any).output
          );
        });
      });

      nodeUnselected();

      return;
    }

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
  <p class="title-node">Script - {node.id}</p>
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
  <p>Enter input count</p>
  <div class="flex justify-between items-center">
    <input
      class="flex-1 min-w-[2rem]"
      type="checkbox"
      bind:checked={setInputCount}
    />
    <input
      class="flex-auto"
      type="number"
      placeholder="0"
      name="input_count"
      disabled={!setInputCount}
      bind:value={inputCount}
    />
  </div>
  <NodeSave />
</form>

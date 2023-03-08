<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { logData } from "@/models/nodes/log";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: logData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as logData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.message = formData.get("message") as string;
    v.data = formData.get("data") != null;
    v.level = formData.get("level") as string;

    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>Message</p>
  <input
    type="text"
    name="message"
    placeholder="awesome log message"
    bind:value={data.message}
  />
  <p>Log Level</p>
  <select name="level" bind:value={data.level}>
    <option value="debug">Debug</option>
    <option value="info">Info</option>
    <option value="warn">Warn</option>
    <option value="error">Error</option>
    <option value="">NoLevel</option>
  </select>
  <hr />
  <label>
    <span>Print data</span>
    <input
      type="checkbox"
      name="data"
      data-action="checkbox"
      bind:checked={data.data}
    />
  </label>
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

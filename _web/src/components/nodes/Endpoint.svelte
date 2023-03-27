<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { endpointData } from "@/models/nodes/endpoint";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: endpointData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as endpointData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.endpoint = formData.get("endpoint") as string;
    v.methods = formData.get("methods") as string;
    v.public = formData.get("public") != null;

    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p class="title-node">Endpoint - {node.id}</p>
  <p>Enter endpoint name</p>
  <input
    type="text"
    placeholder="create"
    name="endpoint"
    bind:value={data.endpoint}
  />
  <p>Methods</p>
  <input
    type="text"
    placeholder="POST, GET"
    name="methods"
    bind:value={data.methods}
  />
  <label>
    <span>Public</span>
    <input
      type="checkbox"
      name="public"
      data-action="checkbox"
      bind:checked={data.public}
    />
  </label>
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

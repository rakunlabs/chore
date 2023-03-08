<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { controlData } from "@/models/nodes/control";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: controlData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as controlData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.control = formData.get("control") as string;
    v.endpoint = formData.get("endpoint") as string;
    v.method = formData.get("method") as string;
    v.info = formData.get("info") as string;

    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>Info for UI</p>
  <input type="text" placeholder="info" name="info" bind:value={data.info} />
  <p>Enter control name</p>
  <input type="text" name="control" bind:value={data.control} />
  <p>Enter endpoint name</p>
  <input type="text" name="endpoint" bind:value={data.endpoint} />
  <p>Enter method name</p>
  <input type="text" name="method" bind:value={data.method} />
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

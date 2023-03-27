<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { forLoopData } from "@/models/nodes/forLoop";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: forLoopData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as forLoopData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.for = formData.get("for") as string;
    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p class="title-node">For - {node.id}</p>
  <p>Expression</p>
  <input type="text" placeholder="data" name="for" bind:value={data.for} />
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

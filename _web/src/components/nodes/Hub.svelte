<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { hubData } from "@/models/nodes/hub";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: hubData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as hubData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p class="title-node">Hub - {node.id}</p>
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { respondData } from "@/models/nodes/respond";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: respondData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as respondData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.get = formData.get("get") != null;
    v.headers = formData.get("headers") as string;
    v.status = formData.get("status") as string;
    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>Enter respond status code</p>
  <input class="mr-2" type="number" name="status" bind:value={data.status} />
  <p>Enter headers</p>
  <textarea placeholder="json/yaml key:value" bind:value={data.headers} />
  <hr />
  <label>
    <span>Get respond in data</span>
    <input
      type="checkbox"
      name="get"
      data-action="checkbox"
      bind:checked={data.get}
    />
  </label>
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

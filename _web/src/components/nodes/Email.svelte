<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { emailData } from "@/models/nodes/email";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let data: emailData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as emailData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.to = formData.get("email-to") as string;
    v.bcc = formData.get("email-bcc") as string;
    v.cc = formData.get("email-cc") as string;
    v.from = formData.get("email-from") as string;
    v.subject = formData.get("email-subject") as string;

    v.tags = formData.get("tags") as string;

    editor.updateNodeDataFromId(node.id, v);
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>From</p>
  <input
    type="text"
    name="email-from"
    placeholder="system defined"
    bind:value={data.from}
  />
  <p>To</p>
  <input type="text" name="email-to" bind:value={data.to} />
  <p>CC</p>
  <input type="text" name="email-cc" bind:value={data.cc} />
  <p>BCC</p>
  <input type="text" name="email-bcc" bind:value={data.bcc} />
  <p>Subject</p>
  <input type="text" name="email-subject" bind:value={data.subject} />
  <p>Enter tags</p>
  <input type="text" placeholder="tags" name="tags" bind:value={data.tags} />
  <NodeSave />
</form>

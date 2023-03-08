<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import type { noteData } from "@/models/nodes/note";
  import NodeSave from "../ui/NodeSave.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;

  let formElement: HTMLFormElement;

  let data: noteData;
  const getData = (nodeV: DrawflowNode) => {
    data = nodeV.data as noteData;
  };

  $: getData(node);

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const v = Object.assign({}, data);

    v.note = formData.get("note") as string;
    v.backgroundColor = formData.get("note-bg-color") as string;
    v.textColor = formData.get("note-text-color") as string;
    v.width = +(formData.get("note-width") as string);
    v.height = +(formData.get("note-height") as string);

    editor.updateNodeDataFromId(node.id, v);

    formElement.dispatchEvent(
      new CustomEvent("nodeUpdated", {
        detail: node.id,
        bubbles: true,
        cancelable: true,
      })
    );
  };

  const reset = () => {
    data = editor.getNodeFromId(node.id).data;
  };
</script>

<form
  on:submit|preventDefault={submit}
  on:reset|preventDefault={reset}
  bind:this={formElement}
>
  <p>Note</p>
  <label>
    <span>Background color</span>
    <input
      type="color"
      placeholder="note"
      name="note-bg-color"
      bind:value={data.backgroundColor}
    />
  </label>
  <label>
    <span>Text color</span>
    <input
      type="color"
      placeholder="note"
      name="note-text-color"
      bind:value={data.textColor}
    />
  </label>
  <br />
  <label>
    <span>Width</span>
    <input
      type="number"
      placeholder="width"
      name="note-width"
      min="10"
      bind:value={data.width}
    />
  </label>
  <label>
    <span>Height</span>
    <input
      type="number"
      placeholder="height"
      name="note-height"
      min="10"
      bind:value={data.height}
    />
  </label>
  <NodeSave />
</form>

<script lang="ts">
  import type Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import Endpoint from "@/components/nodes/Endpoint.svelte";
  import Template from "@/components/nodes/Template.svelte";
  import Request from "@/components/nodes/Request.svelte";
  import Script from "@/components/nodes/Script.svelte";
  import ForLoop from "@/components/nodes/ForLoop.svelte";
  import IfCase from "@/components/nodes/IfCase.svelte";
  import Control from "@/components/nodes/Control.svelte";
  import Respond from "@/components/nodes/Respond.svelte";
  import Wait from "@/components/nodes/Wait.svelte";
  import Log from "@/components/nodes/Log.svelte";
  import Email from "@/components/nodes/Email.svelte";
  import Note from "@/components/nodes/Note.svelte";
  import NodeView from "@/components/ui/NodeView.svelte";

  export let node: DrawflowNode;
  export let editor: Drawflow;
</script>

{#if node == null}
  <NodeView />
{/if}
{#if node?.name == "endpoint"}
  <Endpoint {node} {editor} />
{/if}
{#if node?.name == "template"}
  <Template {node} {editor} />
{/if}
{#if node?.name == "request"}
  <Request {node} {editor} />
{/if}
{#if node?.name == "script"}
  <Script {node} {editor} />
{/if}
{#if node?.name == "forLoop"}
  <ForLoop {node} {editor} />
{/if}
{#if node?.name == "ifCase"}
  <IfCase {node} {editor} />
{/if}
{#if node?.name == "control"}
  <Control {node} {editor} />
{/if}
{#if node?.name == "respond"}
  <Respond {node} {editor} />
{/if}
{#if node?.name == "log"}
  <Log {node} {editor} />
{/if}
{#if node?.name == "email"}
  <Email {node} {editor} />
{/if}
{#if node?.name == "note"}
  <Note {node} {editor} />
{/if}
{#if node?.name == "wait"}
  <Wait {node} {editor} />
{/if}

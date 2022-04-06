<script lang="ts">
  import { hljs } from "@/helper/highlight";
  import Icon from "./Icon.svelte";

  export let lang = "sh";
  export let code = "";

  let value = "";

  const hightlightValue = (c: string) => {
    value = hljs.highlight(c, { language: lang }).value;
  };

  $: hightlightValue(code);

  let copied = false;

  // copy to clipboard part
  const copy = () => {
    console.log(code);
    navigator.clipboard.writeText(code).then(
      () => {
        copied = true;
        setTimeout(() => (copied = false), 500);
      },
      () => {
        console.warn("failed copy");
      }
    );
  };
</script>

<div class="w-full flex justify-center my-2 [font-size:80%]">
  <button
    class="px-1 bg-gray-100 border-r border-gray-400 hover:bg-gray-200"
    title="copy to clipboard"
    on:click|stopPropagation={copy}
  >
    {#if copied}
      <Icon icon="ok" width="20" class="w-5" />
    {:else}
      <Icon icon="copy" vWidth="20" vHeight="14" class="w-5" />
    {/if}
  </button>
  <pre class="flex-1 grid"><code class={`${lang} language-${lang} hljs`}
      >{@html value}</code
    ></pre>
</div>

<script lang="ts">
  import { hljs, stringToHTML, converter } from "@/helper/md";

  export let md = "";

  let introCnv: HTMLElement;

  const render = (v: string) => {
    introCnv = stringToHTML(converter.makeHtml(v));
    introCnv.querySelectorAll("pre code").forEach((block: HTMLElement) => {
      hljs.highlightElement(block);
    });
  };

  $: render(md);
</script>

<div class="prose prose-stone">
  {@html introCnv?.innerHTML}
</div>

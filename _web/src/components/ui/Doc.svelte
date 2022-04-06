<script lang="ts">
  import { stringToHTML, converter } from "@/helper/md";
  import { hljs } from "@/helper/highlight";
  import { afterUpdate } from "svelte";
  import { generateToc } from "@/helper/toc";
  import type { tocT } from "@/helper/toc";
  import Toc from "./Toc.svelte";

  export let md = "";

  let introCnv: HTMLElement;
  let articles: tocT[];

  const render = (v: string) => {
    const converted = stringToHTML(converter.makeHtml(v));
    // highlight codes
    converted.querySelectorAll("pre code").forEach((block: HTMLElement) => {
      hljs.highlightElement(block);
    });

    converted.querySelectorAll("pre, img").forEach((block: HTMLElement) => {
      block.classList.add("md-margin-y");
    });
    introCnv = converted;
  };

  $: render(md);

  afterUpdate(() => {
    if (location.hash) {
      window.location.href = location.href;
      articles = generateToc(introCnv);
    }
  });
</script>

<div class="flex">
  <div class="md flex-1">
    {@html introCnv?.innerHTML}
  </div>

  <div class="pl-4">
    <Toc {articles} class="sticky top-5" />
  </div>
</div>

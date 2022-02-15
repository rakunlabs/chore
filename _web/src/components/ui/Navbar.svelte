<script lang="ts">
  import TokenRenew from "@/components/ui/TokenRenew.svelte";

  let className = "";

  export { className as class };
</script>

<div
  class={`${className} bg-slate-600 flex flex-row items-center h-full justify-between pl-2 pr-2`}
>
  <a
    href="#/"
    class="font-sans font-bold text-3xl text-gray-100 hover:text-yellow-200"
  >
    chore
  </a>
  <div class="flex items-center gap-2">
    <TokenRenew />
    <a
      href="./api/v1/swagger"
      target="_blank"
      class="px-4 py-1 my-auto bg-transparent border-2 border-green-500 text-green-500 text-sm hover:bg-green-500 hover:text-gray-100"
    >
      swagger-ui
    </a>
  </div>
</div>

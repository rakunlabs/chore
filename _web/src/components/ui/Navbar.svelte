<script lang="ts">
  let className = "";

  export { className as class };
</script>

<div
  class={`${className} bg-gray-800 flex flex-row align-middle h-full justify-between pl-2 pr-2`}
>
  <button data-side="" class="font-sans font-bold text-3xl text-gray-100">
    chore
  </button>
  <a
    href="./api/v1/swagger"
    target="_blank"
    class="px-4 py-1 my-auto bg-transparent border-2 border-green-500 text-green-500 text-sm hover:bg-green-500 hover:text-gray-100 focus:border-4 focus:border-green-300"
  >
    swagger-ui
  </a>
</div>

<style lang="scss">
</style>

<script lang="ts">
  import active from "svelte-spa-router/active";
  export let link = "";

  let className = "";

  export { className as class };
</script>

<a
  href={`#/${link}`}
  use:active={{
    path: new RegExp(`/${link}(/(.*))*`),
    className: "sidebar-active",
    inactiveClassName: "sidebar-inactive",
  }}
  class={className}
>
  {link}
</a>

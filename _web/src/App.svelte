<script>
  import Home from "./components/Home.svelte";
  import Router from "svelte-spa-router";
  import Login from "./components/Login.svelte";

  const routes = {
    "/login": Login,
    "*": Home,
  };
</script>

<Router {routes} />

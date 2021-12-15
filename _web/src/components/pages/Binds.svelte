<script>
  import Table from "../ui/Table.svelte";

  const head = ["id", "authentication", "template"];

  const data = [
    {
      id: "bind1",
      authentication: "test1",
      template: "temp",
    },
  ];
</script>

<span>Bind authentication and template</span>

<Table {head} {data} />

<script lang="ts">
  export let head = [] as Array<string>;
  export let data = [] as Array<object>;
</script>

<div class="flex flex-col">
  <div class="overflow-x-auto">
    <div class="py-2 inline-block min-w-full">
      <div class="overflow-hidden">
        <table class="min-w-full">
          <thead class="border-b bg-gray-800">
            <tr>
              {#each head as h}
                <th
                  scope="col"
                  class="text-sm font-medium text-white px-6 py-4 text-left"
                >
                  {h}
                </th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each data as d}
              <tr
                class="bg-white border-b transition duration-300 ease-in-out hover:bg-yellow-100"
              >
                {#each head as h}
                  <td
                    class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900"
                  >
                    {d[h]}
                  </td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>

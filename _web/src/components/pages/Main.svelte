<script lang="ts">
  import Doc from "@/components/ui/Doc.svelte";
  import { requestSender } from "@/helper/api";
  import { storeHead } from "@/store/store";
  import { addToast } from "@/store/toast";
  import axios from "axios";
  import { onMount } from "svelte";

  storeHead.set("Main Page");

  let data = {} as Record<string, string>;

  const getIntro = async () => {
    try {
      const response = await requestSender(
        "./info/intro.md",
        null,
        "GET",
        null,
        true,
        {
          rawArea: true,
        }
      );
      data = { ...data, intro: response.data };
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  onMount(() => getIntro());
</script>

<div class="bg-slate-50 p-5 h-full">
  <Doc md={data["intro"]} />
</div>

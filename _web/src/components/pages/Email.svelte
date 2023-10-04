<script lang="ts">
  import { requestSender } from "@/helper/api";
  import { storeHead } from "@/store/store";
  import { onMount } from "svelte";
  import { formToObject } from "@/helper/codec";
  import axios from "axios";
  import { addToast } from "@/store/toast";

  storeHead.set("Users");

  let data: Record<string, any> = {};
  const error = "";

  const getSettings = async () => {
    try {
      const l = await requestSender(
        "settings",
        { namespace: "email", name: "email-1" },
        "GET",
        null,
        true,
        {
          noAlert: true,
        }
      );
      data = l.data.data?.data;
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  const setSettings = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const data = formToObject(e.currentTarget);

    // delete unused fields
    for (const key of ["password"]) {
      if (data[key] == "") {
        delete data[key];
      }
    }

    data["no_auth"] = !!data["no_auth"];

    try {
      await requestSender(
        "settings",
        { namespace: "email", name: "email-1" },
        "PATCH",
        data,
        true
      );
      addToast("settings saved", "info");
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  onMount(() => {
    getSettings();
  });
</script>

<div class="bg-slate-50 p-5 mb-3">
  <div class="flex flex-row flex-wrap gap-4">
    <div class="flex-1">
      <div class="flex justify-between">
        <span class="font-bold block">Email Setting</span>
        <div>
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() => getSettings()}>Reload</button
          >
        </div>
      </div>

      <hr class="mb-4" />

      <form on:submit|preventDefault|stopPropagation={setSettings}>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Email</span>
          <input
            type="email"
            name="email"
            placeholder="user@ingenico.com"
            value={data?.email ?? ""}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Host</span>
          <input
            type="text"
            name="host"
            placeholder="smtp.office365.com"
            value={data?.host ?? ""}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Port</span>
          <input
            type="number"
            name="port"
            value={data?.port ? data?.port : undefined}
            placeholder="587"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Password</span>
          <input
            type="password"
            name="password"
            autocomplete="off"
            value={data?.password ?? ""}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">NoAuth</span>
          <input
            type="checkbox"
            name="no_auth"
            autocomplete="off"
            checked={!!data?.no_auth}
            class="self-center px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <button
          type="submit"
          name="action"
          value="edit"
          class="w-full inline-flex items-center justify-center px-4 py-1 text-black bg-yellow-200 font-semibold capitalize hover:text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
        >
          Edit
        </button>
        <div
          class={`mt-2 bg-red-200 w-full h-6 ${error != "" ? "" : "invisible"}`}
        >
          <span class="break-all">{error}</span>
        </div>
      </form>
    </div>
  </div>
</div>

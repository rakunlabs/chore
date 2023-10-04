<script lang="ts">
  import { requestSender } from "@/helper/api";
  import { storeHead } from "@/store/store";
  import { onMount } from "svelte";
  import { formToObject } from "@/helper/codec";
  import axios from "axios";
  import { addToast } from "@/store/toast";

  storeHead.set("Oauth2 settings");

  // data => name, data
  let datas: Record<string, any>[] = [];
  const error = "";

  const newSetting = () => {
    if (datas.some((v) => v == null)) {
      return;
    }
    datas = [...datas, null];
  };

  const deleteSetting = async (i: number) => {
    if (!confirm(`Are you sure to delete ${datas[i]?.name}?`)) {
      return;
    }

    try {
      await requestSender(
        "settings",
        { namespace: "oauth2", name: datas[i]?.name },
        "DELETE",
        null,
        true,
        {
          noAlert: true,
        }
      );

      datas.splice(i, 1);
      datas = datas;
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  const getSettings = async () => {
    try {
      const l = await requestSender(
        "settings",
        { namespace: "oauth2" },
        "GET",
        null,
        true,
        {
          noAlert: true,
        }
      );
      datas = l.data.data ?? [];
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  const getSetting = async (name: string) => {
    try {
      const l = await requestSender(
        "settings",
        { namespace: "oauth2", name: name },
        "GET",
        null,
        true,
        {
          noAlert: true,
        }
      );

      return l.data?.data;
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }

    return null;
  };

  const setSettings = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const data = formToObject(e.currentTarget);

    let name = data["name"];
    delete data["name"];

    try {
      await requestSender(
        "settings",
        { namespace: "oauth2", name: name },
        "PATCH",
        data,
        true
      );

      const dataCreated = await getSetting(name);
      if (dataCreated == null) {
        return;
      }

      datas = datas.filter((data) => data != null);
      datas.unshift(dataCreated);

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
        <span class="font-bold block">Oauth2 Settings</span>
        <div>
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={newSetting}>New</button
          >
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={getSettings}>Reload</button
          >
        </div>
      </div>

      <hr class="mb-4" />

      <div>
        {#each datas as data, i}
          <form on:submit|preventDefault|stopPropagation={setSettings}>
            <hr class="mb-2" />
            <div class="flex justify-end">
              <button
                type="button"
                class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
                on:click|stopPropagation={() => deleteSetting(i)}
              >
                Delete
              </button>
            </div>
            <label class="mb-1 flex">
              <span class="w-20 inline-block">Name</span>
              <input
                type="text"
                name="name"
                autocomplete="off"
                placeholder="my-auth"
                value={data?.name ?? ""}
                class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              />
            </label>
            <label class="mb-1 flex">
              <span class="w-20 inline-block">Client ID</span>
              <input
                type="text"
                name="client_id"
                placeholder="auth_ui"
                value={data?.data?.client_id ?? ""}
                class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              />
            </label>
            <label class="mb-1 flex">
              <span class="w-20 inline-block">Client Secret</span>
              <input
                type="password"
                name="client_secret"
                autocomplete="off"
                value={data?.data?.client_secret ?? ""}
                class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              />
            </label>
            <label class="mb-1 flex">
              <span class="w-20 inline-block">Token URL</span>
              <input
                type="text"
                name="token_url"
                placeholder="http://localhost:8082/realms/master/protocol/openid-connect/token"
                value={data?.data?.token_url ?? ""}
                class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              />
            </label>
            <label class="mb-1 flex">
              <span class="w-20 inline-block">Scope</span>
              <input
                type="text"
                name="scope"
                placeholder="profile, x-scope"
                value={data?.data?.scope ?? ""}
                class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              />
            </label>
            <button
              type="submit"
              name="action"
              value="save"
              class="w-full inline-flex items-center justify-center px-4 py-1 text-black bg-yellow-200 font-semibold capitalize hover:text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
            >
              Save
            </button>
            <div
              class={`mt-2 bg-red-200 w-full h-6 ${
                error != "" ? "" : "invisible"
              }`}
            >
              <span class="break-all">{error}</span>
            </div>
          </form>
        {/each}
      </div>
    </div>
  </div>
</div>

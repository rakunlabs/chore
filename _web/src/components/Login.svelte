<script lang="ts">
  import { banner } from "@/helper/banner";
  import { formToObject } from "@/helper/codec";
  import { login } from "@/helper/login";
  import { pushRedirect } from "@/helper/push";
  import { tokenSet } from "@/helper/token";
  import { storeInfo } from "@/store/store";
  import axios from "axios";
  import { querystring } from "svelte-spa-router";

  let error = "";
  let working = false;

  let selected = "";

  const setSelected = (v: string) => {
    selected = v;
  };

  const signin = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    // prevent multiple click
    if (working) {
      return;
    }

    working = true;
    const data = formToObject(e.currentTarget);
    try {
      const response = await login(
        data,
        selected ? { provider: selected } : null
      );
      tokenSet(response.data.data, selected);
      pushRedirect($querystring);
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        error = reason?.response?.data?.error ?? reason.message;
      } else {
        error = reason as any;
      }
    }

    working = false;
  };
</script>

<div
  class="w-full min-h-screen bg-gray-50 flex flex-col items-center pt-6 sm:pt-0"
>
  <div class="w-full sm:max-w-md p-5 mx-auto">
    <h2 class="mb-8 text-center text-sm font-extrabold [line-height:1.2]">
      {banner}
    </h2>
    {#if $storeInfo.providers.length > 0}
      <div>
        <button
          class={`border border-b-0 py-1 px-3 ${
            selected == "" ? "bg-yellow-50" : "bg-gray-100"
          }`}
          on:click={() => setSelected("")}
        >
          chore
        </button>
        {#each $storeInfo.providers as provider}
          <button
            class={`border border-b-0 py-1 px-3 ${
              selected == provider ? "bg-yellow-50" : "bg-gray-100"
            }`}
            on:click={() => setSelected(provider)}
          >
            {provider}
          </button>
        {/each}
      </div>
    {/if}

    <div class="border p-4 bg-yellow-50 relative">
      <span class="absolute top-0 right-0 bg-slate-100 px-1"
        >{$storeInfo.version +
          ($storeInfo.environment ? " - " + $storeInfo.environment : "")}</span
      >
      <form on:submit|preventDefault|stopPropagation={signin}>
        <div class="mb-4">
          <label class="block mb-1" for="login"> Username </label>
          <input
            id="login"
            type="text"
            name="login"
            class="py-2 px-3 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100 mt-1 block w-full"
          />
        </div>
        <div class="mb-4">
          <label class="block mb-1" for="password">Password</label>
          <input
            id="password"
            type="password"
            name="password"
            class="py-2 px-3 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100 mt-1 block w-full"
          />
        </div>
        <div class="mt-6">
          <button
            type="submit"
            class="w-full inline-flex items-center justify-center px-4 py-2 bg-red-400 border border-transparent font-semibold capitalize text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:bg-gray-400 transition"
            disabled={working}
          >
            Sign In
          </button>
        </div>
        {#if error != ""}
          <div class="mt-4 bg-red-200">
            <span class="break-all">{error}</span>
          </div>
        {/if}
      </form>
    </div>
  </div>
</div>

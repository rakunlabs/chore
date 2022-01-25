<script lang="ts">
  import { banner } from "@/helper/banner";
  import { formToObject } from "@/helper/codec";
  import { login } from "@/helper/login";
  import { tokenSet } from "@/helper/token";
  import axios from "axios";
  import { push } from "svelte-spa-router";

  let form: HTMLFormElement;
  let error = "";

  const signin = async () => {
    const data = formToObject(form);
    try {
      const response = await login(data);
      tokenSet(response.data.data.token, null);
      // TODO: use pop for history
      push("/");
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        if (reason.response.status == 401) {
          error = reason.response.data.error;
        } else {
          error = reason.message;
        }
      }
    }
  };
</script>

<div
  class="w-full min-h-screen bg-gray-50 flex flex-col items-center pt-6 sm:pt-0"
>
  <div class="w-full sm:max-w-md p-5 mx-auto">
    <h2 class="mb-12 text-center text-sm font-extrabold [line-height:1.2]">
      {banner}
    </h2>
    <form bind:this={form} on:submit|preventDefault|stopPropagation={signin}>
      <div class="mb-4">
        <label class="block mb-1" for="name">Username</label>
        <input
          id="name"
          type="text"
          name="name"
          class="py-2 px-3 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 rounded-md shadow-sm disabled:bg-gray-100 mt-1 block w-full"
        />
      </div>
      <div class="mb-4">
        <label class="block mb-1" for="password">Password</label>
        <input
          id="password"
          type="password"
          name="password"
          class="py-2 px-3 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 rounded-md shadow-sm disabled:bg-gray-100 mt-1 block w-full"
        />
      </div>
      <div class="mt-6">
        <button
          type="submit"
          class="w-full inline-flex items-center justify-center px-4 py-2 bg-red-600 border border-transparent rounded-md font-semibold capitalize text-white hover:bg-red-700 active:bg-red-700 focus:outline-none focus:border-red-700 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
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

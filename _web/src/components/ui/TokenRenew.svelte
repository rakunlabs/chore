<script lang="ts">
  import { login, renew } from "@/helper/login";

  import { tokenGet, tokenSet } from "@/helper/token";
  import { addToast } from "@/store/toast";
  import axios from "axios";
  import { onDestroy, onMount } from "svelte";
  import Icon from "./Icon.svelte";

  let expDate: number;
  let remain = null as string;
  let showLogin = false;
  let pass = "";

  const getExpDate = () => {
    try {
      const [, claims] = tokenGet();
      expDate = claims.exp * 1000;
    } catch (error) {
      addToast(error, "warn");
      expDate = -1;
    }
  };

  let countdown: ReturnType<typeof setTimeout>;

  const listenStoreEvent = () => {
    getExpDate();

    if (remain == null) {
      setCountDown();
      cancelButton();
    }
  };

  const cancelButton = () => {
    showLogin = false;
    pass = "";
  };

  const renewButton = async () => {
    if (remain) {
      // get new token based on old token
      try {
        const response = await renew();

        tokenSet(response.data.data.token);
        getExpDate();

        cancelButton();
      } catch (reason: unknown) {
        let error = "";
        if (axios.isAxiosError(reason)) {
          error = reason.response.data.error ?? reason.message;
        } else {
          error = reason as any;
        }

        addToast(error, "alert");
      }

      return;
    }

    // login with form
    showLogin = true;
  };

  const loginButton = async () => {
    try {
      const [, claims] = tokenGet();

      const response = await login({
        login: claims.user,
        password: pass,
      });

      tokenSet(response.data.data.token);
      getExpDate();
      setCountDown();

      cancelButton();
    } catch (reason: unknown) {
      let error = "";
      if (axios.isAxiosError(reason)) {
        error = reason.response.data.error ?? reason.message;
      } else {
        error = reason as any;
      }

      addToast(error, "alert");
    }
  };

  const setCountDown = () => {
    countdown = setInterval(() => {
      const now = new Date().getTime();
      const distance = expDate - now;

      if (distance < 0) {
        clearInterval(countdown);
        remain = null;
        return;
      }

      const days = Math.floor(distance / (1000 * 60 * 60 * 24));
      const hours = Math.floor(
        (distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
      );
      const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((distance % (1000 * 60)) / 1000);

      let secondsShow: string;

      if (isNaN(seconds)) {
        remain = "infinity";
        return;
      }

      secondsShow = seconds + "s";
      if (seconds < 10) {
        secondsShow = "0" + secondsShow;
      }

      remain =
        (days ? days + "d " : "") +
        (hours ? hours + "h " : "") +
        (minutes ? minutes + "m " : "") +
        secondsShow;
    });
  };

  onMount(() => {
    getExpDate();
    setCountDown();

    window.addEventListener("storage", listenStoreEvent);
  });

  onDestroy(() => {
    window.removeEventListener("storage", listenStoreEvent);
    clearInterval(countdown);
  });
</script>

<div class="text-yellow-200 flex items-center gap-2">
  {#if remain == "infinity"}
    <span>Infinity Token</span>
  {:else}
    {#if remain != null}
      <span>
        Token Expire in {remain}
      </span>
    {:else}
      <span class="text-white bg-red-500 px-2 py-1">Token Expired</span>
    {/if}
    {#if showLogin}
      <input
        type="password"
        bind:value={pass}
        placeholder="password"
        class="py-1 px-2 text-black"
      />
      <button
        on:click={loginButton}
        class="py-1 px-1 bg-transparent border-2 border-green-500 text-sm hover:bg-green-500 fill-white"
        ><Icon icon="ok" height="1.25rem" /></button
      >
      <button
        on:click={cancelButton}
        class="py-1 px-1 bg-transparent border-2 border-red-500 text-sm hover:bg-red-500 fill-white"
        ><Icon icon="close" height="1.25rem" /></button
      >
    {:else}
      <button
        class="px-4 py-1 my-auto bg-transparent border-2 border-yellow-200 text-sm hover:bg-yellow-200 hover:text-black"
        on:click|stopPropagation={renewButton}
      >
        {remain != null ? "Renew Token" : "Login"}
      </button>
    {/if}
  {/if}
</div>

import { writable } from "svelte/store";

const ui = {
  "sidebar": "",
};

export const storeData = writable(ui);

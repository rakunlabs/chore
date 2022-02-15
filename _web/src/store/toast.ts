import { writable } from "svelte/store";
import update from "immutability-helper";

type Types = "info" | "warn" | "alert";

type Toast = {
  message: string
  type: Types
  timeout: ReturnType<typeof setTimeout>
  id: number
}

const toast = [] as Array<Toast>;

export const storeToast = writable(toast);

export const addToast = (message: string, type: Types = "info", timeout = 3000) => {
  storeToast.update((v) => {
    const id = v.length == 0 ? 0 : v[v.length - 1].id + 1;
    return update(v, {
      $push: [{
        message,
        type,
        timeout: timeout > 0 ? setTimeout(()=>removeToast(id), timeout) : null,
        id: id,
      } as Toast],
    });
  });
};

export const removeToast = (id: number) => {
  storeToast.update((v) => {
    if (v[id]?.timeout) {
      clearTimeout(v[id].timeout);
    }

    return update(v, {
      $set: v.filter((vf) => vf.id != id),
    });
  });
};

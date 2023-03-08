import { endpoint } from "./nodes/endpoint";
import { template } from "./nodes/template";
import { request } from "./nodes/request";
import { script } from "./nodes/script";
import { forLoop } from "./nodes/forLoop";
import { ifCase } from "./nodes/ifCase";
import { control } from "./nodes/control";
import { respond } from "./nodes/respond";
import { log } from "./nodes/log";
import { email } from "./nodes/email";
import { note } from "./nodes/note";
import { wait } from "./nodes/wait";

export type node = {
  name: string
  html: string
  data: any
  input: number
  output: number
  optionalInput?: boolean
  class?: string
  style?: string
};

export const nodes = { endpoint, template, request, script, forLoop, ifCase, control, respond, log, email, note, wait } as Record<string, node>;

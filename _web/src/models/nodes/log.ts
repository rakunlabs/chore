import type { node } from "@/models/node";

export type logData = {
  message: string
  level: string
  data: boolean
  tags: string
};

export const log: node = {
  name: "log",
  html: `
  <div>
    <div class="title-box">Log</div>
  </div>
  `,
  data: {
    message: "",
    level: "debug",
    data: false,
    tags: "",
  },
  input: 1,
  output: 1,
  class: "title-box-alone",
};

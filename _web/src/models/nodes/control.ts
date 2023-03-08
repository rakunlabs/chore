import type { node } from "@/models/node";

export type controlData = {
  info: string
  control: string
  endpoint: string
  method: string
  tags: string
};

export const control: node = {
  name: "control",
  html: `
  <div>
    <div class="title-box">Control</div>
    <div class="box">
      <input type="text" placeholder="info" name="info" readonly disabled df-info>
    </div>
  </div>
  `,
  data: {
    info: "",
    control: "",
    endpoint: "",
    method: "POST",
    tags: "",
  },
  input: 1,
  output: 1,
};

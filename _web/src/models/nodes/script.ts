import type { node } from "@/models/node";

export type scriptData = {
  info: string,
  inputs: string,
  script: string
  tags: string
};

export const script: node = {
  name: "script",
  html: `
  <div>
    <div class="title-box">Script</div>
    <div class="box">
      <input type="text" placeholder="info" name="info" readonly disabled df-info>
    </div>
  </div>
  `,
  data: {
    info: "",
    inputs: "",
    script: `function main(data) {
  return data;
}
`,
    tags: "",
  } as scriptData,
  input: 1,
  output: 3,
  optionalInput: true,
  class: "node-script",
};

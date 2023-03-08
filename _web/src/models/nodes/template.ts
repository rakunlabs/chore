import type { node } from "@/models/node";

export type templateData = {
  template: string
  tags: string
};

export const template: node = {
  name: "template",
  html: `
  <div>
    <div class="title-box">Template</div>
    <div class="box">
      <input type="text" placeholder="deepcore/create-issue" name="template" readonly disabled df-template>
    </div>
  </div>
  `,
  data: {
    template: "",
    tags: "",
  } as templateData,
  input: 1,
  output: 1,
};

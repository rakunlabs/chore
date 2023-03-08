import type { node } from "@/models/node";

export type emailData = {
  from: string
  to: string
  cc: string
  bcc: string
  subject: string
  tags: string
};

export const email: node = {
  name: "email",
  html: `
  <div>
    <div class="title-box">Email</div>
  </div>
  `,
  data: {
    from: "",
    to: "",
    cc: "",
    bcc: "",
    subject: "",
    tags: "",
  },
  input: 2,
  output: 0,
  class: "node-email title-box-alone",
};

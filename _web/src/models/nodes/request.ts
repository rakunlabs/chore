import type { node } from "@/models/node";

export type requestData = {
  info: string,
  skip_verify: boolean,
  payload_nil: boolean,
  pool_client: boolean,
  url: string,
  method: string,
  auth: string,
  headers: string,
  retry_codes: string,
  retry_decodes: string,
  tags: string
};

export const request: node = {
  name: "request",
  html: `
  <div>
    <div class="title-box">Request</div>
    <div class="box">
      <input type="text" placeholder="info" name="info" readonly disabled df-info>
    </div>
  </div>
  `,
  data: {
    info: "",
    skip_verify: false,
    payload_nil: false,
    pool_client: false,
    url: "",
    method: "",
    auth: "",
    headers: "",
    retry_codes: "",
    retry_decodes: "",
    tags: "",
  } as requestData,
  input: 2,
  output: 3,
  class: "node-request",
};

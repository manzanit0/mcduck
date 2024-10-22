import { Handlers, PageProps } from "$fresh/server.ts";
import { encodeBase64 } from "$std/encoding/base64.ts";
import ReceiptBreakdown from "../../islands/ReceiptBreakdown.tsx";
import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ReceiptsService } from "../../gen/api/receipts.v1/receipts_connect.ts";
import {
  mapFullReceiptToSerializable,
  SerializableFullReceipt,
} from "../../lib/types.ts";

const url = Deno.env.get("API_HOST") || "";

const transport = createConnectTransport({
  baseUrl: url,
});

const client = createPromiseClient(ReceiptsService, transport);

export const handler: Handlers = {
  async GET(_, ctx) {
    if (!ctx.state || !ctx.state.loggedIn) {
      return ctx.renderNotFound({});
    }

    const res = await client.getReceipt(
      { id: BigInt(ctx.params.id) },
      { headers: { authorization: `Bearer ${ctx.state.authToken}` } }
    );

    const receipt = mapFullReceiptToSerializable(res.receipt!);
    return await ctx.render({ receipt: receipt, file: res.receipt!.file });
  },
};

interface Props {
  receipt: SerializableFullReceipt;
  file: Uint8Array;
}

export default function Single({ data: { receipt, file } }: PageProps<Props>) {
  const encoded = encodeBase64(file);

  let receiptView = (
    <img
      class="object-contain"
      src={`data:image/png;base64, ${encoded}`}
      alt="Receipt image"
    />
  );

  const contentType = detectContentType(file);
  if (contentType === "application/pdf") {
    receiptView = (
      <embed
        class=" overflow-visible"
        src={`data:application/pdf;base64, ${encoded}`}
        type="application/pdf"
        frameborder="0"
        scrolling="auto"
        height="123%"
        width="100%"
      />
    );
  }

  return (
    <div class="m-6 h-screen">
      <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
        Receipt #{receipt.id}
      </h2>
      <div class="mt-10 grid grid-cols-3 gap-4 items-center">
        <div class="col-span-1 h-full">{receiptView}</div>
        <div class="col-span-2 p-5">
          <ReceiptBreakdown receipt={receipt} url={url} />
        </div>
      </div>
    </div>
  );
}

function detectContentType(data: Uint8Array): string {
  const inspected = Deno.inspect(data).replaceAll(" ", "").replaceAll("\n", "");

  if (inspected.includes("[255,216,255")) {
    return "image/jpeg";
  } else if (inspected.includes("[137,80,78,71,13,10,26,10")) {
    return "image/png";
  } else if (inspected.includes("[71,73,70,56")) {
    return "image/gif";
  } else if (inspected.includes("[37,80,68,70")) {
    return "application/pdf";
  } else {
    return "application/octet-stream"; // Default binary data type
  }
}

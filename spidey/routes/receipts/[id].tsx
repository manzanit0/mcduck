import { RouteContext } from "$fresh/server.ts";
import { encodeBase64 } from "$std/encoding/base64.ts";
import ExpensesTable from "../../islands/ExpensesTable.tsx";
import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ReceiptsService } from "../../gen/api/receipts.v1/receipts_connect.ts";
import {
  mapExpensesToSerializable,
  mapReceiptsToSerializable,
} from "../../lib/types.ts";

import { AuthState } from "../../lib/auth.ts";
import ReceiptForm from "../../islands/ReceiptForm.tsx";

const url = Deno.env.get("API_HOST")!;

export default async function Single(_: Request, ctx: RouteContext<AuthState>) {
  console.log("get receipt");
  const transport = createConnectTransport({
    baseUrl: url!,
  });

  const client = createPromiseClient(ReceiptsService, transport);

  const res = await client.getReceipt(
    { id: BigInt(ctx.params.id) },
    {
      headers: { authorization: `Bearer ${ctx.state.authToken}` },
    },
  );

  const receipt = mapReceiptsToSerializable([res.receipt!])[0];
  const encoded = encodeBase64(res.receipt!.file);

  let receiptView = (
    <img
      class="object-contain"
      src={`data:image/png;base64, ${encoded}`}
      alt="Receipt image"
    />
  );

  const contentType = detectContentType(res.receipt!.file);
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
        Receipt #{ctx.params.id}
      </h2>
      <div class="mt-10 grid grid-cols-3 gap-4 items-center">
        <div class="col-span-1 h-full">
          {receiptView}
        </div>
        <div class="col-span-2 p-5">
          <ReceiptForm receipt={receipt} url={url} />
          <div class="mt-10">
            <h2 class="text-base font-semibold leading-7 text-gray-900">
              Expenses
            </h2>
            <div class="mt-2">
              <ExpensesTable
                expenses={mapExpensesToSerializable(res.receipt!.expenses)}
                url={url}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function detectContentType(data: Uint8Array): string {
  const inspected = Deno.inspect(data).replaceAll(" ", "").replaceAll("\n", "");

  if (inspected.includes("[255,216,255")) {
    return "image/jpeg";
  } else if (
    inspected.includes("[137,80,78,71,13,10,26,10")
  ) {
    return "image/png";
  } else if (inspected.includes("[71,73,70,56")) {
    return "image/gif";
  } else if (inspected.includes("[37,80,68,70")) {
    return "application/pdf";
  } else {
    return "application/octet-stream"; // Default binary data type
  }
}

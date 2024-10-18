import { Handlers, PageProps } from "$fresh/server.ts";
import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ReceiptsService } from "../../gen/api/receipts.v1/receipts_connect.ts";
import { ListReceiptsSince } from "../../gen/api/receipts.v1/receipts_pb.ts";
import ReceiptsTable from "../../islands/ReceiptsTable.tsx";
import {
  mapReceiptsToSerializable,
  SerializableReceipt,
} from "../../lib/types.ts";

const url = Deno.env.get("API_HOST")!;
const transport = createConnectTransport({
  baseUrl: url!,
});

const client = createPromiseClient(ReceiptsService, transport);

export const handler: Handlers = {
  async GET(_, ctx) {
    if (!ctx.state || !ctx.state.loggedIn) {
      return ctx.renderNotFound({});
    }

    const res = await client.listReceipts(
      { since: ListReceiptsSince.ALL_TIME },
      { headers: { authorization: `Bearer ${ctx.state.authToken}` } },
    );

    console.log("got", res.receipts.length, "receipts");
    const serializable = mapReceiptsToSerializable(res.receipts);
    return await ctx.render({ receipts: serializable });
  },
  async POST(req, ctx) {
    const form = await req.formData();

    // NOTE: this "receipt-ids" lives in the ReceiptsTable island.
    const receiptIds = (form.get("receipt-ids") as string).split(",");
    console.log(receiptIds)

    const requests = [];
    for (const id of receiptIds) {
      const r = client.deleteReceipt(
        { id: BigInt(id) },
        { headers: { authorization: `Bearer ${ctx.state.authToken}` } },
      );

      requests.push(r);
    }

    try {
      await Promise.all(requests);
    } catch (err) {
      console.log(err);
    }

    const res = await client.listReceipts(
      { since: ListReceiptsSince.ALL_TIME },
      { headers: { authorization: `Bearer ${ctx.state.authToken}` } },
    );

    console.log("got", res.receipts.length, "receipts");
    const serializable = mapReceiptsToSerializable(res.receipts);
    return await ctx.render({ receipts: serializable });
  },
};

interface Props {
  receipts: SerializableReceipt[];
}

export default function List({ data: { receipts } }: PageProps<Props>) {
  return (
    <div class="m-6">
      <ReceiptsTable
        receipts={receipts}
        url={url}
      />
    </div>
  );
}

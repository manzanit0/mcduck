import { createConnectTransport } from "@connectrpc/connect-web";
import { createPromiseClient } from "@connectrpc/connect";
import { ReceiptsService } from "../gen/api/receipts.v1/receipts_connect.ts";
import { UpdateReceiptRequest } from "../gen/api/receipts.v1/receipts_pb.ts";
import { getAuthTokenFromBrowser } from "./auth.ts";
import { PartialMessage } from "@bufbuild/protobuf";
import { CreateExpenseRequest, DeleteExpenseRequest, UpdateExpenseRequest } from "../gen/api/expenses.v1/expenses_pb.ts";
import { ExpensesService } from "../gen/api/expenses.v1/expenses_connect.ts";

export function updateReceipt(host: string, body: PartialMessage<UpdateReceiptRequest>) {
  const client = createPromiseClient(
    ReceiptsService,
    createConnectTransport({
      baseUrl: host,
    }),
  );

  const { authToken } = getAuthTokenFromBrowser();

  return client.updateReceipt(body, {
    headers: { authorization: `Bearer ${authToken}` },
  });
}

export function createExpense(host: string, body: PartialMessage<CreateExpenseRequest>) {
  const client = createPromiseClient(
    ExpensesService,
    createConnectTransport({
      baseUrl: host,
    }),
  );

  const { authToken } = getAuthTokenFromBrowser();

  return client.createExpense(body, {
    headers: { authorization: `Bearer ${authToken}` },
  });
}

export function updateExpense(host: string, body: PartialMessage<UpdateExpenseRequest>) {
  const client = createPromiseClient(
    ExpensesService,
    createConnectTransport({
      baseUrl: host,
    }),
  );

  const { authToken } = getAuthTokenFromBrowser();

  return client.updateExpense(body, {
    headers: { authorization: `Bearer ${authToken}` },
  });
}

export function deleteExpense(host: string, body: PartialMessage<DeleteExpenseRequest>) {
  const client = createPromiseClient(
    ExpensesService,
    createConnectTransport({
      baseUrl: host,
    }),
  );

  const { authToken } = getAuthTokenFromBrowser();

  return client.deleteExpense(body, {
    headers: { authorization: `Bearer ${authToken}` },
  });
}
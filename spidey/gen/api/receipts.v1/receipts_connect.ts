// @generated by protoc-gen-connect-es v1.5.0 with parameter "target=ts,import_extension=.ts"
// @generated from file api/receipts.v1/receipts.proto (package api.receipts.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { CreateReceiptsRequest, CreateReceiptsResponse, DeleteReceiptRequest, DeleteReceiptResponse, GetReceiptRequest, GetReceiptResponse, ListReceiptsRequest, ListReceiptsResponse, UpdateReceiptRequest, UpdateReceiptResponse } from "./receipts_pb.ts";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service api.receipts.v1.ReceiptsService
 */
export const ReceiptsService = {
  typeName: "api.receipts.v1.ReceiptsService",
  methods: {
    /**
     * @generated from rpc api.receipts.v1.ReceiptsService.CreateReceipts
     */
    createReceipts: {
      name: "CreateReceipts",
      I: CreateReceiptsRequest,
      O: CreateReceiptsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.receipts.v1.ReceiptsService.UpdateReceipt
     */
    updateReceipt: {
      name: "UpdateReceipt",
      I: UpdateReceiptRequest,
      O: UpdateReceiptResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.receipts.v1.ReceiptsService.DeleteReceipt
     */
    deleteReceipt: {
      name: "DeleteReceipt",
      I: DeleteReceiptRequest,
      O: DeleteReceiptResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.receipts.v1.ReceiptsService.ListReceipts
     */
    listReceipts: {
      name: "ListReceipts",
      I: ListReceiptsRequest,
      O: ListReceiptsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.receipts.v1.ReceiptsService.GetReceipt
     */
    getReceipt: {
      name: "GetReceipt",
      I: GetReceiptRequest,
      O: GetReceiptResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

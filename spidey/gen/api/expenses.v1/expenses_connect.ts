// @generated by protoc-gen-connect-es v1.5.0 with parameter "target=ts,import_extension=.ts"
// @generated from file api/expenses.v1/expenses.proto (package api.expenses.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { CreateExpenseRequest, CreateExpenseResponse, DeleteExpenseRequest, DeleteExpenseResponse, ListExpensesRequest, ListExpensesResponse, UpdateExpenseRequest, UpdateExpenseResponse } from "./expenses_pb.ts";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service api.expenses.v1.ExpensesService
 */
export const ExpensesService = {
  typeName: "api.expenses.v1.ExpensesService",
  methods: {
    /**
     * @generated from rpc api.expenses.v1.ExpensesService.CreateExpense
     */
    createExpense: {
      name: "CreateExpense",
      I: CreateExpenseRequest,
      O: CreateExpenseResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.expenses.v1.ExpensesService.UpdateExpense
     */
    updateExpense: {
      name: "UpdateExpense",
      I: UpdateExpenseRequest,
      O: UpdateExpenseResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.expenses.v1.ExpensesService.DeleteExpense
     */
    deleteExpense: {
      name: "DeleteExpense",
      I: DeleteExpenseRequest,
      O: DeleteExpenseResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.expenses.v1.ExpensesService.ListExpenses
     */
    listExpenses: {
      name: "ListExpenses",
      I: ListExpensesRequest,
      O: ListExpensesResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;


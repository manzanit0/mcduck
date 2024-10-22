import { Timestamp } from "@bufbuild/protobuf";
import {
  Expense,
  FullReceipt,
  Receipt,
} from "../gen/api/receipts.v1/receipts_pb.ts";

export interface SerializableReceipt {
  id: bigint;
  status: number;
  vendor: string;
  date?: string;
  totalAmount: bigint;
  // expenses: SerializableExpense[];
}

export interface SerializableFullReceipt {
  id: bigint;
  status: number;
  vendor: string;
  date?: string;
  totalAmount: bigint;
  expenses: SerializableExpense[];
}

export interface SerializableExpense {
  id: bigint;
  date?: string;
  category: string;
  subcategory: string;
  description: string;
  amount: bigint;
}

export function mapReceiptsToSerializable(
  receipts: Receipt[],
): SerializableReceipt[] {
  return receipts.map((r) => {
    return {
      id: r.id,
      status: r.status,
      vendor: r.vendor,
      date: r.date?.toDate().toISOString(),
      totalAmount: r.totalAmount,
    };
  });
}

export function mapFullReceiptToSerializable(
  r: FullReceipt,
): SerializableFullReceipt {
  return {
    id: r.id,
    status: r.status,
    vendor: r.vendor,
    date: r.date?.toDate().toISOString(),
    totalAmount: r.expenses.reduce((acc, ex) => (acc += ex.amount), 0n),
    expenses: mapExpensesToSerializable(r.expenses),
  };
}

export function mapExpensesToSerializable(
  expenses: Expense[],
): SerializableExpense[] {
  return expenses.map((e) => {
    return {
      id: e.id,
      date: e.date?.toDate().toISOString(),
      category: e.category,
      subcategory: e.subcategory,
      description: e.description,
      amount: e.amount,
    };
  });
}

export function toStringDate(date: Timestamp): string {
  return date.toDate().toISOString()
}

export function toTimestamp(date: string): Timestamp {
  return Timestamp.fromDate(new Date(date))
}
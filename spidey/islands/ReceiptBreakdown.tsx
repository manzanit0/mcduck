import { useSignal, useComputed } from "@preact/signals";
import ExpensesTable from "./ExpensesTable.tsx";
import ReceiptForm from "./ReceiptForm.tsx";
import { SerializableFullReceipt } from "../lib/types.ts";

interface CounterProps {
  receipt: SerializableFullReceipt;
  url: string;
}

export default function ReceiptBreakdown({ receipt, url }: CounterProps) {
  const expensesSignal = useSignal(receipt.expenses.map((x) => useSignal(x)));

  const totalAmount = useComputed(() => {
    const a = expensesSignal.value.reduce(
      (acc, ex) => (acc += ex.value.amount),
      0n
    );
    return Number(a);
  });

  return (
    <div>
      <ReceiptForm receipt={receipt} url={url} totalAmount={totalAmount} />
      <div class="mt-10">
        <h2 class="text-base font-semibold leading-7 text-gray-900">
          Expenses
        </h2>
        <div class="mt-2">
          <ExpensesTable
            expenses={expensesSignal}
            url={url}
            receiptId={receipt.id}
            receiptDate={receipt.date!}
          />
        </div>
      </div>
    </div>
  );
}

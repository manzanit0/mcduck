import Checkbox from "../components/Checkbox.tsx";
import FormattedMoney from "../components/FormattedMoney.tsx";
import GenericTable from "../components/GenericTable.tsx";
import TextInput from "../components/TextInput.tsx";
import { createExpense } from "../lib/receipts.ts";
import { SerializableExpense, toTimestamp } from "../lib/types.ts";
import { useSignal } from "@preact/signals";

interface TableProps {
  receiptId: bigint;
  receiptDate: string;
  expenses: SerializableExpense[];
  url: string;
}

interface CheckeableExpense extends SerializableExpense {
  checked: boolean;
}

export default function ExpensesTable(props: TableProps) {
  const mapped = useSignal(
    props.expenses.map((x) => {
      return useSignal<CheckeableExpense>({
        ...x,
        checked: false,
      });
    })
  );

  const globallySelected = useSignal(false);

  const checkExpenses = () => {
    globallySelected.value = !globallySelected.value;

    for (const r of mapped.value) {
      r.value.checked = globallySelected.value;
    }
  };

  const addExpense = async () => {
    await createExpense(props.url, {
      receiptId: props.receiptId,
      amount: BigInt(0),
      date: toTimestamp(props.receiptDate),
    });

    // NOTE: Hooks (useSignal) can't be used outside of preact components, so
    // I'm unsure what would be the right solution/pattern here.
    location.reload();

    // const expense = {
    //   id: res.expense!.id,
    //   amount: BigInt(0),
    //   category: "",
    //   subcategory: "",
    //   description: "",
    //   date: props.receiptDate,
    //   checked: false,
    // };

    // mapped.value = [...mapped.value, useSignal(expense)];
  };

  return (
    <div class="sm:rounded-lg">
      <GenericTable
        data={mapped.value}
        columns={[
          {
            header: (
              <Checkbox
                onInput={checkExpenses}
                checked={globallySelected.value}
              />
            ),
            accessor: (r) => (
              <Checkbox
                checked={r.value.checked}
                onInput={() => (r.value.checked = !r.value.checked)}
              />
            ),
          },
          {
            header: <span>Category</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.category}
                onfocusout={() => Promise.resolve()}
              />
            ),
          },
          {
            header: <span>Subcategory</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.subcategory}
                onfocusout={() => Promise.resolve()}
              />
            ),
          },
          {
            header: <span>Description</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.description}
                onfocusout={() => Promise.resolve()}
              />
            ),
          },
          {
            header: <span>Amount</span>,
            accessor: (r) => (
              <FormattedMoney currency="EUR" amount={Number(r.value.amount)} />
            ),
          },
          {
            header: <span>Action</span>,
            accessor: (r) => (
              <a
                href={`receipts/${r.value.id}`}
                class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
              >
                Delete
              </a>
            ),
          },
        ]}
      />

      <div class="pt-3 float-right">
        <button
          type="submit"
          class="flex flex-column rounded-md bg-gray-800 px-8 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-gray-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600"
          onClick={addExpense}
        >
          Add Expense
        </button>
      </div>
    </div>
  );
}

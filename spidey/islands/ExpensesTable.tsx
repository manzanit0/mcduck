import { JSX } from "preact/jsx-runtime";
import GenericTable from "../components/GenericTable.tsx";
import TextInput from "../components/TextInput.tsx";
import { createExpense, updateExpense } from "../lib/receipts.ts";
import { SerializableExpense, toCents, toTimestamp } from "../lib/types.ts";
import { Signal, useSignal } from "@preact/signals";
import MoneyInput from "../components/MoneyInput.tsx";

interface TableProps {
  receiptId: bigint;
  receiptDate: string;
  expenses: SerializableExpense[];
  url: string;
}

export default function ExpensesTable(props: TableProps) {
  const mapped = useSignal(
    props.expenses.map((x) => {
      return useSignal(x);
    })
  );

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

  const updateCategory = async (
    e: JSX.TargetedEvent<HTMLInputElement>,
    r: Signal<SerializableExpense>
  ) => {
    if (!e.currentTarget || e.currentTarget.value === "") {
      return;
    }

    const value = e.currentTarget.value;
    r.value = { ...r.value, category: value };
    await updateExpense(props.url, { id: r.value.id, category: value });

    console.log("updated category to", value);
  };

  const updateSubcategory = async (
    e: JSX.TargetedEvent<HTMLInputElement>,
    r: Signal<SerializableExpense>
  ) => {
    if (!e.currentTarget || e.currentTarget.value === "") {
      return;
    }

    const value = e.currentTarget.value;
    r.value = { ...r.value, subcategory: value };
    await updateExpense(props.url, { id: r.value.id, subcategory: value });

    console.log("updated subcategory to", value);
  };

  const updateDescription = async (
    e: JSX.TargetedEvent<HTMLInputElement>,
    r: Signal<SerializableExpense>
  ) => {
    if (!e.currentTarget || e.currentTarget.value === "") {
      return;
    }

    const value = e.currentTarget.value;
    r.value = { ...r.value, description: value };
    await updateExpense(props.url, { id: r.value.id, description: value });

    console.log("updated description to", value);
  };

  const updateAmount = async (
    e: JSX.TargetedEvent<HTMLInputElement>,
    r: Signal<SerializableExpense>
  ) => {
    if (!e.currentTarget || e.currentTarget.value === "") {
      return;
    }

    const value = BigInt(toCents(e.currentTarget.valueAsNumber));
    r.value = { ...r.value, amount: value };
    await updateExpense(props.url, { id: r.value.id, amount: r.value.amount });

    console.log("updated amount to", value);
  };

  return (
    <div class="sm:rounded-lg">
      <GenericTable
        data={mapped.value}
        columns={[
          {
            header: <span>Category</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.category}
                onfocusout={(e) => updateCategory(e, r)}
              />
            ),
          },
          {
            header: <span>Subcategory</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.subcategory}
                onfocusout={(e) => updateSubcategory(e, r)}
              />
            ),
          },
          {
            header: <span>Description</span>,
            accessor: (r) => (
              <TextInput
                value={r.value.description}
                onfocusout={(e) => updateDescription(e, r)}
              />
            ),
          },
          {
            header: <span>Amount</span>,
            accessor: (r) => (
              <MoneyInput
                currency="EUR"
                amount={Number(r.value.amount) / 100}
                onfocusout={(e) => updateAmount(e, r)}
              />
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

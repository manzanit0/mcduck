import { useSignal } from "@preact/signals";
import { JSX } from "preact";

type Props = {
  amount: number;
  currency: string;
  onfocusout: (e: JSX.TargetedEvent<HTMLInputElement>) => Promise<void>;
};

export default function MoneyInput({ amount, currency, onfocusout }: Props) {
  const focusRingClass = useSignal("focus:ring-gray-600");
  const actualAmount = useSignal(`${amount}`);

  return (
    <div>
      <div class="relative mt-2 rounded-md shadow-sm">
        <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
          <span class="text-gray-500 sm:text-sm">{currency}</span>
        </div>
        <input
          type="number"
          name="price"
          id="price"
          lang="de-DE"
          class={`block w-full rounded-md border-0 py-1.5 pl-12 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset ${focusRingClass.value} sm:text-sm sm:leading-6 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none`}
          placeholder="0.00"
          value={actualAmount.value}
          onFocus={() => {
            const value = actualAmount.peek();
            const valid = /^(\d{1,})(\.\d{1,2})?$/.test(value);
            if (!valid) {
              focusRingClass.value = "focus:ring-red-600";
            } else {
              focusRingClass.value = "focus:ring-gray-600";
            }
          }}
          onInput={(e) => {
            const value = e.currentTarget.value;
            console.log("onInput ", value);
            const valid = /^(\d{1,})(\.\d{1,2})?$/.test(value);
            if (!valid) {
              focusRingClass.value = "focus:ring-red-600";
            } else {
              focusRingClass.value = "focus:ring-gray-600";
            }

            actualAmount.value = value;
          }}
          onfocusout={(e) => {
            const valid = /^(\d{1,})(\.\d{1,2})?$/.test(e.currentTarget.value);
            if (valid) {
              onfocusout(e);
            } else {
              focusRingClass.value = "ring-red-600";
            }
          }}
          onKeyDown={(e) => {
            if (
              e.key === "Backspace" ||
              e.key === "ArrowLeft" ||
              e.key === "ArrowRight"
            ) {
              return;
            }

            // NOTE: this doesn't really validate the input to it's full length.
            // The user can still type inputs such as "9,9,90,12" or "9,999999"
            const allowed = e.key.length === 1 && /[0-9,]/.test(e.key);
            if (!allowed) {
              e.preventDefault();
              return;
            }
          }}
        />
      </div>
    </div>
  );
}

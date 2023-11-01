import { LitElement, unsafeCSS } from "lit";

import style from "./tailwind.global.css";

// this is something I found online for using tailwind
const tailwindElement = unsafeCSS(style);

export const TailwindElement = (style: any) =>
  class extends LitElement {
    static styles = [tailwindElement, unsafeCSS(style)];
  };

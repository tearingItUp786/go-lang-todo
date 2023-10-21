import { html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { TailwindElement } from "../shared/tailwind.element";

@customElement("my-button")
export class MyButton extends TailwindElement("") {
  @property({ type: Boolean, reflect: true })
  disabled = false;

  @property({ type: Function, reflect: true })
  onClick: Function = () => {};

  render() {
    return html`<button
      @click=${this.disabled ? null : this.onClick}
      ?disabled=${this.disabled}
      class="
      disabled:text-white
      disabled:cursor-not-allowed	disabled:bg-slate-500 bg-white text-gray-700 hover:bg-gray-700 hover:text-white font-bold py-2 px-4 rounded"
    >
      <slot></slot>
    </button>`;
  }
}

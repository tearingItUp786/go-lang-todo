import { html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { TailwindElement } from "../shared/tailwind.element";

@customElement("my-button")
export class MyButton extends TailwindElement("") {
  @property({ type: String, reflect: true })
  type = "";

  @property({ type: Boolean, reflect: true })
  disabled = false;

  @property({ type: Function, reflect: true })
  onClick: Function = () => {};

  handleClick(event: any) {
    if (this.type === "submit") {
      // Fire the default submit behavior
      const form = this.closest("form");
      if (form) {
        form.dispatchEvent(
          new Event("submit", { bubbles: true, cancelable: true }),
        );
      }
    } else {
      // Handle the click event for non-submit buttons
      // You can add your custom logic here
      this.onClick(event);
    }
  }

  render() {
    return html`<button
      type=${this.type}
      @click=${this.disabled ? null : this.handleClick}
      ?disabled=${this.disabled}
      class="
      disabled:text-white
      disabled:cursor-not-allowed	disabled:bg-slate-200 bg-white text-gray-700 hover:bg-gray-700 hover:text-white font-bold py-2 px-4 rounded"
    >
      <slot></slot>
    </button>`;
  }
}

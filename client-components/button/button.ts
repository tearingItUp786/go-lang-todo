import { html } from "lit";
import { customElement } from "lit/decorators.js";
import { TailwindElement } from "../shared/tailwind.element";

@customElement("my-button")
export class MyButton extends TailwindElement("") {
  onClick() {
    // Get the current query parameters
    var queryParams = new URLSearchParams(window.location.search);
    // Update the isOpen parameter to true
    queryParams.set("isOpen", "true");
    // Update the URL with the modified query parameters
    window.history.pushState({}, "", "?" + queryParams.toString());

    var event = new Event("historystatechanged");
    window.dispatchEvent(event);
  }

  render() {
    return html`<button
      id="modal-button"
      class="bg-white text-gray-700 hover:bg-gray-700 hover:text-white font-bold py-2 px-4 rounded"
      @click=${this.onClick}
    >
      Bulk Edit
    </button>`;
  }
}

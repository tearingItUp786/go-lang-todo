import { html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { TailwindElement } from "../shared/tailwind.element";
import "../button/button";

@customElement("my-open-modal-button")
export class MyOpenModalButton extends TailwindElement("") {
  @property({ type: String, reflect: true })
  modalId = "";

  ownOnClick() {
    // Get the current query parameters
    var queryParams = new URLSearchParams(window.location.search);
    // Update the isOpen parameter to true
    queryParams.set(this.modalId, "true");
    // Update the URL with the modified query parameters
    window.history.pushState({}, "", "?" + queryParams.toString());

    var event = new Event("historystatechanged");
    window.dispatchEvent(event);
  }

  render() {
    return html`
      <my-button id="${this.modalId}--button" @click=${this.ownOnClick}>
        <slot></slot>
      </my-button>
    `;
  }
}

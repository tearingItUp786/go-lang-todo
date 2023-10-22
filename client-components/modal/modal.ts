import { html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { TailwindElement } from "../shared/tailwind.element";
import styles from "./modal.css?inline";

@customElement("my-modal")
export class MyModal extends TailwindElement(styles) {
  @property({ type: String, reflect: true })
  title: string = "";

  @property({ type: Boolean, reflect: true })
  isOpen: boolean = false;

  // give me an onClose function
  @property({ type: Function, reflect: true })
  onClose: Function | undefined;

  constructor() {
    super();
    this.handleQueryParamChange();
  }

  connectedCallback() {
    super.connectedCallback();
    document.addEventListener("keyup", this.handleEscapeKey.bind(this));
    window.addEventListener("historystatechanged", this.handleQueryParamChange);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    document.removeEventListener("keyup", this.handleEscapeKey.bind(this));

    window.removeEventListener(
      "historystatechanged",
      this.handleQueryParamChange,
    );
  }

  handleQueryParamChange = () => {
    const searchParams = new URLSearchParams(window.location.search);
    this.isOpen = searchParams.get(this.id) === "true";
    this.requestUpdate();
  };

  createRenderRoot() {
    const root = super.createRenderRoot();

    root.addEventListener("click", (e: Event) => {
      if ((e.target as Element).id !== "container") return;
      window.history.replaceState({}, "", window.location.pathname);
      this.isOpen = false;
    });

    return root;
  }

  handleEscapeKey(event: any) {
    // close if user presses escape key
    if (event.key === "Escape") {
      window.history.replaceState({}, "", window.location.pathname);
      this.isOpen = false;
    }
  }

  render() {
    if (!this.isOpen) return html``;

    return html`
      <div
        class="animate-fade-in-long left-1/2 top-[35%] transform -translate-x-1/2 -translate-y-1/2 z-50 absolute bg-white rounded-lg shadow-lg p-6 min-w-[50%]"
      >
        <h2 class="text-xl font-bold mb-4">${this.title}</h2>

        <slot name="content"></slot>

        <div class="flex justify-end">
          <slot
            @click="${() => {
              window.history.replaceState({}, "", window.location.pathname);
              this.isOpen = !this.isOpen;
              this?.onClose?.();
            }}"
            name="action"
          ></slot>
          <button
            @click="${() => {
              window.history.replaceState({}, "", window.location.pathname);
              this.isOpen = !this.isOpen;
              this?.onClose?.();
            }}"
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          >
            Close
          </button>
        </div>
      </div>
      <!-- Overlay -->
      <div id="container" class="fixed inset-0 bg-black opacity-50 z-40"></div>
    `;
  }
}

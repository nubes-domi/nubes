import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  connect() {
    console.log("hi")
  }
  
  close() {
    this.element.parentElement.removeAttribute("src")
    this.element.remove()
  }
}

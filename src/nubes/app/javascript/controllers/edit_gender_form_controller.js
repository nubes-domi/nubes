import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "select", "other", "otherField", "pronouns" ]

  connect() {
    this.update()
  }

  update() {
    let value = this.selectTarget.value

    if (value != "other") {
      this.selectTarget.name = "user[gender]"
      this.otherFieldTarget.name = ""
      this.otherFieldTarget.required = false
      this.otherTarget.style.display = "none"
      this.pronounsTarget.style.display = "none"
    } else {
      this.selectTarget.name = ""
      this.otherFieldTarget.name = "user[gender]"
      this.otherFieldTarget.required = true
      this.otherTarget.style.display = "block"
      this.pronounsTarget.style.display = "block"
    }
  }
}

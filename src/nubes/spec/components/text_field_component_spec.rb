require "rails_helper"

RSpec.describe TextFieldComponent, type: :component do
  it "renders component with name attribute" do
    render_inline(described_class.new(name: "username"))

    expect(rendered_component).to have_css "input[name='username']"
  end

  it "uses name attribute as id when otherwise unspecified" do
    render_inline(described_class.new(name: "username"))

    expect(rendered_component).to have_css "input#username"
  end

  it "prefers specific id" do
    render_inline(described_class.new(name: "username", id: "my-username"))

    expect(rendered_component).to have_css "input#my-username[name='username']"
  end

  it "shows errors" do
    render_inline(described_class.new(name: "username", error: "nice try"))

    expect(rendered_component).to have_text "nice try"
  end

  it "supports sizes" do
    render_inline(described_class.new(name: "username", size: "lg"))

    expect(rendered_component).to have_css ".text-control--lg"
  end

  it "supports custom attributes" do
    render_inline(described_class.new(name: "username", "data-magic" => "alohomora"))

    expect(rendered_component).to have_css "input[data-magic='alohomora']"
  end

  it "mixes classes in" do
    render_inline(described_class.new(name: "username", klass: "meh"))

    expect(rendered_component).to have_css ".text-control.meh"
  end
end

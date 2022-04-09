require "rails_helper"

RSpec.describe Contacts::Create do
  fixtures :users

  let(:arguments) do
    {
      user: users(:joe),
      attributes: {
        given_name: "Alice"
      }
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "creates a new contact" do
      expect { subject.call(**arguments) }.to(change { Contact.count }.by(1))
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:jane))).to be_failure
    end
  end
end

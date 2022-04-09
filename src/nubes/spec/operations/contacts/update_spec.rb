require "rails_helper"

RSpec.describe Contacts::Update do
  fixtures :contacts, :users

  let(:arguments) do
    {
      user: users(:joe),
      id: contacts(:joe_contact1).id,
      attributes: {
        family_name: "Doe"
      }
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "updates the contact" do
      contact = subject.call(**arguments).value!
      expect(contact.family_name).to eq "Doe"
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:jane))).to be_failure
    end
  end
end

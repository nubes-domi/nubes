require "rails_helper"

RSpec.describe Contacts::AddAddress do
  fixtures :contacts, :users

  let(:contact) { contacts(:joe_contact1) }

  let(:arguments) do
    {
      user: users(:joe),
      contact_id: contact.id,
      attributes: {
        name: "Work",
        kind: "PHONE",
        data: "+447597777777"
      }
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "creates a new contact" do
      expect { subject.call(**arguments) }.to(change { contact.addresses.count }.by(1))
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:jane))).to be_failure
    end
  end
end

require "rails_helper"

RSpec.describe Contacts::DestroyAddress do
  fixtures :contacts, :contact_addresses, :users

  let(:contact) { contacts(:joe_contact1) }

  let(:arguments) do
    {
      user: users(:joe),
      id: contact_addresses(:joe_contact1_addr).id,
      contact_id: contact.id
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "destroys the address" do
      expect { subject.call(**arguments) }.to(change { contact.addresses.count }.by(-1))
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:jane))).to be_failure
    end
  end
end

require "rails_helper"

RSpec.describe Contacts::UpdatePostalAddress do
  fixtures :contacts, :contact_postal_addresses, :users

  let(:contact) { contacts(:joe_contact1) }

  let(:arguments) do
    {
      user: users(:joe),
      id: contact_postal_addresses(:joe_contact1_postal_addr).id,
      contact_id: contact.id,
      attributes: {
        country: "ua"
      }
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "updates the address" do
      subject.call(**arguments)

      address = contact.postal_addresses.first
      expect(address.country).to eq "ua"
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:jane))).to be_failure
    end
  end
end

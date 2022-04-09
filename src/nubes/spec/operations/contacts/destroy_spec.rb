require "rails_helper"

RSpec.describe Contacts::Destroy do
  fixtures :contacts, :users

  let(:arguments) do
    {
      user: users(:joe),
      id: contacts(:joe_contact1).id
    }
  end

  context "with valid arguments" do
    it "succeeds" do
      expect(subject.call(**arguments)).to be_success
    end

    it "deletes a contact" do
      expect { subject.call(**arguments) }.to(change { Contact.count }.by(-1))
    end
  end

  context "with missing required arguments" do
    it "fails" do
      expect(subject.call(user: users(:joe))).to be_failure
    end
  end

  context "when the user is not the owner" do
    it "fails" do
      expect(subject.call(**arguments, user: users(:jane))).to be_failure
    end
  end
end

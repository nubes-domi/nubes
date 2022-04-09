require "rails_helper"

RSpec.describe "GraphQL" do
  fixtures :users, :contacts

  def run(variables = nil)
    NubesSchema.execute(mutation, variables:, context: { current_user: users(:joe) })
  end

  describe "Contacts" do
    describe "mutations" do
      describe "addContact" do
        let(:mutation) do
          <<-GRAPHQL
            mutation {
              addContact(input: { givenName: "Arthur" }) {
                contact {
                  id
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        it "creates a new contact" do
          expect { run }.to change { Contact.count }.by(1)
        end

        it "returns the new contact id" do
          expect(run.dig("data", "addContact", "contact", "id")).not_to be_nil
        end
      end

      describe "addContactPostalAddress" do
        let(:mutation) do
          <<-GRAPHQL
            mutation($contactID: ID!) {
              addContactAddress(input: { contactId: $contactID, kind: "phone", data: "+447777777777" }) {
                contactAddress {
                  id
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) { { "contactID" => contacts(:joe_contact1).id } }

        it "creates a new contact" do
          expect { run(variables) }.to change { ContactAddress.count }.by(1)
        end

        it "returns the new contact id" do
          expect(run(variables).dig("data", "addContactAddress", "contactAddress", "id")).not_to be_nil
        end
      end

      describe "addContactPostalAddress" do
        let(:mutation) do
          <<-GRAPHQL
            mutation($contactID: ID!) {
              addContactPostalAddress(input: { contactId: $contactID, name: "home", postalCode: "E14 4AF" }) {
                contactPostalAddress {
                  id, postalCode
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) { { "contactID" => contacts(:joe_contact1).id } }

        it "creates a new contact postal address" do
          expect { run(variables) }.to change { ContactPostalAddress.count }.by(1)
        end

        it "returns the new contact postal address id" do
          expect(run(variables).dig("data", "addContactPostalAddress", "contactPostalAddress",
                                    "postalCode")).to eq "E14 4AF"
        end
      end

      describe "destroyContact" do
        let(:mutation) do
          <<-GRAPHQL
            mutation($id: ID!) {
              destroyContact(input: { id: $id }) {
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) { { "id" => contacts(:joe_contact1).id } }

        it "destroys the record" do
          expect { run(variables) }.to(change { Contact.count }.by(-1))
        end
      end

      describe "destroyContactAddress" do
        fixtures :contact_addresses

        let(:mutation) do
          <<-GRAPHQL
            mutation($id: ID!, $contactId: ID!) {
              destroyContactAddress(input: { id: $id, contactId: $contactId }) {
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) do
          {
            "id" => contact_addresses(:joe_contact1_addr).id,
            "contactId" => contacts(:joe_contact1).id
          }
        end

        it "destroys the record" do
          expect { run(variables) }.to(change { ContactAddress.count }.by(-1))
        end
      end

      describe "destroyContactPostalAddress" do
        fixtures :contact_postal_addresses

        let(:mutation) do
          <<-GRAPHQL
            mutation($id: ID!, $contactId: ID!) {
              destroyContactPostalAddress(input: { id: $id, contactId: $contactId }) {
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) do
          {
            "id" => contact_postal_addresses(:joe_contact1_postal_addr).id,
            "contactId" => contacts(:joe_contact1).id
          }
        end

        it "destroys the record" do
          expect { run(variables) }.to(change { ContactPostalAddress.count }.by(-1))
        end
      end

      describe "updateContact" do
        let(:mutation) do
          <<-GRAPHQL
            mutation($contactId: ID!) {
              updateContact(input: { id: $contactId, givenName: "Johanna" }) {
                contact {
                  givenName
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) { { "contactId" => contacts(:joe_contact1).id } }

        it "updates the record" do
          run(variables)
          expect(contacts(:joe_contact1).reload.given_name).to eq "Johanna"
        end

        it "returns the new contact details" do
          expect(run(variables).dig("data", "updateContact", "contact", "givenName")).to eq "Johanna"
        end
      end

      describe "updateContactPostalAddress" do
        fixtures :contact_postal_addresses

        let(:mutation) do
          <<-GRAPHQL
            mutation($contactId: ID!, $contactPostalAddressId: ID!) {
              updateContactPostalAddress(input: { id: $contactPostalAddressId, contactId: $contactId, name: "wooork" }) {
                contactPostalAddress {
                  name
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) do
          { "contactId" => contacts(:joe_contact1).id,
            "contactPostalAddressId" => contact_postal_addresses(:joe_contact1_postal_addr).id }
        end

        it "updates the record" do
          run(variables)
          expect(contact_postal_addresses(:joe_contact1_postal_addr).reload.name).to eq "wooork"
        end

        it "returns the new contact details" do
          expect(run(variables).dig("data", "updateContactPostalAddress", "contactPostalAddress",
                                    "name")).to eq "wooork"
        end
      end

      describe "updateContactAddress" do
        fixtures :contact_addresses

        let(:mutation) do
          <<-GRAPHQL
            mutation($contactId: ID!, $contactAddressId: ID!) {
              updateContactAddress(input: { id: $contactAddressId, contactId: $contactId, name: "wooork" }) {
                contactAddress {
                  name
                }
                errors {
                  path
                  message
                }
              }
            }
          GRAPHQL
        end

        let(:variables) do
          { "contactId" => contacts(:joe_contact1).id, "contactAddressId" => contact_addresses(:joe_contact1_addr).id }
        end

        it "updates the record" do
          run(variables)
          expect(contact_addresses(:joe_contact1_addr).reload.name).to eq "wooork"
        end

        it "returns the new contact details" do
          expect(run(variables).dig("data", "updateContactAddress", "contactAddress", "name")).to eq "wooork"
        end
      end
    end
  end
end

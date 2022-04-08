module Mutations
  class CreateContactPostalAddress < Mutations::BaseMutation
    null true

    ContactPostalAddress.define_graphql_mutation(self, Types::Models::ContactPostalAddressType)

    def resolve(contact_id:, **params)
      contact_postal_address = context[:current_user].contacts.find(contact_id).postal_addresses.build(params)
      if contact_postal_address.save
        # Successful creation, return the created object with no errors
        {
          contact_postal_address:
        }
      else
        # Failed save, return the errors to the client
        {
          errors: contact_postal_address.errors.full_messages
        }
      end
    end
  end
end

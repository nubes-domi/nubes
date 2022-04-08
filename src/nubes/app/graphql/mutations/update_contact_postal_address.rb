module Mutations
  class UpdateContactPostalAddress < Mutations::BaseMutation
    null true

    ContactPostalAddress.define_graphql_mutation(self, Types::Models::ContactPostalAddressType, type: :update)

    def resolve(id:, **params)
      contact_postal_address = context[:current_user].contact_postal_addresses.find(id)
      if contact_postal_address.update(params)
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

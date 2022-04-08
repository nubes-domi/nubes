module Mutations
  class UpdateContactAddress < Mutations::BaseMutation
    null true

    ContactAddress.define_graphql_mutation(self, Types::Models::ContactAddressType, type: :update)

    def resolve(id:, **params)
      contact_address = context[:current_user].addresses.find(id)
      if contact_address.update(params)
        # Successful creation, return the created object with no errors
        {
          contact_address:
        }
      else
        # Failed save, return the errors to the client
        {
          errors: contact_address.errors.full_messages
        }
      end
    end
  end
end

module Mutations
  class CreateContactAddress < Mutations::BaseMutation
    null true

    ContactAddress.define_graphql_mutation(self, Types::Models::ContactAddressType)

    def resolve(contact_id:, **params)
      contact_address = context[:current_user].contacts.find(contact_id).addresses.build(params)
      if contact_address.save
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

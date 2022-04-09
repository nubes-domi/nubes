module Mutations
  class AddContactAddress < Mutations::BaseMutation
    null true

    ContactAddress.define_graphql_mutation(self, Types::Models::ContactAddressType)

    def resolve(contact_id:, **attributes)
      result = Contacts::AddAddress.call(user: context[:current_user], contact_id:, attributes:)
      handle_failures(result) do |contact_address|
        { contact_address: }
      end
    end
  end
end

module Mutations
  class UpdateContactAddress < Mutations::BaseMutation
    null true

    argument :contact_id, ID
    ContactAddress.define_graphql_mutation(self, Types::Models::ContactAddressType, type: :update)

    def resolve(id:, contact_id:, **attributes)
      result = Contacts::UpdateAddress.call(user: context[:current_user], id:, contact_id:, attributes:)
      handle_failures(result) do |contact_address|
        { contact_address: }
      end
    end
  end
end
